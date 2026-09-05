# PENDING_EXPIRATIONS_SOLUTION

## Problem

When the bot restarts, `scheduleExpiry` goroutines die. Messages sent before the restart retain their inline buttons forever because no cleanup mechanism persists across restarts.

## Root Cause

The current approach spawns a goroutine per message that sleeps for 60 minutes then edits the message to remove buttons. This is entirely in-memory — no DB persistence.

## Solution

Replace the goroutine-based approach with a DB-backed cron job. Store pending expirations in a dedicated table, and run a cron job every 10 minutes to clean up expired messages.

---

## Implementation Plan

### 1. New Migration

**Create:** `migrations/20260905000000_create_pending_expirations.sql`

```sql
-- +goose Up
CREATE TABLE pending_expirations (
    id SERIAL PRIMARY KEY,
    chat_id BIGINT NOT NULL REFERENCES users(chat_id),
    message_id INT NOT NULL,
    message_text TEXT NOT NULL,
    expires_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE pending_expirations;
```

### 2. New Model

**Edit:** `models/models.go` — Add after `Task` struct:

```go
type PendingExpiration struct {
    ID          int       `json:"id" db:"id"`
    ChatID      int64     `json:"chat_id" db:"chat_id"`
    MessageID   int       `json:"message_id" db:"message_id"`
    MessageText string    `json:"message_text" db:"message_text"`
    ExpiresAt   time.Time `json:"expires_at" db:"expires_at"`
}
```

### 3. Storage Methods

**Edit:** `storage/store.go` — Add 3 methods:

```go
func (s *Store) SavePendingExpiration(chatID int64, messageID int, messageText string, expiresAt time.Time) error {
    query := `INSERT INTO pending_expirations (chat_id, message_id, message_text, expires_at) VALUES ($1, $2, $3, $4)`
    _, err := s.db.Exec(query, chatID, messageID, messageText, expiresAt)
    return err
}

func (s *Store) GetExpiredExpirations() ([]models.PendingExpiration, error) {
    var expirations []models.PendingExpiration
    query := `SELECT * FROM pending_expirations WHERE expires_at < NOW()`
    err := s.db.Select(&expirations, query)
    return expirations, err
}

func (s *Store) DeletePendingExpiration(chatID int64, messageID int) error {
    query := `DELETE FROM pending_expirations WHERE chat_id = $1 AND message_id = $2`
    _, err := s.db.Exec(query, chatID, messageID)
    return err
}
```

### 4. Scheduler Changes

**Edit:** `scheduler/cron.go`

#### 4a. Add `expireOldButtons` function

```go
func (s *Scheduler) expireOldButtons() {
    expirations, err := s.store.GetExpiredExpirations()
    if err != nil {
        slog.Error("Failed to get expired expirations", "error", err)
        return
    }

    for _, exp := range expirations {
        msg := &tg.StoredMessage{MessageID: fmt.Sprintf("%d", exp.MessageID), ChatID: exp.ChatID}
        markup := &tg.ReplyMarkup{}
        _, err := s.bot.Edit(msg, exp.MessageText, markup)
        if err != nil {
            slog.Error("Failed to revoke action buttons", "chat_id", exp.ChatID, "message_id", exp.MessageID, "error", err)
        }
        s.store.DeletePendingExpiration(exp.ChatID, exp.MessageID)
    }
}
```

#### 4b. Register cron job in `Start()`

```go
s.cron.AddFunc("*/10 * * * *", func() {
    s.expireOldButtons()
})
```

#### 4c. Store message on send in `sendMessageToAllUsersInTimeZone`

Replace `s.scheduleExpiry(msg)` with:

```go
expiresAt := time.Now().Add(EXPIRATION_TIME * time.Minute)
s.store.SavePendingExpiration(user.ChatID, msg.ID, formattedMsg, expiresAt)
```

#### 4d. Same change in `testMessage()`

Replace `s.scheduleExpiry(msg)` with:

```go
expiresAt := time.Now().Add(EXPIRATION_TIME * time.Minute)
s.store.SavePendingExpiration(user.ChatID, msg.ID, formattedMsg, expiresAt)
```

#### 4e. Delete `scheduleExpiry` function

Remove lines 183-192 (the entire `scheduleExpiry` function).

### 5. Bot Handler — Clean up on click

**Edit:** `bot/handler.go`

In `handleTaskCompleted` (after line 172):

```go
h.store.DeletePendingExpiration(c.Chat().ID, c.Message().ID)
```

In `handleTaskSkipped` (after line 204):

```go
h.store.DeletePendingExpiration(c.Chat().ID, c.Message().ID)
```

---

## Files Changed

| File | Action |
|------|--------|
| `migrations/20260905000000_create_pending_expirations.sql` | Create new |
| `models/models.go` | Add `PendingExpiration` struct |
| `storage/store.go` | Add 3 methods |
| `scheduler/cron.go` | Add `expireOldButtons`, register cron, store on send, delete `scheduleExpiry` |
| `bot/handler.go` | Add `DeletePendingExpiration` call in 2 handlers |

---

## Edge Cases

| Scenario | Behavior |
|----------|----------|
| User clicks button before 60 min | Buttons removed + row deleted → expiry cron finds nothing |
| Bot restarts | Rows persist in DB → cron cleans up on next run |
| Edit fails (message deleted, rate limit) | Error logged → row deleted to avoid retry loops |

---

## Technical Notes

- `tg.StoredMessage` implements the `Editable` interface with just `ChatID` + `MessageID` — used for `bot.Edit()` without a full `*tg.Message`
- `c.Message().ID` gives the message ID in callback context
- `message_text` is stored because Telegram's `Edit` method requires the full message text to update markup
