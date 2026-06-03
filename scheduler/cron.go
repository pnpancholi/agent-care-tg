package scheduler

// ToDo : Add a logger for each cron job, type of message/reminder. number of users sent to, time, and error if any
import (
	bot "agent-care-tg/bot"
	"agent-care-tg/models"
	"agent-care-tg/storage"
	"fmt"
	"log/slog"
	"time"

	"github.com/robfig/cron/v3"
	tg "gopkg.in/telebot.v3"
)

type Scheduler struct {
	cron  *cron.Cron
	store *storage.Store
	bot   *tg.Bot
}

const (
	MorningTag      = "daily_morning"
	SunlightTag     = "daily_sunlight"
	ExcercoseTag    = "daily_excercise"
	PersonalGoalTag = "daily_personal"
	HealthyMealTag  = "daily_meal"
)

func New(store *storage.Store, bot *tg.Bot) *Scheduler {
	return &Scheduler{
		cron:  cron.New(),
		store: store,
		bot:   bot,
	}
}

func (s *Scheduler) Start() {
	s.cron.AddFunc("*/2 * * * *", func() {
		s.testMessage()
	})
	s.cron.AddFunc("*/10 * * * *", func() {
		s.sendMorningMessage(7)
		s.checkInForSunlight(14)
		s.checkInForHealthyMeal(14)
		s.checkInForPersonalGoal(21)
		s.checkInForExcercise(17)
	})

	s.cron.Start()
	slog.Info("Scheduler started.")
}

func (s *Scheduler) Stop() {
	s.cron.Stop()
	slog.Info("Scheduler stopped.")
}

func (s *Scheduler) testMessage() {
	taskTag := "daily_morning"
	msg := "Test message for"
	users, err := s.store.GetAllUsers()

	if err != nil {
		slog.Error("Failed to access users from DB for sendMessageToAllUsersInTimeZone", "error", err)
		return
	}

	for _, user := range users {

		markup := &tg.ReplyMarkup{}
		doneBtn := markup.Data("Done", taskTag+"_task_completed")
		skippedBtn := markup.Data("Skipped", taskTag+"_task_skipped")
		markup.Inline(markup.Row(doneBtn, skippedBtn))
		formattedMsg := fmt.Sprintf(msg, user.Username)

		if _, err := s.bot.Send(tg.ChatID(user.ChatID), formattedMsg, markup, tg.ModeMarkdown); err != nil {
			slog.Error("Failed to send message to : ", "username", user.TGUsername, "error", err)
			continue
		}

		if err := s.store.UpdateLastSentAt(&user); err != nil {
			slog.Error("Failed to update last sent at for : ", "username", user.TGUsername, "error", err)
			continue
		}

		slog.Info("Updated last_sent_at timestampe for the user")
	}

}
func (s *Scheduler) sendMorningMessage(localHour uint8) {
	s.sendMessageToAllUsersInTimeZone(localHour, MorningTag, bot.MsgMorningCheckIn)
}

func (s *Scheduler) checkInForSunlight(localHour uint8) {
	s.sendMessageToAllUsersInTimeZone(localHour, SunlightTag, bot.MsgSunlightCheckIn)
}

func (s *Scheduler) checkInForHealthyMeal(localHour uint8) {
	s.sendMessageToAllUsersInTimeZone(localHour, HealthyMealTag, bot.MsgMealCheckIn)
}

func (s *Scheduler) checkInForPersonalGoal(localHour uint8) {
	s.sendMessageToAllUsersInTimeZone(localHour, PersonalGoalTag, bot.MsgPersonalGoalCheckIn)
}

func (s *Scheduler) checkInForExcercise(localHour uint8) {
	s.sendMessageToAllUsersInTimeZone(localHour, ExcercoseTag, bot.MsgExcerciseCheckIn)
}

func lastSentCheckPassed(user *models.User, hour uint8) bool {
	loc, err := time.LoadLocation(user.Timezone)
	if err != nil {
		slog.Error("Failed to load location for lastSentCheckPassed", "error", err)
		return false
	}
	localTime := time.Now().In(loc)

	if !user.LastSentAt.Valid {
		return true
	}
	if localTime.Hour() != int(hour) || localTime.Minute() > 10 {
		return false
	}

	lastSentLocalTime := user.LastSentAt.Time.In(loc)

	if lastSentLocalTime.Year() == localTime.Year() &&
		lastSentLocalTime.Month() == localTime.Month() &&
		lastSentLocalTime.Day() == localTime.Day() &&
		lastSentLocalTime.Hour() == int(hour) &&
		lastSentLocalTime.Minute() >= 10 {
		slog.Warn("Double message check ", "username", user.TGUsername)
		return false
	}

	return true
}

func (s *Scheduler) sendMessageToAllUsersInTimeZone(hour uint8, taskTag string, msg string) {
	users, err := s.store.GetAllUsers()

	if err != nil {
		slog.Error("Failed to access users from DB for sendMessageToAllUsersInTimeZone", "error", err)
		return
	}

	for _, user := range users {

		if !lastSentCheckPassed(&user, hour) {
			continue
		}

		markup := &tg.ReplyMarkup{}
		doneBtn := markup.Data("Done", taskTag+"_task_completed")
		skippedBtn := markup.Data("Skipped", taskTag+"_task_skipped")
		markup.Inline(markup.Row(doneBtn, skippedBtn))
		formattedMsg := fmt.Sprintf(msg, user.Username)

		if _, err := s.bot.Send(tg.ChatID(user.ChatID), formattedMsg, markup, tg.ModeMarkdown); err != nil {
			slog.Error("Failed to send message to : ", "username", user.TGUsername, "error", err)
			continue
		}

		if err := s.store.UpdateLastSentAt(&user); err != nil {
			slog.Error("Failed to update last sent at for : ", "username", user.TGUsername, "error", err)
			continue
		}

		slog.Info("Updated last_sent_at timestampe for the user")
	}
}
