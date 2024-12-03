package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	token := "your-bot-token"
	sess, err := initializeBot(token)
	if err != nil {
		log.Fatal(err)
	}
	defer sess.Close()

	sess.AddHandler(handleMessage)

	err = sess.Open()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Bot is online!")

	// Wait for termination signal
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	fmt.Println("Bot shutting down.")
}

func initializeBot(token string) (*discordgo.Session, error) {
	sess, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	sess.Identify.Intents = discordgo.IntentsAllWithoutPrivileged
	return sess, nil
}

func handleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if len(m.Attachments) > 0 {
		handleFiles(s, m)
	} else if m.Content == "Hello" {
		s.ChannelMessageSend(m.ChannelID, "Hello, World!")
	}
}

func handleFiles(s *discordgo.Session, m *discordgo.MessageCreate) {
	for _, attachment := range m.Attachments {
		// Check if the filename contains a valid prefix
		if strings.HasPrefix(attachment.Filename, "adse") || strings.HasPrefix(attachment.Filename, "dafe") || strings.HasPrefix(attachment.Filename, "dats") {
			// First 8 characters are the channel ID, i.e subjectID
			channelID := attachment.Filename[:8]
			s.ChannelMessageSend(channelID, m.Content)
		}
	}

}
