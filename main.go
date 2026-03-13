package main

import (
	bot "agent-care-tg/bot"
	env "github.com/joho/godotenv"
	tg "gopkg.in/telebot.v3"
	"log"
	"os"
	"time"
)

func main() {
	err := env.Load()
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

	handler := bot.NewHandler(agentBot)
	handler.Register()
	log.Printf("[agent-care-tg]: Authorized on account %s, bot is online", agentBot.Me.Username)
	agentBot.Start()

}
