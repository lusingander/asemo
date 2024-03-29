package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// fixme: pagination?
func (a *application) listMessageHandler(c echo.Context) error {
	ms := a.messageRepository.getAll()

	resp := &listMessageResponse{
		Messages: make([]*listMessageResponseItem, len(ms)),
	}
	for i, m := range ms {
		receivedAt := m.receivedAt.Local().Format(time.RFC3339)
		resp.Messages[i] = &listMessageResponseItem{
			MessageId:  m.messageId,
			From:       m.fromAddress,
			Subject:    m.subject,
			ReceivedAt: receivedAt,
		}
	}
	return c.JSON(http.StatusOK, resp)
}

type listMessageResponse struct {
	Messages []*listMessageResponseItem `json:"messages"`
}

type listMessageResponseItem struct {
	MessageId  string `json:"message_id"`
	From       string `json:"from"`
	Subject    string `json:"subject"`
	ReceivedAt string `json:"received_at"`
}

func (a *application) getMessageHandler(c echo.Context) error {
	messageId := c.Param("id")

	message := a.messageRepository.get(messageId)
	if message == nil {
		errorMsg := fmt.Sprintf("message not found (id = %s)", messageId)
		resp := &errorResponse{
			Message: errorMsg,
		}
		log.Println(errorMsg)
		return c.JSON(http.StatusNotFound, resp)
	}

	body := message.bodyHtml
	if body == "" {
		body = message.bodyText
	}

	receivedAt := message.receivedAt.Local().Format(time.RFC3339)

	resp := &getMessageResponse{
		MessageId:  messageId,
		From:       message.fromAddress,
		ReplyTo:    message.replyToAddresses,
		To:         message.toAddresses,
		Cc:         message.ccAddresses,
		Bcc:        message.bccAddresses,
		Subject:    message.subject,
		Body:       body,
		ReceivedAt: receivedAt,
	}
	return c.JSON(http.StatusOK, resp)
}

type getMessageResponse struct {
	MessageId  string   `json:"message_id"`
	From       string   `json:"from"`
	ReplyTo    []string `json:"reply_to"`
	To         []string `json:"to"`
	Cc         []string `json:"cc"`
	Bcc        []string `json:"bcc"`
	Subject    string   `json:"subject"`
	Body       string   `json:"body"`
	ReceivedAt string   `json:"received_at"`
}

type errorResponse struct {
	Message string `json:"message"`
}
