package bot

import (
	"agent-care-tg/models"
	"agent-care-tg/storage"
	"fmt"
	"log/slog"
	"strings"
	"sync"

	tz "github.com/bradfitz/latlong"
	tg "gopkg.in/telebot.v3"
)

type Handler struct {
	bot      *tg.Bot
	mu       sync.RWMutex
	state    map[int64]string
	userData map[int64]*models.User
	store    *storage.Store
}

var feedbackIndex = 0

func NewHandler(bot *tg.Bot, store *storage.Store) *Handler {
	return &Handler{bot: bot, state: make(map[int64]string), userData: make(map[int64]*models.User), store: store}
}

func (h *Handler) Register() {
	h.bot.Handle("/start", h.handleStart)
	h.bot.Handle("/profile", h.handleProfile)
	h.bot.Handle("/streak", h.handleStreak)
	h.bot.Handle("Learn how it works", h.handleLearnHowItWorks)
	h.bot.Handle("Get Started", h.handleGetStarted)
	// This handles all the answers and their respective responses in user regiastration//
	h.bot.Handle(tg.OnText, h.handleUserRegistration)
	h.bot.Handle(tg.OnLocation, h.handleUserRegistration)
	// handling reponse to task check-ins
	h.bot.Handle(tg.OnCallback, func(c tg.Context) error {
		data := c.Callback().Data
		if strings.Contains(data, "task_completed") {
			return h.handleTaskCompleted(c)
		}
		if strings.Contains(data, "task_skipped") {
			return h.handleTaskSkipped(c)
		}
		return nil
	})
}

func (h *Handler) handleStart(c tg.Context) error {
	// 1. Send welcom message
	c.Send(MsgWelcome, tg.ModeMarkdown)

	// 2. Send data collection message and close keyboard
	removeKeyboard := &tg.ReplyMarkup{RemoveKeyboard: true}
	c.Send(MsgDataCollection, tg.ModeMarkdown, removeKeyboard)

	// 3. Present button for how it works
	markup := &tg.ReplyMarkup{ResizeKeyboard: true}
	btnHowItWorks := markup.Text("Learn how it works")
	markup.Reply(markup.Row(btnHowItWorks))
	return c.Send("Ready to learn how it works ?", markup)
}

func (h *Handler) handleLearnHowItWorks(c tg.Context) error {
	markup := &tg.ReplyMarkup{ResizeKeyboard: true}
	btnGetStarted := markup.Text("Get Started")
	markup.Reply(markup.Row(btnGetStarted))
	return c.Send(MsgHowItWorks, markup, tg.ModeMarkdown)
}

func (h *Handler) handleGetStarted(c tg.Context) error {
	h.mu.Lock()
	h.userData[c.Chat().ID] = models.NewUser()
	h.state[c.Chat().ID] = "waiting_for_name"
	h.mu.Unlock()

	removeKeyboard := &tg.ReplyMarkup{RemoveKeyboard: true}
	return c.Send("What should I call you?", removeKeyboard)
}

func (h *Handler) handleUserRegistration(c tg.Context) error {
	h.mu.RLock()
	state := h.state[c.Chat().ID]
	h.mu.RUnlock()

	switch state {
	case "waiting_for_name":
		h.mu.Lock()
		h.userData[c.Chat().ID].Username = c.Text()
		h.state[c.Chat().ID] = "waiting_for_goal"
		h.mu.Unlock()
		return c.Send("Nice to meet you " + c.Text() + "!" + "\n\nWhat's your personal goal?")

	case "waiting_for_goal":
		h.mu.Lock()
		h.userData[c.Chat().ID].PersonalGoal = c.Text()
		h.state[c.Chat().ID] = "waiting_for_timezone"
		h.mu.Unlock()
		markup := &tg.ReplyMarkup{ResizeKeyboard: true, OneTimeKeyboard: true}
		locationBtn := markup.Location("Share my location")
		markup.Reply(markup.Row(locationBtn))
		return c.Send(("Almost done! Please share your location so you can get reminders in your timezone"), markup)

	case "waiting_for_timezone":
		lat := c.Message().Location.Lat
		lng := c.Message().Location.Lng
		timezone := tz.LookupZoneName(float64(lat), float64(lng))

		if timezone == "" {
			return c.Send("Sorry, I coulnd't detect your timezone. Please try again")
		}

		h.mu.Lock()
		user := h.userData[c.Chat().ID]
		user.ChatID = c.Chat().ID
		user.TGUsername = c.Sender().Username
		user.Timezone = timezone
		delete(h.state, c.Chat().ID)
		delete(h.userData, c.Chat().ID)
		h.mu.Unlock()
		removeKeyboard := &tg.ReplyMarkup{RemoveKeyboard: true}

		if err := h.store.SaveUser(user); err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				return c.Send("You are already registered", removeKeyboard)
			}
			slog.Error("Failed to save user", "error", err)
			return c.Send("Something went wrong with your profile. Please try again later", removeKeyboard)
		}
		//ToDo: Send a prep message
		slog.Info("New user registered", "username", user.TGUsername)
		c.Send("Thanks ! I am now setting up your profile...", removeKeyboard)
		return c.Send("Perfect! You are all setup")

	}
	return nil
}

func (h *Handler) handleTaskCompleted(c tg.Context) error {
	// Gets rid of the button from the markup to avoid race condtion and bad operations.
	c.Edit(c.Callback().Message.Text, &tg.ReplyMarkup{})

	callBackData := strings.TrimSpace(c.Callback().Data)
	taskTag := models.TaskTag(strings.Replace(callBackData, "_task_completed", "", 1))
	chatID := c.Chat().ID

	err := h.store.IncrementStreak(chatID, taskTag)
	if err != nil {
		slog.Error("Failed to update streak", "err", err)
		c.Send("Oops, something went wong. We couldn't update your streak")
		c.Respond()
		return fmt.Errorf("Failed to update streak: %w", err)
	}

	err = h.handleMaxStreak(chatID, taskTag)
	if err != nil {
		slog.Error("Failed to update max streak", "err", err)
		c.Send("Oops, something went wong. We couldn't update your max streak")
		c.Respond()
		return fmt.Errorf("Failed to update max streak: %w", err)
	}

	c.Send(GetFeedbackMessage(taskTag))
	slog.Info("Task completed clicked", "data", taskTag)
	c.Respond()
	return nil
}

func (h *Handler) handleTaskSkipped(c tg.Context) error {
	// Gets rid of the button from the markup to avoid race condtion and bad operations.
	c.Edit(c.Callback().Message.Text, &tg.ReplyMarkup{})

	callBackData := strings.TrimSpace(c.Callback().Data)
	taskTag := models.TaskTag(strings.Replace(callBackData, "_task_skipped", "", 1))
	chatID := c.Chat().ID

	err := h.store.ResetStreak(chatID, taskTag)
	if err != nil {
		slog.Error("Failed to reset streak", "err", err)
		c.Send("Oops, something went wong.")
		c.Respond()
		return fmt.Errorf("Failed to reset streak: %w", err)
	}
	c.Send("Its Okay")
	c.Respond()
	return nil
}

func (h *Handler) handleMaxStreak(chatID int64, taskTag models.TaskTag) error {
	task, err := h.store.GetTask(chatID, taskTag)

	if err != nil {
		slog.Error("Failed to get user data", "error", err)
		return fmt.Errorf("Failed to get user data: %w", err)
	}

	maxStreak := task.MaxStreak
	slog.Info("max streak", "maxStreak", maxStreak)

	if task.MaxStreak > task.CurrentStreak {
		slog.Info("ffff")
		return nil
	}

	err = h.store.UpdateMaxStreak(task.ID, int64(task.CurrentStreak))
	if err != nil {
		slog.Error("Failed to update max streak", "error", err)
		return fmt.Errorf("Failed to update max streak %w", err)
	}
	return nil
}

func (h *Handler) handleProfile(c tg.Context) error {
	chatID := c.Chat().ID

	user, err := h.store.GetUserByChatID(chatID)

	if err != nil {
		slog.Warn("Can not get user data for profile", "warning", err)
		return c.Send("Sorry, cant find ur profile, are u sure u are registered")
	}

	userName := user.Username
	userGoal := user.PersonalGoal
	userTimezone := user.Timezone
	userJoinedAtFormatted := user.JoinedAt.Format("Jan 02, 2006") // Format JoinedAt
	formattedMsg := fmt.Sprintf(MsgProfileData, userName, userGoal, userTimezone, userJoinedAtFormatted)
	c.Send(formattedMsg, tg.ModeMarkdown)
	return nil
}

// ToDo: Refactor the string builder to be its own util function//
func (h *Handler) handleStreak(c tg.Context) error {
	chatID := c.Chat().ID

	tasks, err := h.store.GetAllTasksForUserByChatID(chatID)

	if err != nil {
		slog.Error("Can not get all the tasks for the give user", "error", err)
		// Assuming this message implies registration issue, similar to handleProfile
		return c.Send("Sorry, cant find your streak data. Are you sure you are registered and have tasks?")
	}

	var messageBuilder strings.Builder
	messageBuilder.WriteString(MsgStreaksHeader)

	if len(tasks) == 0 {
		messageBuilder.WriteString("\n\nYou don't have any active tasks yet!\nStart by setting some goals.")
	} else {
		for _, task := range tasks {
			// Only display active tasks for streaks
			if task.IsActive {
				taskEmoji := "✅" // Default emoji
				switch task.Tag {
				case models.TagMorning:
					taskEmoji = "⏰"
				case models.TagSunlight:
					taskEmoji = "☀️"
				case models.TagExercise:
					taskEmoji = "💪"
				case models.TagMeal:
					taskEmoji = "🥗"
				case models.TagPersonal:
					taskEmoji = "📔" // Journal emoji for personal goal
				}

				currentStreakEmoji := "⚡️"
				if task.CurrentStreak == 0 {
					currentStreakEmoji = "🥶"
				}
				formattedTask := fmt.Sprintf(MsgTaskStreak, taskEmoji, task.Name, task.CurrentStreak, currentStreakEmoji, task.MaxStreak)
				messageBuilder.WriteString(formattedTask)
			}
		}
		// If no active tasks were found, add a specific message
		if messageBuilder.Len() == len(MsgStreaksHeader) {
			messageBuilder.WriteString("\n\nYou don't have any active tasks to track streaks for.\nConsider enabling some tasks!")
		}
	}

	return c.Send(messageBuilder.String(), tg.ModeMarkdown)
}
