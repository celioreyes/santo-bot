package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/celioreyes/santo-bot/pkg/howlongtobeat"
	"github.com/robfig/cron/v3"
)

const (
	generalchannel = "699841515610177568"
)

// Variables used for command line parameters
var (
	token string

	hltb *howlongtobeat.Service
)

// TODO: change logger to logrus and add envconf support
func main() {
	token = os.Getenv("DS_TOKEN")
	if token == "" {
		log.Println("DS_TOKEN is required.")
		return
	}

	hltb = howlongtobeat.New(http.DefaultClient)

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Println("error creating Discord session,", err)
		return

	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(handleW2G)
	dg.AddHandler(handleHLTB)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)

	// Open a websocket connection to Discord and begin listening.
	if err := dg.Open(); err != nil {
		log.Println("error opening connection,", err)
		return
	}

	// configure crons
	croninstance := cron.New()

	// On Fridays send itsFriday video at 9am
	croninstance.AddFunc("0 13 * * 5", itsFriday(dg))

	// Start cron
	croninstance.Start()

	// Wait here until CTRL-C or other term signal is received.
	log.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()

	// Cleanly stop all crons
	croninstance.Stop()
}
