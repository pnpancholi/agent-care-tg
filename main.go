package main

import (
	bot "agent-care-tg/bot"
	"agent-care-tg/scheduler"
	"agent-care-tg/storage"
	"log"
	"os"
	"time"

	env "github.com/joho/godotenv"
	tg "gopkg.in/telebot.v3"
)

func main() {
	_, err := os.Stat(".env"); err == nil {
		env.Load()
	}
	if err != nil {
		log.Fatal("[agent-care-tg]: Failed to load environment variables")
	}

	token := os.Getenv("TG_BOT_TOKEN")
	if token == "" {
		log.Fatalf("[agent-care-tg]: TG_BOT_TOKEN is not set")
	}

	agentBot, err := tg.NewBot(tg.Settings{
		Token:  token,
		Poller: &tg.LongPoller{Timeout: 1 * time.Second},
	})
	if err != nil {
		log.Fatalf("[agent-care-tg]: Failed to create bot: %v", err)
	}

	db := storage.Connect()
	defer db.Close()
	store := storage.NewStore(db)
	handler := bot.NewHandler(agentBot, store)
	handler.Register()
	s := scheduler.New(store, agentBot)
	s.Start()
	log.Printf("[agent-care-tg]: Authorized on account %s, bot is online", agentBot.Me.Username)
	agentBot.Start()

}
