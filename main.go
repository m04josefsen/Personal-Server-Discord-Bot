package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"net/http"
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
	fmt.Println("in handle files")

	subjectToChannel := map[string]string{
		"adse1310": "1",
		"adse2100": "1313564569058279506",
		"dafe2200": "2",
		"dats2300": "3",
	}

	for _, attachment := range m.Attachments {

		// filename -> lowercase
		lowerFilename := strings.ToLower(attachment.Filename)
		fmt.Println(lowerFilename)

		if strings.HasPrefix(lowerFilename, "adse") || strings.HasPrefix(lowerFilename, "dafe") || strings.HasPrefix(lowerFilename, "dats") {
			channelID := lowerFilename[:8]
			fmt.Println("Extracted channelID:", channelID)

			actualChannelID, exists := subjectToChannel[channelID]
			fmt.Println("Actual channelID:", actualChannelID)

			if !exists {
				fmt.Println("Invalid channelID:", channelID)
				continue
			}

			if strings.HasSuffix(lowerFilename, ".pdf") {

				// Forward the PDF
				resp, err := http.Get(attachment.URL) // Download the PDF
				if err != nil {
					log.Println("Error downloading PDF:", err)
					continue
				}
				defer resp.Body.Close()

				// Send the PDF to the target channel
				_, err = s.ChannelFileSend(actualChannelID, attachment.Filename, resp.Body) // Corrected here
				if err != nil {
					log.Println("Error forwarding PDF:", err)
				} else {
					log.Println("PDF forwarded successfully:", attachment.Filename)
				}

				fmt.Println("At the end of .pdf if")
			}
		}
	}
}
