[![Go Reference](https://pkg.go.dev/badge/github.com/lusingander/asemo.svg)](https://pkg.go.dev/github.com/lusingander/asemo)

# asemo

Aws SEs MOck

## About

asemo provides [Amazon SES API(v2)](https://docs.aws.amazon.com/ses/latest/APIReference-V2/Welcome.html) mock library and simple standalone command

## Usage

require Go 1.20+

### As standalone command

#### Installation

`$ go install github.com/lusingander/asemo/cmd/asemo@latest`

#### Run

`$ asemo -port 8080`

Then the mock server will start at http://localhost:8080.

#### API

Sent messages can be referenced through the API.

```
$ curl -s http://localhost:8080/api/messages | jq .
{
  "messages": [
    {
      "message_id": "2653e13d-9144-4627-9806-494c90ecc3b4",
      "from": "from@example.com",
      "subject": "foo",
      "received_at": "2023-12-18T08:24:43+09:00"
    },
    {
      "message_id": "092f09ba-bb61-4548-beb0-b9fe9a60c344",
      "from": "from@example.com",
      "subject": "bar",
      "received_at": "2023-12-18T08:24:44+09:00"
    }
  ]
}

$ curl -s http://localhost:8080/api/messages/6d2aed72-62d8-4839-8ee1-19c06751077d | jq .
{
  "message_id": "6d2aed72-62d8-4839-8ee1-19c06751077d",
  "from": "from@example.com",
  "reply_to": null,
  "to": [
    "to@example.com"
  ],
  "cc": null,
  "bcc": null,
  "subject": "hi",
  "body": "hello foo",
  "received_at": "2023-05-24T22:41:42+09:00"
}
```

### As library

```go
server := asemo.NewServer()

server.SetSendEmailHandler(
	func(req *asemo.SendEmailRequest) (*asemo.SendEmailResponse, *asemo.SendEmailError) {
		fmt.Printf("receive: [subject = '%v']\n", req.Content.Simple.Subject.Data)
		return &asemo.SendEmailResponse{
			MessageId: "1",
		}, nil
	},
)

server.Start() // mock server starts on localhost:8080
```

See [./examples](./examples) for more details.

### Supported actions

- [SendEmail](https://docs.aws.amazon.com/ses/latest/APIReference-V2/API_SendEmail.html)

## License

MIT
