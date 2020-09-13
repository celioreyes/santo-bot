package main

import (
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

// itsFriday - go routine that will keep checking until it's friday. Once it's Friday will post desidred video
func itsFriday(s *discordgo.Session) func() {
	return func() {
		log.Println("Friday!")
		videoFile := "its_friday_mufasa.mp4"

		file, err := os.Open("./media/" + videoFile)
		if err != nil {
			log.Println(err.Error())
			return
		}

		// log.Printf("Looks good? %s", msg.ID)
		if _, err := s.ChannelFileSend(generalchannel, videoFile, file); err != nil {
			log.Println("failed to send it's friday video")
		}
	}
}
