package bot

import (
	"agent-care-tg/models"
	tz "github.com/bradfitz/latlong"
	tg "gopkg.in/telebot.v3"
)

type Handler struct {
	bot      *tg.Bot
	state    map[int64]string
	userData map[int64]*models.User
}

func NewHandler(bot *tg.Bot) *Handler {
	return &Handler{bot: bot, state: make(map[int64]string), userData: make(map[int64]*models.User)}
}

func (h *Handler) Register() {
	h.bot.Handle("/start", h.handleStart)
	h.bot.Handle("Get Started", h.handleGetStarted)
	h.bot.Handle(tg.OnText, h.handleUserRegistration)
	h.bot.Handle(tg.OnLocation, h.handleUserRegistration)
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
		if err := c.Send(msg); err != nil {
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
	h.userData[c.Chat().ID] = &models.User{}
	h.state[c.Chat().ID] = "waiting_for_name"
	return c.Send("What should I call you? (Just type your name or whatever you'd like to go by!)")
}

func (h *Handler) handleUserRegistration(c tg.Context) error {
	switch h.state[c.Chat().ID] {
	case "waiting_for_name":
		h.userData[c.Chat().ID].PrefferedName = c.Text()
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
		// ToDo: save to db
		// ToDo: Get rid of buttons
		return c.Send("Perfect! I am now setting up your profile")

	}
	return nil
}
