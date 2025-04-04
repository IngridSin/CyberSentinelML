package email

import (
	"errors"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
	"goServer/models"
	"io"
	"log"
	"strings"
	"time"
)

func fetchParseEmail(c *client.Client, mbox *imap.MailboxStatus) (*models.EmailMetadata, error) {
	seqSet := new(imap.SeqSet)
	seqSet.AddNum(mbox.Messages)

	ch := make(chan *imap.Message, 1)
	err := c.Fetch(seqSet, []imap.FetchItem{"BODY.PEEK[]"}, ch)
	if err != nil {
		return nil, err
	}

	msg := <-ch
	if msg == nil {
		return nil, errors.New("no message returned")
	}

	var meta models.EmailMetadata

	for _, literal := range msg.Body {
		mr, err := mail.CreateReader(literal)
		if err != nil {
			return nil, err
		}

		header := mr.Header

		// Subject
		meta.Subject, _ = header.Subject()

		// From
		if fromList, err := header.AddressList("From"); err == nil && len(fromList) > 0 {
			meta.From = fromList[0].String()
		}

		// To
		if toList, err := header.AddressList("To"); err == nil && len(toList) > 0 {
			meta.To = toList[0].String()
		}

		// Message-ID
		if msgID, err := header.Text("Message-Id"); err == nil {
			meta.MessageID = strings.Trim(msgID, "<>")
		}

		// Date
		if date, err := header.Date(); err == nil {
			meta.Date = date.In(time.Now().Location())
		}

		// DKIM
		if dkim, err := header.Text("DKIM-Signature"); err == nil {
			meta.DKIMSignature = dkim
		}

		// SPF
		if spf, err := header.Text("Received-SPF"); err == nil {
			meta.SPFResult = spf
		}

		// Delivered-To
		if deliveredTo, err := header.Text("Delivered-To"); err == nil {
			meta.DeliveredTo = deliveredTo
		}

		// Return-Path
		if returnPath, err := header.Text("Return-Path"); err == nil {
			meta.ReturnPath = returnPath
		}

		// Received (last hop)
		if received, err := header.Text("Received"); err == nil {
			meta.Received = received
		}

		// Body & Attachments
		for {
			p, err := mr.NextPart()
			if err == io.EOF {
				break
			} else if err != nil {
				log.Println("Error reading part:", err)
				break
			}

			switch h := p.Header.(type) {
			case *mail.InlineHeader:
				b, _ := io.ReadAll(p.Body)
				meta.Body = string(b)

			case *mail.AttachmentHeader:
				filename, _ := h.Filename()
				meta.Attachments = append(meta.Attachments, filename)
			}
		}

		break
	}

	return &meta, nil
}
