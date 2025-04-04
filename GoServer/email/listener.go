package email

import (
	"fmt"
	idle "github.com/emersion/go-imap-idle"
	"github.com/emersion/go-imap/client"
	"goServer/config"
	"goServer/database"
	"goServer/websocket"
	"log"
	"time"
)

func StartEmailListener() {
	// Connect to the server
	c, err := client.DialTLS(config.IMAPServer, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Logout()

	// Login
	if err := c.Login(config.IMUsername, config.IMPassword); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Logged in")

	// Select INBOX
	_, err = c.Select("INBOX", false)
	if err != nil {
		log.Fatal(err)
	}

	idleClient := idle.NewClient(c)

	// Updates channel
	updates := make(chan client.Update)
	c.Updates = updates

	for {
		mboxUpd, err := waitForMailboxUpdate(idleClient, updates)
		if err != nil {
			log.Fatal(err)
		}
		if meta, err := fetchParseEmail(c, mboxUpd.Mailbox); err == nil {
			log.Println("Received new email:", meta.Subject)
			err := database.InsertEmail(meta)
			if err != nil {
				log.Println("Failed to insert into DB:", err)
			} else {
				go triggerMLPrediction(meta.MessageID)
				go func() {
					time.Sleep(10 * time.Second)
					websocket.BroadcastStats()
				}()
			}
		}
	}
}

func waitForMailboxUpdate(c *idle.Client, updates chan client.Update) (*client.MailboxUpdate, error) {
	done := make(chan error, 1)
	stop := make(chan struct{})

	// Start IDLE mode in a goroutine
	go func() {
		done <- c.IdleWithFallback(stop, 5*time.Minute)
	}()

	var mboxUpd *client.MailboxUpdate
waitLoop:
	for {
		select {
		case upd := <-updates:
			if mboxUpd = asMailboxUpdate(upd); mboxUpd != nil {
				break waitLoop
			}
		case err := <-done:
			if err != nil {
				return nil, fmt.Errorf("error while idling: %s", err.Error())
			}
			return nil, nil // timeout without updates
		}
	}

	close(stop)
	<-done

	return mboxUpd, nil
}

func asMailboxUpdate(upd client.Update) *client.MailboxUpdate {
	if v, ok := upd.(*client.MailboxUpdate); ok {
		return v
	}
	return nil
}
