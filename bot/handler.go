package bot

import (
	"agent-care-tg/models"
	"agent-care-tg/storage"
	"log"
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
	h.bot.Handle("Get Started", h.handleGetStarted)
	h.bot.Handle(tg.OnText, h.handleUserRegistration)
	h.bot.Handle(tg.OnLocation, h.handleUserRegistration)
	// handling reponse to task check-ins
	h.bot.Handle(tg.OnCallback, func(c tg.Context) error {
		data := c.Callback().Data
		if strings.Contains(data, "task_done") {
			return h.handleTaskCompleted(c)
		}
		if strings.Contains(data, "task_skipped") {
			return h.handleTaskSkipped(c)
		}
		return nil
	})
}

func (h *Handler) handleStart(c tg.Context) error {
	messages := []string{
		MsgWelcome,
		MsgDataCollection,
		MsgHowItWorks,
		MsgCheckInRule,
		MsgInvite,
	}

	for _, msg := range messages {
		if err := c.Send(msg, tg.ModeMarkdown); err != nil {
			return err
		}
	}
	// trigger next step
	markup := &tg.ReplyMarkup{ResizeKeyboard: true, OneTimeKeyboard: true}
	btnGetStarted := markup.Text("Get Started")
	markup.Reply(markup.Row(btnGetStarted))
	return c.Send("Are you ready ?", markup)
}

func (h *Handler) handleGetStarted(c tg.Context) error {
	h.userData[c.Chat().ID] = models.NewUser()
	h.state[c.Chat().ID] = "waiting_for_name"
	return c.Send("What should I call you? (Just type your name or whatever you'd like to go by!)")
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
			log.Println("error", err)
			return c.Send("Something went wrong with your profile. Please try again")
		}
		// ToDo: Get rid of buttons
		//ToDo: Send a prep message
		log.Println("New user registered", user.TGUsername)
		c.Send("Thanks ! I am now setting up your profile...")
		return c.Send("Perfect! You are all setup")

	}
	return nil
}

func (h *Handler) handleTaskCompleted(c tg.Context) error {
	log.Println("Task completed clicked")
	// mark streak
	// send a positive message
	// use c.respond //
	return nil
}

func (h *Handler) handleTaskSkipped(c tg.Context) error {
	log.Println("Task skipped clicked")
	// send a supportive message
	// use c.respond
	return nil
}
