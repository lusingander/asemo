package asemo

import (
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/labstack/echo/v4"
)

type SendEmailHandler func(SendEmailRequest) SendEmailResponse

// https://docs.aws.amazon.com/ses/latest/APIReference-V2/API_SendEmail.html
func (s *Server) sendEmail(c echo.Context) error {
	var req SendEmailRequest
	err := c.Bind(&req)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request") // fixme
	}
	resp := s.sendEmailHandler(req)
	return c.JSON(http.StatusOK, resp)
}

type SendEmailRequest struct {
	ConfigurationSetName string `json:"ConfigurationSetName"`
	Content              struct {
		Raw struct {
			Data string `json:"Data"`
		} `json:"Raw"`
		Simple struct {
			Body struct {
				HTML struct {
					Charset string `json:"Charset"`
					Data    string `json:"Data"`
				} `json:"Html"`
				Text struct {
					Charset string `json:"Charset"`
					Data    string `json:"Data"`
				} `json:"Text"`
			} `json:"Body"`
			Subject struct {
				Charset string `json:"Charset"`
				Data    string `json:"Data"`
			} `json:"Subject"`
		} `json:"Simple"`
		Template struct {
			TemplateArn  string `json:"TemplateArn"`
			TemplateData string `json:"TemplateData"`
			TemplateName string `json:"TemplateName"`
		} `json:"Template"`
	} `json:"Content"`
	Destination struct {
		BccAddresses []string `json:"BccAddresses"`
		CcAddresses  []string `json:"CcAddresses"`
		ToAddresses  []string `json:"ToAddresses"`
	} `json:"Destination"`
	EmailTags []struct {
		Name  string `json:"Name"`
		Value string `json:"Value"`
	} `json:"EmailTags"`
	FeedbackForwardingEmailAddress            string `json:"FeedbackForwardingEmailAddress"`
	FeedbackForwardingEmailAddressIdentityArn string `json:"FeedbackForwardingEmailAddressIdentityArn"`
	FromEmailAddress                          string `json:"FromEmailAddress"`
	FromEmailAddressIdentityArn               string `json:"FromEmailAddressIdentityArn"`
	ListManagementOptions                     struct {
		ContactListName string `json:"ContactListName"`
		TopicName       string `json:"TopicName"`
	} `json:"ListManagementOptions"`
	ReplyToAddresses []string `json:"ReplyToAddresses"`
}

type SendEmailResponse struct {
	MessageId string `json:"MessageId"`
}

var defaultSendEmailCounter uint32

func defaultSendEmailHandler(req SendEmailRequest) SendEmailResponse {
	atomic.AddUint32(&defaultSendEmailCounter, 1)
	messageId := fmt.Sprintf("%v", defaultSendEmailCounter)
	return SendEmailResponse{
		MessageId: messageId,
	}
}
