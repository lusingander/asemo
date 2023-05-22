package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	"github.com/lusingander/asemo/asemo"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {

	go asemo.Run()

	messageId, err := callSendEmail()
	if err != nil {
		return err
	}
	fmt.Printf("success: messageId = %v\n", messageId)

	return nil
}

func callSendEmail() (string, error) {
	ctx := context.Background()
	endpointResolver := config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
		func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{URL: "http://localhost:8080"}, nil
		},
	))
	credentialsProvider := config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("dummy", "dummy", "dummy"))
	cfg, err := config.LoadDefaultConfig(ctx, endpointResolver, credentialsProvider)
	if err != nil {
		return "", err
	}
	client := sesv2.NewFromConfig(cfg)
	input := &sesv2.SendEmailInput{
		FromEmailAddress: ptr("from@example.com"),
		Destination: &types.Destination{
			ToAddresses: []string{"to@example.com"},
		},
		Content: &types.EmailContent{
			Simple: &types.Message{
				Body: &types.Body{
					Text: &types.Content{
						Data: ptr("hello ses"),
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

func ptr[T any](v T) *T {
	return &v
}
