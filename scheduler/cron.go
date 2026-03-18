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
	s.sendMessageToAllUsers("Morning Message", "Morning Message")
}

func (s *Scheduler) CheckInForSunlight() {
	s.sendMessageToAllUsers("Sunlight Check-In", "Sunlight Check-In Message")
}

func (s *Scheduler) sendMessageToAllUsers(jobName string, msg string) error {
	users, err := s.store.GetAllUsers()
	if err != nil {
		log.Println("Failed to access users from DB for : ", jobName)
		log.Println("Erroe : ", err)
		return err
	}

	for _, user := range users {
		//ToDo: Add a filtering system for message for safety
		_, err := s.bot.Send(tg.ChatID(user.ChatID), msg)
		// ToDo: handle partialfails
		if err != nil {
			log.Println("Failed to send users message for : ", jobName)
			// ToDO: Clean up,this is here just for initial stage. Any sort personal info should needs to be removed post testing
			log.Println("User : ", user.ChatID)
			log.Println("Error : ", err)
			return err
		} else {
			log.Println("Successfully sent message to all users for job : ", jobName)
			return nil
		}
	}
	return nil
}
