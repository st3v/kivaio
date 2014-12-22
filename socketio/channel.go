package socketio

import "fmt"

type Channel interface {
	Received(message string)
	Messages() <-chan string
}

type channel struct {
	name     string
	messages chan string
}

func newChannel(name string, sender Sender) (Channel, error) {
	err := sender.Send(fmt.Sprintf("%d::%s", CONNECT, name))
	if err != nil {
		fmt.Printf("Error sending connect message for channel '%s': %s\n", name, err.Error())
		return nil, err
	}

	return &channel{
		name:     name,
		messages: make(chan string),
	}, nil
}

func (c *channel) Received(message string) {
	c.messages <- message
}

func (c *channel) Messages() <-chan string {
	return c.messages
}
