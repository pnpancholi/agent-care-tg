package scheduler

// ToDo : Add a logger for each cron job, type of message/reminder. number of users sent to, time, and error if any
import (
	"agent-care-tg/storage"
	"github.com/robfig/cron/v3"
	tg "gopkg.in/telebot.v3"
	"log"
)

type Scheduler struct {
	cron  *cron.Cron
	store *storage.Store
	bot   *tg.Bot
}

func New(store *storage.Store, bot *tg.Bot) *Scheduler {
	return &Scheduler{
		cron:  cron.New(cron.WithSeconds()),
		store: store,
		bot:   bot,
	}
}

func (s *Scheduler) Start() {
	s.cron.AddFunc("*/10 * * * * *", func() {
		log.Println("Scheduler Fired")
	})
	//Truggering Cron Jobs
	s.cron.AddFunc("*/10 * * * * *", s.SendMorningMessage)
	s.cron.AddFunc("*/10 * * * * *", s.CheckInForSunlight)

	s.cron.Start()
	log.Println("Scheduler Started...")
}

func (s *Scheduler) Stop() {
	s.cron.AddFunc("", func() {
		log.Println("Scheduler Stopped...")
	})
}

func (s *Scheduler) SendMorningMessage() {
	users, err := s.store.GetAllUsers()
	if err != nil {
		log.Println("Can not access users from DB", err)
		return
	}

	for _, user := range users {
		msg := "morning msg"
		_, err := s.bot.Send(tg.ChatID(user.ChatID), msg)
		if err != nil {
			log.Println("Failed to send user their morning message", user.ChatID, err)
		} else {
			log.Println("Successfully sent user their morning message", user.ChatID)
		}
	}
}

func (s *Scheduler) CheckInForSunlight() {
	users, err := s.store.GetAllUsers()
	if err != nil {
		log.Println("Failed to access users from DB", err)
		return
	}

	for _, user := range users {
		msg := "sunlight checkin"
		_, err := s.bot.Send(tg.ChatID(user.ChatID), msg)
		if err != nil {
			log.Println("Failed to send user their sunlight check-in", user.ChatID, err)
		} else {
			log.Println("Successfully sent users their sunlight check-in", user.ChatID)
		}
	}
}
