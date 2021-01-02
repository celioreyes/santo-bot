package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const hltblistener = "!hltb"

func handleHLTB(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	// If the message doesn't contain our listener just ignore
	if !strings.Contains(m.Content, hltblistener) {
		return
	}

	// Remove the listener and use the remaining text as the game to search for
	gameInputted := strings.Replace(m.Content, hltblistener+" ", "", -1)

	if gameInputted == "" || gameInputted == hltblistener {
		s.ChannelMessageSend(m.ChannelID, "Please provide a game title to search for.")
		return
	}

	// Looks up the game in HLTB
	games, err := hltb.SearchGame(gameInputted)
	if err != nil {
		log.Printf("Error: %s", err.Error())
		s.ChannelMessageSend(m.ChannelID, "There was an error looking game up in HowLongToBeat, try again or contact admin.")
		return
	}

	if len(games) == 0 {
		s.ChannelMessageSend(m.ChannelID, "No games found in HowLongToBeat database")
		return
	}

	msg := "```md"

	for _, game := range games {
		timeString := ""
		for _, time := range game.Times {
			timeString += fmt.Sprintf("* %s: %s\n", time.Type, time.Value)
		}

		gameString := fmt.Sprintf("\n# %s:\n%s", game.Title, timeString)
		msg += gameString
	}

	s.ChannelMessageSend(m.ChannelID, msg+"```")
	return
}
