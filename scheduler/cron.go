package scheduler

// ToDo : Add a logger for each cron job, type of message/reminder. number of users sent to, time, and error if any
import (
	"agent-care-tg/storage"
	"log"

	"github.com/robfig/cron/v3"
	tg "gopkg.in/telebot.v3"
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
	s.cron.AddFunc("*/10 * * * * *", func() {
		log.Println("Scheduler Fired")
	})
	//Triggering Cron Jobs
	s.cron.AddFunc("0 */10 * * * *", func() {
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
	s.sendMessageToAllUsers("Morning Message", "Morning Message")
}

func (s *Scheduler) checkInForSunlight(localHour uint8) {
	s.sendMessageToAllUsers("Sunlight Check-In", "Sunlight Check-In Message")
}

func (s *Scheduler) checkInForHealthyMeal(localHour uint8) {
	s.sendMessageToAllUsers("Healthy Meal Check-In", "Did you have a nutritious meal today?")
}

func (s *Scheduler) checkInForPersonalGoal(localHour uint8) {
	s.sendMessageToAllUsers("Personal Goal Check-In", "Did u work on that personal goal today?")
}

func (s *Scheduler) checkInForExcercise(localHour uint8) {
	s.sendMessageToAllUsers("Excercise Check-In", "Did you get a chance to workout today?")
}

func (s *Scheduler) sendMessageToAllUsersInTimeZone() {
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

//func checkForTime() {
// get users time and then check against last_sent//
// if itswithin 60mins of last sent
// send a good job message
// add streak
// if not, mentioned failure and send a asupportive message
//}
