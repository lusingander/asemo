package main

import (
	"time"

	"github.com/lusingander/asemo/asemo"
)

func (a *application) sendEmailHandler(req *asemo.SendEmailRequest) (*asemo.SendEmailResponse, *asemo.SendEmailError) {
	messageId, _ := generateMessageId()

	message := &message{
		messageId:        messageId,
		fromAddress:      req.FromEmailAddress,
		replyToAddresses: req.ReplyToAddresses,
		toAddresses:      req.Destination.ToAddresses,
		ccAddresses:      req.Destination.CcAddresses,
		bccAddresses:     req.Destination.BccAddresses,
		subject:          req.Content.Simple.Subject.Data,
		bodyHtml:         req.Content.Simple.Body.HTML.Data,
		bodyText:         req.Content.Simple.Body.Text.Data,
		receivedAt:       time.Now(),
	}

	a.messageRepository.set(messageId, message)

	resp := &asemo.SendEmailResponse{
		MessageId: message.messageId,
	}
	return resp, nil
}
