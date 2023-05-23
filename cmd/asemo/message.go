package main

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

type message struct {
	messageId string

	fromAddress string

	replyToAddresses []string
	toAddresses      []string
	ccAddresses      []string
	bccAddresses     []string

	subject  string
	bodyHtml string
	bodyText string

	receivedAt time.Time
}

func generateMessageId() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}

type messageRepository struct {
	m  map[string]*message
	mu sync.RWMutex
}

func newMessageRepository() *messageRepository {
	return &messageRepository{
		m: make(map[string]*message),
	}
}

func (r *messageRepository) set(messageId string, messsage *message) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.m[messageId] = messsage
}

func (r *messageRepository) get(messageId string) *message {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.m[messageId]
}
