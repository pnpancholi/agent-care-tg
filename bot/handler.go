package bot

import (
	"agent-care-tg/models"
	"agent-care-tg/storage"
	"fmt"
	"log/slog"
	"strings"

	tz "github.com/bradfitz/latlong"
	tg "gopkg.in/telebot.v3"
)

type Handler struct {
	bot      *tg.Bot
	state    map[int64]string
	userData map[int64]*models.User
	store    *storage.Store
}

func NewHandler(bot *tg.Bot, store *storage.Store) *Handler {
	return &Handler{bot: bot, state: make(map[int64]string), userData: make(map[int64]*models.User), store: store}
}

func (h *Handler) Register() {
	h.bot.Handle("/start", h.handleStart)
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
	c.Send(MsgDataCollection, removeKeyboard, tg.ModeMarkdown)

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
	return c.Send(MsgHowItWorks, markup)
}

func (h *Handler) handleGetStarted(c tg.Context) error {
	h.userData[c.Chat().ID] = models.NewUser()
	h.state[c.Chat().ID] = "waiting_for_name"
	return c.Send("What should I call you?")
}

func (h *Handler) handleUserRegistration(c tg.Context) error {
	switch h.state[c.Chat().ID] {
	case "waiting_for_name":
		h.userData[c.Chat().ID].Username = c.Text()
		h.state[c.Chat().ID] = "waiting_for_goal"
		return c.Send("Nice to meet you " + c.Text() + "!" + "\n\nWhat's your personal goal?")

	case "waiting_for_goal":
		h.userData[c.Chat().ID].PersonalGoal = c.Text()
		h.state[c.Chat().ID] = "waiting_for_timezone"
		markup := &tg.ReplyMarkup{ResizeKeyboard: true, OneTimeKeyboard: true}
		locationBtn := markup.Location("Share my location")
		markup.Reply(markup.Row(locationBtn))
		return c.Send(("Almost done! Please share your location so you can get reminders in your timezone"), markup)

	case "waiting_for_timezone":
		//capture response
		lat := c.Message().Location.Lat
		lng := c.Message().Location.Lng
		timezone := tz.LookupZoneName(float64(lat), float64(lng))
		if timezone == "" {
			return c.Send("Sorry, I coulnd't detect your timezone. Please try again")
		}
		h.userData[c.Chat().ID].Timezone = timezone
		// ToDo: clear out state
		user := h.userData[c.Chat().ID]
		user.ChatID = c.Chat().ID
		user.TGUsername = c.Sender().Username
		user.Timezone = timezone

		if err := h.store.SaveUser(user); err != nil {
			slog.Error("Failed to save user", "error", err)
			return c.Send("Something went wrong with your profile. Please try again later")
		}
		// ToDo: Get rid of buttons
		//ToDo: Send a prep message
		slog.Info("New user registered", "username", user.TGUsername)
		c.Send("Thanks ! I am now setting up your profile...")
		return c.Send("Perfect! You are all setup")

	}
	return nil
}

func (h *Handler) handleTaskCompleted(c tg.Context) error {
	callBackData := strings.TrimSpace(c.Callback().Data)
	taskTag := strings.Replace(callBackData, "_task_completed", "", 1)
	chatID := c.Chat().ID

	err := h.store.IncrementStreak(chatID, taskTag)

	if err != nil {
		slog.Error("Failed to update streak", "err", err)
		c.Send("Oops, something went wong. We couldn't update your streak")
		c.Respond()
		return fmt.Errorf("Failed to update streak: %w", err)
	}
	c.Send("Great job, keep it up!")
	slog.Info("Task completed clicked", "data", taskTag)
	c.Respond()
	return nil
}

func (h *Handler) handleTaskSkipped(c tg.Context) error {
	callBackData := strings.TrimSpace(c.Callback().Data)
	taskTag := strings.Replace(callBackData, "_task_skipped", "", 1)
	chatID := c.Chat().ID

	err := h.store.ResetStreak(chatID, taskTag)
	if err != nil {
		slog.Error("Failed to reset streak", "err", err)
		c.Send("Oops, something went wong.")
		c.Respond()
		return fmt.Errorf("Failed to reset streak: %w", err)
	}
	c.Send("Its Okay")
	slog.Info("Task skipped clicked", "data", c)
	c.Respond()
	return nil
}
