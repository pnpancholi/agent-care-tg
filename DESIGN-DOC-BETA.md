# DESIGN-DOC-BETA: Agent Care (Telegram Bot)

## 1. Overview
Agent Care is a specialized Telegram bot built to act as a gentle, low-friction accountability companion. It is specifically designed to help individuals navigating depressive episodes or long-term ruts (such as NEETs) to slowly build momentum and reclaim their sense of agency through incremental daily habits.

## 2. Architecture & Tech Stack
- **Language**: Go (v1.22/v1.25)
- **Framework**: `telebot v3` for Telegram Bot API interactions.
- **Database**: PostgreSQL (hosted on Supabase) utilizing `sqlx` for ORM capabilities.
- **Migrations**: Version-controlled SQL migrations managed via `goose`.
- **Scheduling**: In-memory cron jobs using `robfig/cron/v3`.
- **Timezone Management**: Automated timezone detection from coordinates using `bradfitz/latlong`.
- **Deployment**: Dockerized container deployed to a Google Cloud VM via GitHub Actions (CI/CD).

## 3. Core Features

### 3.1. Timezone-Aware Routine Scheduling
The bot triggers background jobs every 10 minutes to check if users have reached specific local hours in their respective timezones. It schedules the following core check-ins:
- **Morning Check-In**: 07:00 AM
- **Sunlight Check-In**: 02:00 PM (14:00)
- **Exercise Check-In**: 05:00 PM (17:00)
- **Healthy Meal Check-In**: 08:00 PM (20:00)
- **Personal Goal Check-In**: 09:00 PM (21:00)

### 3.2. Interactive Inline Responses
Messages are sent with Inline Keyboards (buttons) allowing users to easily respond with:
- **"Done"**: Registers task completion.
- **"Skipped"**: Acknowledges that the task was missed.

To avoid retroactive polling and database pollution, message buttons are configured to expire and disappear after a short window (currently 2 minutes).

### 3.3. Empathetic Feedback Loop
- **Positive Reinforcement**: Clicking "Done" triggers randomized positive reinforcement messages tailored to the specific activity.
- **Setback Support**: Clicking "Skipped" replies with gentle, non-judgmental setback messages emphasizing that recovery is non-linear and tomorrow is another opportunity.

### 3.4. Guardrails & Concurrency
- **Double-Send Prevention**: A `last_sent_at` timestamp is enforced per user. The bot cross-references the current day and hour to ensure a user receives exactly one message per scheduled block.
- **Concurrency**: Operations like message expirations (`scheduleExpiry`) run concurrently in goroutines to avoid blocking the main scheduler timeline.

## 4. Data Model (High-Level)
The primary data entity is the **User**:
- `chat_id` / `tg_username`: For routing Telegram messages.
- `timezone`: Stored string (e.g., "Asia/Kolkata") parsed from user location sharing.
- `last_sent_at`: Timestamp (Guards against duplicate cron-triggered messages).
- *Tracking/Streaks*: Fields to record adherence (current streak, max streak, tasks completed).

## 5. Deployment Flow
1. **GitHub Actions**: Triggers on push to `main` and `dev` branches.
2. **Build**: Packages the Go binary into a lightweight `alpine` Docker Image.
3. **Transfer**: Pushes the compressed image to a GCP Virtual Machine via SSH/SCP.
4. **Run**: Replaces the old running container with the newly built Docker container seamlessly passing in production environment variables (`.env`).
