# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

### Feature Under Planning
1. Custom reminders for users. Useful for things like medication, study sessions, water intake, etc.
    a. Create Task
    b. Delete Task
    c. Update Task - Typos, Time, etc.
2. Steak Check
    a. Ability to check streaks for each tasks and 
    b. Abilty to check days a task was performed since beginning
---
## [v 0.1.0] - 2026-06-23 (Initial Release)

This release marks the initial public version of Agent Care, a Telegram bot designed to help users build momentum and re-engage with daily tasks. It establishes the core framework for user interaction, task management, and deployment.

### Added
-   **Core Bot Infrastructure**: Initial project setup, database connection, and Telegram bot bootstrapping.
-   **User Registration Flow**: Comprehensive conversation flow for new user registration and saving user profiles to the database, including `joined_at` timestamp.
-   **Daily Check-in System**:
    -   Scheduler implementation with cron jobs to send messages to users with correct timezones.
    -   Functionality to send morning messages and sun light check-ins.
    -   Base logic for handling user responses to daily tasks.
    -   Buttons for user responses to interactive messages.
    -   Default tasks automatically created for each new user.
-   **Task and Streak Management**:
    -   Database schema and functionality for storing user tasks.
    -   Base logic for tracking and updating streaks when tasks are completed.
    -   Function to reset the current streak if a task is skipped.
    -   Functions to manage and track maximum streaks.
-   **Feedback Messages**:
    -   Initial positive feedback messages for morning routines, sunlight, exercise, meals, and personal tasks.
    -   `setbackMessages` array for encouraging feedback during challenges.
-   **Bot Commands**:
    -   `/profile` command with a dedicated message template.
    -   `/streak` command with a message template.
-   **Logging**: Implemented `slog` for structured logging in handler, main, and config files.
-   **Deployment Automation**:
    -   GitHub Actions workflow (`deploy.yml`) for automated deployment on `push` to `main`.
    -   GCP authentication, Docker image building, saving, and copying to VM.
    -   Secure handling of `.env` variables for deployment.
    -   Deployment script (`deploy.sh`) for managing Docker containers on the VM.
    -   `DEPLOYMENT.md` documentation detailing the deployment process.
-   **Database Migrations**: Integrated `goose` for safe and version-controlled database migrations, including a migration script for `last_sent_at` on users.
-   **Configuration**: Added environment variable checks for deployment.
-   **Helper Utilities**: Added functions to manage message button clean-ups and edit messages after expiry time.
-   **Error Handling**: Basic error handling for existing users during registration.
-   **Constants**: Implemented constants for task tags for centralization and typesafety.

### Changed
-   **Database Driver**: Switched to `pgx` for better performance and `jsonb` support with `scan` and `value` interfaces.
-   **Message Handling**: Refactored routine messages to use `sendMessageToAllUsers`.
-   **User Registration**: Cleaned up user registration logic.
-   **`GetStarted` Function**: Refactored for cleaner implementation.
-   **Positive Feedback Logic**: Improved message retrieval in `GetFeedbackMessage` to prevent out-of-bounds errors and ensure cycling.
-   **Deployment Scripts**: Refactored deployment scripts for cleaner VM state and safer environment variable handling, adapted for `main` branch deployment.
-   **README.md**: Multiple refactorings and updates for better clarity and structure, including the addition of dynamic badges for project metadata.

### Removed

### Fixed
---
Template

## [Version] - Release Date
### Added
### Changed
### Deprecated
### Removed
### Fixed
### Security
