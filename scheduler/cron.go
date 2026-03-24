package scheduler

// ToDo : Add a logger for each cron job, type of message/reminder. number of users sent to, time, and error if any
import (
	bot "agent-care-tg/bot"
	"agent-care-tg/storage"
	"fmt"
	"github.com/robfig/cron/v3"
	tg "gopkg.in/telebot.v3"
	"log"
	"time"
)

type Scheduler struct {
	cron  *cron.Cron
	store *storage.Store
	bot   *tg.Bot
}

func New(store *storage.Store, bot *tg.Bot) *Scheduler {
	return &Scheduler{
		cron:  cron.New(),
		store: store,
		bot:   bot,
	}
}

func (s *Scheduler) Start() {
	s.cron.AddFunc("*/10 * * * *", func() {
		log.Println("Scheduler Triggered...")
	})
	//Triggering Cron Jobs
	s.cron.AddFunc("*/10 * * * *", func() {
		s.sendMorningMessage(7)
		s.checkInForSunlight(14)
		s.checkInForHealthyMeal(14)
		s.checkInForPersonalGoal(21)
		s.checkInForExcercise(17)
	})

	s.cron.Start()
	log.Println("Scheduler Started...")
}

func (s *Scheduler) Stop() {
	s.cron.AddFunc("", func() {
		log.Println("Scheduler Stopped...")
	})
}

func (s *Scheduler) sendMorningMessage(localHour uint8) {
	s.sendMessageToAllUsersInTimeZone(localHour, bot.MsgMorningCheckIn)
}

func (s *Scheduler) checkInForSunlight(localHour uint8) {
	s.sendMessageToAllUsersInTimeZone(localHour, bot.MsgSunlightCheckIn)
}

func (s *Scheduler) checkInForHealthyMeal(localHour uint8) {
	s.sendMessageToAllUsersInTimeZone(localHour, bot.MsgMealCheckIn)
}

func (s *Scheduler) checkInForPersonalGoal(localHour uint8) {
	s.sendMessageToAllUsersInTimeZone(localHour, bot.MsgPersonalGoalCheckIn)
}

func (s *Scheduler) checkInForExcercise(localHour uint8) {
	s.sendMessageToAllUsersInTimeZone(localHour, bot.MsgExcerciseCheckIn)
}

func (s *Scheduler) sendMessageToAllUsersInTimeZone(hour uint8, msg string) {
	// check if last_sent was within last 5 mins//

	// timezone match//
	users, err := s.store.GetAllUsers()
	if err != nil {
		log.Println("Failed to access users from DB for sendMessageToAllUsersInTimeZone : ", err)
		return
	}

	for _, user := range users {
		loc, err := time.LoadLocation(user.Timezone)
		if err != nil {
			log.Println("Failed to convert timezone from DB object : ", err)
		}

		localTime := time.Now().In(loc)

		if localTime.Hour() == int(hour) && localTime.Minute() < 10 {
			formattedMsg := fmt.Sprintf(msg, user.Username)
			_, err := s.bot.Send(tg.ChatID(user.ChatID), formattedMsg, tg.ModeMarkdown)
			if err != nil {
				log.Println("Failed to send message to : ", user.TGUsername)
			}
		}
	}
	// write to last_sent
}

func (s *Scheduler) sendMessageToAllUsers(jobName string, msg string) error {
	users, err := s.store.GetAllUsers()
	if err != nil {
		log.Println("Failed to access users from DB for : ", jobName)
		log.Println("Erroe : ", err)
		return err
	}

	for _, user := range users {
		// ToDo: Make msg make context for yes or no response
		//ToDo: Add a filtering system for message for safety
		// Rendering buttons for each task with call back//
		markup := &tg.ReplyMarkup{}
		taskDoneBtnKey := jobName + "_task_done"
		taskSkippedBtnKey := jobName + "_task_skipped"
		taskDoneBtn := markup.Data("Yes", taskDoneBtnKey)
		taskSkippedBtn := markup.Data("Skipped", taskSkippedBtnKey)
		markup.Inline(markup.Row(taskDoneBtn, taskSkippedBtn))

		_, err := s.bot.Send(tg.ChatID(user.ChatID), msg, markup)
		// ToDo: handle partialfails for sending message
		if err != nil {
			log.Println("Failed to send users message for : ", jobName)
			// ToDO: Clean up,this is here just for initial stage. Any sort personal info should needs to be removed post testing
			log.Println("User : ", user.ChatID)
			log.Println("Error : ", err)
			// Continue to next user on error
		} else {
			log.Println("Successfully sent message to user ", user.ChatID, " for job : ", jobName)
		}
	}
	return nil
}
