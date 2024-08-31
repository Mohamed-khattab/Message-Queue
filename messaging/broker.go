package messaging

import (
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID           int
	Body         string `db:"body" json:"body"`
	PUBLISHED_AT string `db:"published_at" json:"published_at"`
	CONSUMED_AT  string `db:"consumed_at" json:"consumed_at"`
}

type Broker struct {
	mu     sync.Mutex
	queues map[string][]Message
	topics map[string][]Subscriber
}

func NewBroker() *Broker {
	return &Broker{
		queues: make(map[string][]Message),
		topics: make(map[string][]Subscriber),
	}
}

type SubscribeRequest struct {
	Endpoint string   `json:"endpoint"`
	Topics   []string `json:"topics"`
}
type UnsubscribeRequest struct {
	SubscriberID string   `json:"subscriber_id"`
	Topics       []string `json:"topics"`
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

func (b *Broker) Subscribe(endpoint string, topics []string) (string, error) {
	if !isValidURL(endpoint) {
		return "", fmt.Errorf("invalid endpoint URL: %s", endpoint)
	}

	// Create a new subscriber
	subscriber := Subscriber{
		ID:            uuid.New().String(),
		EndPoint:      endpoint,
		Topics:        topics,
		LastHeartbeat: time.Now().Format(time.RFC3339),
		status:        ACTIVE,
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	for _, topic := range topics {
		// Initialize the topic if it doesn't exist
		if _, exists := b.topics[topic]; !exists {
			b.topics[topic] = []Subscriber{}
		}

		subscriberExists := false
		for _, sub := range b.topics[topic] {
			if sub.EndPoint == endpoint {
				subscriberExists = true
				break
			}
		}

		if !subscriberExists {
			b.topics[topic] = append(b.topics[topic], subscriber)
		}
	}

	return subscriber.ID, nil
}

func (b *Broker) Unsubscribe(subscriberID string, topics []string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	for _, topic := range topics {
		subscribers, exist := b.topics[topic]
		if !exist {
			return fmt.Errorf("topic %s not found", topic)
		}

		var index int = 0
		for index, sub := range subscribers {
			if sub.ID == subscriberID {
				b.topics[topic] = append(b.topics[topic][:index], b.topics[topic][index+1:]...)
				break
			}
		}
		if len(b.topics[topic]) == index {
			return fmt.Errorf("subscriber %s not subscribed to topic %s", subscriberID, topic)
		}
	}
	return nil
}

func (b *Broker) Publish(topic string, msg Message) int {

	return 0

}

func (b *Broker) Retrieve(queue string) (Message, bool) {

	return Message{
		ID:   12,
		Body: "hello this is me ",
	}, true
}

func isValidURL(str string) bool {
	parsedURL, err := url.Parse(str)
	if err != nil || parsedURL == nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return false
	}
	return true
}
