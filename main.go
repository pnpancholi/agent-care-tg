package main

import (
	bot "agent-care-tg/bot"
	"agent-care-tg/scheduler"
	"agent-care-tg/storage"
	"log/slog"
	"os"
	"time"

	env "github.com/joho/godotenv"
	tg "gopkg.in/telebot.v3"
)

func main() {
	if _, err := os.Stat(".env"); err == nil {
		env.Load()
	}

	token := os.Getenv("TG_BOT_TOKEN")
	if token == "" {
		slog.Error("[agent-care-tg]: TG_BOT_TOKEN is not set")
		os.Exit(1)
	}

	agentBot, err := tg.NewBot(tg.Settings{
		Token:  token,
		Poller: &tg.LongPoller{Timeout: 1 * time.Second},
	})
	if err != nil {
		slog.Error("[agent-care-tg]: Failed to create bot", err)
		os.Exit(1)
	}

	db := storage.Connect()
	defer db.Close()
	store := storage.NewStore(db)
	handler := bot.NewHandler(agentBot, store)
	handler.Register()
	s := scheduler.New(store, agentBot)
	s.Start()
	slog.Info("[agent-care-tg]: Authorized on account", agentBot.Me.Username)
	slog.Info("[agent-care-tg]: Bot is online")
	agentBot.Start()

}
