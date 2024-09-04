package messaging

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Mohamed-khattab/Message-Queue/utils"
	"github.com/google/uuid"
)

func NewBroker() *Broker {
	return &Broker{
		topics: make(map[string]Topic),
	}
}

func NewSubscriber(endpoint string, topics []string) Subscriber {
	return Subscriber{
		ID:            uuid.New().String(),
		EndPoint:      endpoint,
		Topics:        topics,
		LastHeartbeat: time.Now().Format(time.RFC3339),
		status:        ACTIVE,
	}
}

func NewMessage(body string) Message {
	return Message{
		ID:           uuid.New().String(),
		Body:         body,
		PUBLISHED_AT: time.Now().Format(time.RFC3339),
		CONSUMED_AT:  "",
	}
}

func (b *Broker) Subscribe(endpoint string, topics []string) (string, error) {
	if !utils.IsValidURL(endpoint) {
		return "", fmt.Errorf("invalid endpoint URL: %s", endpoint)
	}

	subscriber := NewSubscriber(endpoint, topics)

	b.mu.Lock()
	defer b.mu.Unlock()

	for _, topicName := range topics {
		// Initialize the topic if it doesn't exist
		if _, exists := b.topics[topicName]; !exists {
			newTopic := Topic{
				queue:       []Message{},
				subscribers: []Subscriber{},
			}
			b.topics[topicName] = newTopic
		}

		isSubscriberExists := false
		for _, sub := range b.topics[topicName].subscribers {
			if sub.EndPoint == endpoint {
				isSubscriberExists = true
				break
			}
		}

		if !isSubscriberExists {
			topic := b.topics[topicName]
			topic.subscribers = append(topic.subscribers, subscriber)
			b.topics[topicName] = topic
		}
	}

	return subscriber.ID, nil
}

func (b *Broker) Unsubscribe(subscriberID string, topics []string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	for _, topicName := range topics {

		topic, exist := b.topics[topicName]
		if !exist {
			return fmt.Errorf("topic %s not found", topicName)
		}

		var index int = 0
		for index, sub := range topic.subscribers {
			if sub.ID == subscriberID {
				topic.subscribers = append(topic.subscribers[:index], topic.subscribers[index+1:]...)
				break
			}
		}
		if len(topic.subscribers) == index {
			return fmt.Errorf("subscriber %s not subscribed to topic %s", subscriberID, topicName)
		}
	}
	return nil
}

func (b *Broker) Publish(topicName string, messageBody string) error {

	topic, exist := b.topics[topicName]
	if !exist {
		return fmt.Errorf("topic %s not found", topicName)
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	message := NewMessage(messageBody)
	topic.queue = append(topic.queue, message)

	b.propagateMessages(&topic)
	b.topics[topicName] = topic

	return nil
}

func (b *Broker) Retrieve(subId string, topicName string, startDate string) ([]Message, error) {

	b.mu.Lock()
	defer b.mu.Unlock()

	topic, exist := b.topics[topicName]
	if !exist {
		return nil, fmt.Errorf("topic %s not found", topicName)
	}

	// Check if the subscriber exists
	var subscriberExists bool
	for _, sub := range topic.subscribers {
		if sub.ID == subId {
			subscriberExists = true
			break
		}
	}

	if !subscriberExists {
		return nil, fmt.Errorf("subscriber %s does not exist for topic %s", subId, topicName)
	}

	if startDate == "" {
		return topic.queue, nil
	}

	var messages []Message
	for _, msg := range topic.queue {
		if msg.PUBLISHED_AT >= startDate {
			messages = append(messages, msg)
		}
	}

	return messages, nil
}

func (b *Broker) propagateMessages(topic *Topic) error {

	for msgIndex := range topic.queue {
		msg := &topic.queue[msgIndex]
		if msg.CONSUMED_AT == "" {
			for subIndex := range topic.subscribers {
				sub := &topic.subscribers[subIndex]
				if sub.status == ACTIVE {

					// send the message to the subscriber
					body := []byte(`{
						"Message": "` + msg.Body + `",
						"Published_at": "` + msg.PUBLISHED_AT + `",
						}`)

					request, err := http.NewRequest("POST", sub.EndPoint, bytes.NewBuffer(body))
					if err != nil {
						return fmt.Errorf("failed to construct request for subescriber%s: %v", sub.ID, err)
					}
					request.Header.Set("Content-Type", "application/json; charset=utf-8")

					client := &http.Client{}
					res, err := client.Do(request)

					if err != nil {
						return fmt.Errorf("failed to send message to subscriber %s: %v", sub.ID, err)
					}

					_, err = io.ReadAll(res.Body)

					if err != nil {
						return fmt.Errorf("failed to parse the reponse from subscriber %s: %v", sub.ID, err)
					}

					defer res.Body.Close()
					fmt.Println("Message sent to subscriber", sub.ID)
					sub.LastHeartbeat = time.Now().Format(time.RFC3339)
				}
			}
			msg.CONSUMED_AT = time.Now().Format(time.RFC3339)
		}
	}
	return nil
}
