package main

import (
	"fmt"
	"sync"
	"time"
)

// Message represents a message that will be passed through the broker.
type Message struct {
	Topic   string
	Payload interface{}
}

// Subscriber represents a subscriber to a topic.
type Subscriber struct {
	Channel     chan interface{}
	Unsubscribe chan bool
}

// Broker is the in-memory message broker.
type Broker struct {
	subscribers map[string][]*Subscriber
	mutex       sync.Mutex
}

// NewBroker creates a new instance of the in-memory broker.
func NewBroker() *Broker {
	return &Broker{
		subscribers: make(map[string][]*Subscriber),
	}
}

// Subscribe allows a client to subscribe to a topic and receive messages.
func (b *Broker) Subscribe(topic string) *Subscriber {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	subscriber := &Subscriber{
		Channel:     make(chan interface{}, 1),
		Unsubscribe: make(chan bool),
	}

	b.subscribers[topic] = append(b.subscribers[topic], subscriber)

	return subscriber
}

// Unsubscribe removes a subscriber from a topic.
func (b *Broker) Unsubscribe(topic string, subscriber *Subscriber) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if subscribers, found := b.subscribers[topic]; found {
		for i, sub := range subscribers {
			if sub == subscriber {
				close(sub.Channel)
				b.subscribers[topic] = append(subscribers[:i], subscribers[i+1:]...)
				return
			}
		}
	}
}

// Publish sends a message to all subscribers of a specific topic.
func (b *Broker) Publish(topic string, payload interface{}) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if subscribers, found := b.subscribers[topic]; found {
		for _, sub := range subscribers {
			select {
			case sub.Channel <- payload:
			case <-time.After(time.Second):
				fmt.Printf("Subscriber slow. Unsubscribing from topic: %s\n", topic)
				b.Unsubscribe(topic, sub)
			}
		}
	}
}

func main() {
	broker := NewBroker()

	go func() {
		for {
			time.Sleep(1 * time.Second)
			fmt.Println("Hello there")
		}
	}()

	subscriber := broker.Subscribe("example_topic")
	go func() {
		for {
			select {
			case msg, ok := <-subscriber.Channel:
				if !ok {
					fmt.Println("Subscriber channel closed.")
					return
				}
				fmt.Printf("Received: %v\n", msg)
			case <-subscriber.Unsubscribe:
				fmt.Println("Unsubscribed.")
				return
			}
		}
	}()

	broker.Publish("example_topic", "Hello, World!")
	broker.Publish("example_topic", "This is a test message.")

	time.Sleep(2 * time.Second)
	broker.Unsubscribe("example_topic", subscriber)

	broker.Publish("example_topic", "This message won't be received.")

	time.Sleep(time.Second)
}
