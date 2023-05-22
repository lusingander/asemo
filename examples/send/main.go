package main

import (
	"context"
	"fmt"
	"math/rand"

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

	server := asemo.NewServer()

	server.SetPort(8081)
	server.SetSendEmailHandler(
		func(asemo.SendEmailRequest) asemo.SendEmailResponse {
			n := rand.Intn(1000000)
			messageId := fmt.Sprintf("%06d", n)
			return asemo.SendEmailResponse{
				MessageId: messageId,
			}
		},
	)

	go server.Start()

	ctx := context.Background()
	client, err := setupClient(ctx)
	if err != nil {
		return err
	}

	fmt.Println("Press 'q' to quit...")
	for {
		if s, err := scan(); err != nil {
			return err
		} else if s == "q" {
			fmt.Println("quit")
			return nil
		}

		messageId, err := callSendEmail(ctx, client)
		if err != nil {
			return err
		}
		fmt.Printf("success: messageId = %v\n", messageId)
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
			return aws.Endpoint{URL: "http://localhost:8081"}, nil
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

func callSendEmail(ctx context.Context, client *sesv2.Client) (string, error) {
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
