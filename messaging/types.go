package messaging

import (
	"sync"
)

type Broker struct {
	mu     sync.Mutex
	topics map[string]Topic
}

type Topic struct {
	queue       []Message
	subscribers []Subscriber
}

type Message struct {
	ID           string
	Body         string `db:"body" json:"body"`
	PUBLISHED_AT string `db:"published_at"`
	CONSUMED_AT  string `db:"consumed_at"`
}

type Subscriber struct {
	ID            string
	EndPoint      string
	Topics        []string
	LastHeartbeat string
	status        SubscriberStatus
}

type SubscriberStatus string

const (
	ACTIVE   SubscriberStatus = "active"
	INACTIVE SubscriberStatus = "inactive"
	FAILED   SubscriberStatus = "failed"
)

// Requests
type SubscribeRequest struct {
	Endpoint string   `json:"endpoint"`
	Topics   []string `json:"topics"`
}
type UnsubscribeRequest struct {
	SubscriberID string   `json:"subscriber_id"`
	Topics       []string `json:"topics"`
}

type PublishRequest struct {
	Topic   string `json:"topic"`
	Message string `json:"message"`
}

type RetrieveRequest struct {
	Topic        string `json:"topic"`
	SubscriberId string `json:"subId"`
	StartDate    string `json:"startDate"`
}