package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const (
	w2gbaseurl = "https://w2g.tv/rooms"
	w2glistner = "!w2g"
)

type W2GPayload struct {
	APIKey string `json:"w2g_api_key,omitempty"`
	Share  string `json:"share,omitempty"`
}

func handleW2G(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return

	}

	// If the message doesn't contain our listener just ignore
	if !strings.Contains(m.Content, w2glistner) {
		return
	}

	var payload W2GPayload

	// Let's parse the messages
	msgParts := strings.Split(m.Content, " ")

	// If the length is larger than one then the user
	// probably provided a link to their video
	// make sure the url is valid and add it to the payload
	if len(msgParts) > 1 {
		url, err := url.ParseRequestURI(msgParts[1])
		if err != nil {
			log.Println(err)
			s.ChannelMessageSend(m.ChannelID, "URL provided is malformed, please provided a valid URL.")
			return
		}

		payload.Share = url.String()
	}

	b, err := json.Marshal(payload)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Oops something went wrong preparing request for W2G. Try Again or contact a server admin")
		return
	}

	log.Println(string(b))
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/create.json", w2gbaseurl), bytes.NewBuffer(b))
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Oops something went wrong preparing request for W2G. Try Again or contact a server admin")
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := s.Client.Do(req)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Oops something went wrong preparing request for W2G. Try Again or contact a server admin")
		return
	}

	if resp.StatusCode != http.StatusOK {
		s.ChannelMessageSend(m.ChannelID, "Oops something went wrong preparing request for W2G. Try Again or contact a server admin")
		return
	}

	// Get response from W2G
	var body map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		s.ChannelMessageSend(m.ChannelID, "Oops something went wrong preparing request for W2G. Try Again or contact a server admin")
		return
	}

	// Grab the stream key from the body
	streamkey, ok := body["streamkey"]
	if !ok {
		s.ChannelMessageSend(m.ChannelID, "Oops something went wrong preparing request for W2G. Try Again or contact a server admin")
		return
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s/%s", w2gbaseurl, streamkey))
}
