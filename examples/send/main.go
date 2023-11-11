package main

import (
	"context"
	"fmt"
	"math/rand"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	"github.com/lusingander/asemo/asemo"
)

const (
	port = 8081
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {

	// mock server
	server := asemo.NewServer()

	server.E.HidePort = true
	server.SetPort(port)
	server.SetSendEmailHandler(
		func(req *asemo.SendEmailRequest) (*asemo.SendEmailResponse, *asemo.SendEmailError) {
			sub := req.Content.Simple.Subject.Data
			body := req.Content.Simple.Body.Text.Data
			fmt.Printf("[server] received request: subject = '%v', body = '%v'\n", sub, body)

			// return the error information if the body contains an error name
			if e := findError(body); e != nil {
				return nil, e
			}

			// generate random message id
			n := rand.Intn(1000000)
			messageId := fmt.Sprintf("%06d", n)
			resp := &asemo.SendEmailResponse{
				MessageId: messageId,
			}
			return resp, nil
		},
	)

	go server.Start()

	// sender
	ctx := context.Background()
	client, err := setupClient(ctx)
	if err != nil {
		return err
	}

	fmt.Println("Press 'q' to quit...")
	for {
		fmt.Printf(">> ")
		s, err := scan()
		if err != nil {
			return err
		}
		if s == "q" {
			fmt.Println("quit")
			return nil
		}

		if messageId, err := callSendEmail(ctx, client, s); err == nil {
			fmt.Printf("[sender] success: messageId = %v\n", messageId)
		} else {
			fmt.Printf("[sender] error: %v\n", err)
		}
	}
}

func scan() (string, error) {
	var s string
	if _, err := fmt.Scan(&s); err != nil {
		return "", err
	}
	return s, nil
}

func setupClient(ctx context.Context) (*sesv2.Client, error) {
	endpointResolver := config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
		func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{URL: fmt.Sprintf("http://localhost:%d", port)}, nil
		},
	))
	credentialsProvider := config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("dummy", "dummy", "dummy"))
	cfg, err := config.LoadDefaultConfig(ctx, endpointResolver, credentialsProvider)
	if err != nil {
		return nil, err
	}
	client := sesv2.NewFromConfig(cfg)
	return client, nil
}

func callSendEmail(ctx context.Context, client *sesv2.Client, s string) (string, error) {
	input := &sesv2.SendEmailInput{
		FromEmailAddress: ptr("from@example.com"),
		Destination: &types.Destination{
			ToAddresses: []string{"to@example.com"},
		},
		Content: &types.EmailContent{
			Simple: &types.Message{
				Body: &types.Body{
					Text: &types.Content{
						Data: ptr(fmt.Sprintf("hello %v", s)),
					},
				},
				Subject: &types.Content{
					Data: ptr("hi"),
				},
			},
		},
	}
	res, err := client.SendEmail(ctx, input)
	if err != nil {
		return "", err
	}
	return *res.MessageId, nil
}

func findError(s string) *asemo.SendEmailError {
	m := map[string]*asemo.SendEmailError{
		"AccessDeniedException":              &asemo.AccessDeniedException,
		"ExpiredTokenException":              &asemo.ExpiredTokenException,
		"IncompleteSignature":                &asemo.IncompleteSignature,
		"InternalFailure":                    &asemo.InternalFailure,
		"MalformedHttpRequestException":      &asemo.MalformedHttpRequestException,
		"NotAuthorized":                      &asemo.NotAuthorized,
		"OptInRequired":                      &asemo.OptInRequired,
		"RequestAbortedException":            &asemo.RequestAbortedException,
		"RequestEntityTooLargeException":     &asemo.RequestEntityTooLargeException,
		"RequestExpired":                     &asemo.RequestExpired,
		"RequestTimeoutException":            &asemo.RequestTimeoutException,
		"ServiceUnavailable":                 &asemo.ServiceUnavailable,
		"ThrottlingException":                &asemo.ThrottlingException,
		"UnrecognizedClientException":        &asemo.UnrecognizedClientException,
		"UnknownOperationException":          &asemo.UnknownOperationException,
		"ValidationError":                    &asemo.ValidationError,
		"AccountSuspendedException":          &asemo.AccountSuspendedException,
		"BadRequestException":                &asemo.BadRequestException,
		"LimitExceededException":             &asemo.LimitExceededException,
		"MailFromDomainNotVerifiedException": &asemo.MailFromDomainNotVerifiedException,
		"MessageRejected":                    &asemo.MessageRejected,
		"NotFoundException":                  &asemo.NotFoundException,
		"SendingPausedException":             &asemo.SendingPausedException,
		"TooManyRequestsException":           &asemo.TooManyRequestsException,
	}
	for k, v := range m {
		if strings.Contains(s, k) {
			return v
		}
	}
	return nil
}

func ptr[T any](v T) *T {
	return &v
}
