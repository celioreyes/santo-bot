package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
var (
	token string
)

func main() {
	token = os.Getenv("DS_TOKEN")
	if token == "" {
		log.Println("DS_TOKEN is required.")
		return
	}
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Println("error creating Discord session,", err)
		return

	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(handleW2G)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)

	// Open a websocket connection to Discord and begin listening.
	if err := dg.Open(); err != nil {
		log.Println("error opening connection,", err)
		return
	}

	// Ideally we want this to be inside of some loop or cron / channel to check for day of week
	// msg, err := dg.ChannelMessageSend("753977130714660963", "It's saturdayyyyy")
	// if err != nil {
	//	log.Println("fuckkk")
	//	log.Println(err)
	//	return
	//}

	// log.Printf("Looks good? %s", msg.ID)

	// Wait here until CTRL-C or other term signal is received.
	log.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}
