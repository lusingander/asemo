package asemo

import (
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/labstack/echo/v4"
)

type SendEmailHandler func(*SendEmailRequest) (*SendEmailResponse, *SendEmailError)

// https://docs.aws.amazon.com/ses/latest/APIReference-V2/API_SendEmail.html
func (s *Server) sendEmail(c echo.Context) error {
	var req SendEmailRequest
	err := c.Bind(&req)
	if err != nil {
		return sendEmailErrorResponse(c, BadRequestException)
	}
	resp, e := s.sendEmailHandler(&req)
	if e != nil {
		return sendEmailErrorResponse(c, e)
	}
	return sendEmailSuccessResponse(c, resp)
}

func sendEmailSuccessResponse(c echo.Context, resp *SendEmailResponse) error {
	return c.JSON(http.StatusOK, resp)
}

func sendEmailErrorResponse(c echo.Context, e *SendEmailError) error {
	c.Response().Header().Set("x-amzn-errortype", e.ErrorName)
	return c.JSON(e.StatusCode, e.SesErrorResponse)
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

type SesErrorResponse struct {
	Message string `json:"message"`
}

var defaultSendEmailCounter uint32

func defaultSendEmailHandler(req *SendEmailRequest) (*SendEmailResponse, *SendEmailError) {
	atomic.AddUint32(&defaultSendEmailCounter, 1)
	messageId := fmt.Sprintf("%v", defaultSendEmailCounter)
	resp := &SendEmailResponse{
		MessageId: messageId,
	}
	return resp, nil
}

type SendEmailError struct {
	StatusCode int
	ErrorName  string
	SesErrorResponse
}

var (
	// common errors
	// https://docs.aws.amazon.com/ses/latest/APIReference-V2/CommonErrors.html
	AccessDeniedException = &SendEmailError{
		StatusCode: 403,
		ErrorName:  "AccessDeniedException",
	}
	ExpiredTokenException = &SendEmailError{
		StatusCode: 403,
		ErrorName:  "ExpiredTokenException",
	}
	IncompleteSignature = &SendEmailError{
		StatusCode: 403,
		ErrorName:  "IncompleteSignature",
	}
	InternalFailure = &SendEmailError{
		StatusCode: 500,
		ErrorName:  "InternalFailure",
	}
	MalformedHttpRequestException = &SendEmailError{
		StatusCode: 400,
		ErrorName:  "MalformedHttpRequestException",
	}
	NotAuthorized = &SendEmailError{
		StatusCode: 401,
		ErrorName:  "NotAuthorized",
	}
	OptInRequired = &SendEmailError{
		StatusCode: 403,
		ErrorName:  "OptInRequired",
	}
	RequestAbortedException = &SendEmailError{
		StatusCode: 400,
		ErrorName:  "RequestAbortedException",
	}
	RequestEntityTooLargeException = &SendEmailError{
		StatusCode: 413,
		ErrorName:  "RequestEntityTooLargeException",
	}
	RequestExpired = &SendEmailError{
		StatusCode: 400,
		ErrorName:  "RequestExpired",
	}
	RequestTimeoutException = &SendEmailError{
		StatusCode: 408,
		ErrorName:  "RequestTimeoutException",
	}
	ServiceUnavailable = &SendEmailError{
		StatusCode: 503,
		ErrorName:  "ServiceUnavailable",
	}
	ThrottlingException = &SendEmailError{
		StatusCode: 400,
		ErrorName:  "ThrottlingException",
	}
	UnrecognizedClientException = &SendEmailError{
		StatusCode: 403,
		ErrorName:  "UnrecognizedClientException",
	}
	UnknownOperationException = &SendEmailError{
		StatusCode: 404,
		ErrorName:  "UnknownOperationException",
	}
	ValidationError = &SendEmailError{
		StatusCode: 400,
		ErrorName:  "ValidationError",
	}

	// special errors
	AccountSuspendedException = &SendEmailError{
		StatusCode: 400,
		ErrorName:  "AccountSuspendedException",
	}
	BadRequestException = &SendEmailError{
		StatusCode: 400,
		ErrorName:  "BadRequestException",
	}
	LimitExceededException = &SendEmailError{
		StatusCode: 400,
		ErrorName:  "LimitExceededException",
	}
	MailFromDomainNotVerifiedException = &SendEmailError{
		StatusCode: 400,
		ErrorName:  "MailFromDomainNotVerifiedException",
	}
	MessageRejected = &SendEmailError{
		StatusCode: 400,
		ErrorName:  "MessageRejected",
	}
	NotFoundException = &SendEmailError{
		StatusCode: 404,
		ErrorName:  "NotFoundException",
	}
	SendingPausedException = &SendEmailError{
		StatusCode: 400,
		ErrorName:  "SendingPausedException",
	}
	TooManyRequestsException = &SendEmailError{
		StatusCode: 429,
		ErrorName:  "TooManyRequestsException",
	}
)
