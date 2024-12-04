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
	"time"
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

	// go makes it run concurrently with the rest of the code
	go sendPeriodicMessages(sess)

	// Wait for termination signal
	// TODO: hva gjør den?
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
	// If the bot sends a message
	if m.Author.ID == s.State.User.ID {
		return
	}

	if len(m.Attachments) > 0 {
		handleFiles(s, m)
	}
}

func handleFiles(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println("in handle files")

	subjectToChannel := map[string]string{
		// First year
		// First semester
		"data1100": "1313886461480927242",
		"dape1300": "1313886404949966898",
		"dape1400": "1313886439766884374",

		// Second semester
		"adse1310": "1313564389097603153",
		"data1500": "1313886488991367188",
		"data1700": "1313886509652377620",

		// Second year
		// Third semester
		"adse2100": "1313564569058279506",
		"dafe2200": "1313564534132310047",
		"dats2300": "1313564591740944394",
	}

	for _, attachment := range m.Attachments {

		// filename -> lowercase
		lowerFilename := strings.ToLower(attachment.Filename)
		fmt.Println(lowerFilename)

		// TODO: må kunne forkorte denne
		// For Documents under OsloMet category
		if strings.HasPrefix(lowerFilename, "adse") || strings.HasPrefix(lowerFilename, "dafe") || strings.HasPrefix(lowerFilename, "dats") || strings.HasPrefix(lowerFilename, "data") || strings.HasPrefix(lowerFilename, "dape") {
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
				resp, err := http.Get(attachment.URL)
				if err != nil {
					log.Println("Error downloading PDF:", err)
					continue
				}
				defer resp.Body.Close()

				// Send the PDF to the target channel
				_, err = s.ChannelFileSend(actualChannelID, attachment.Filename, resp.Body)
				if err != nil {
					log.Println("Error forwarding PDF:", err)
				} else {
					log.Println("PDF forwarded successfully:", attachment.Filename)
				}
			}
		}
	}
}

func sendPeriodicMessages(s *discordgo.Session) {
	channelID := "1313564463579791450"

	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		_, err := s.ChannelMessageSend(channelID, "OH MAH GAH!!")
		if err != nil {
			log.Println("Failed to send periodic message:", err)
		}
	}
}
