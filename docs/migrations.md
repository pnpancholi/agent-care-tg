# Database Migrations

This project uses [goose](https://github.com/pressly/goose) to manage schema changes.

---

## Rules

- **Never edit existing migration files** once they have been run
- **Always create a new migration** for any schema change
- **Commit migration files** to version control
- **Never modify the database manually** — always through a migration

---

## Setup

Make sure your `DATABASE_URL` is exported:

```bash
export $(cat .env | xargs)
```

---

## Commands

**Run all pending migrations**
```bash
goose -dir migrations postgres "$DATABASE_URL" up
```

**Rollback last migration**
```bash
goose -dir migrations postgres "$DATABASE_URL" down
```

**Check status of all migrations**
```bash
goose -dir migrations postgres "$DATABASE_URL" status
```

**Create a new migration**
```bash
goose -dir migrations create migration_name sql
```

---

## Migration File Structure

Every migration file has two sections:

```sql
-- +goose Up
-- SQL to apply the change
CREATE TABLE users (...);

-- +goose Down
-- SQL to reverse the change
DROP TABLE IF EXISTS users;
```

`Up` runs when you migrate forward. `Down` runs when you rollback.

---

## Migration History

| Version | Name | Description |
|---|---|---|
| 20260316080325 | create_users | Initial users table |
