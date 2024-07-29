package messaging

import (
	"sync"
)

type Message struct {
	ID int 
	Body string
}

type Broker struct {
	mu sync.Mutex
	queues map[string][]Message
	topics map[string][]string
}


func NewBroker () *Broker {
	return &Broker{
		queues: make(map[string][]Message),
		topics: make(map[string][]string),
	}
}

func (b *Broker) Publish (topic string, msg Message) int {

 return 0 ;	 

}



func (b * Broker) Retrieve (queue string ) (Message, bool) {

	return Message{
		ID: 12,
		Body: "hello this is me ", 
	}, true 
}

