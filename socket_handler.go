package kivaio

import (
	"fmt"
	"log"
	"time"

	"github.com/st3v/tracerr"
)

type SocketHandler interface {
	Handle()
	OpenChannel(name string) (<-chan string, error)
}

type socketHandler struct {
	listener       Listener
	sender         Sender
	channels       map[string]Channel
	socketMessages <-chan message
}

func newSocketHandler(listener Listener, sender Sender) SocketHandler {
	return &socketHandler{
		listener: listener,
		sender:   sender,
		channels: make(map[string]Channel),
	}
}

func (s *socketHandler) OpenChannel(name string) (<-chan string, error) {
	if s.channels[name] == nil {
		channel, err := newChannel(name, s.sender)
		if err != nil {
			return nil, tracerr.Wrap(err)
		}
		s.channels[name] = channel
	}
	return s.channels[name].Messages(), nil
}

func (s *socketHandler) Handle() {
	socketMessages, ready := s.listener.Listen()

	if <-ready == true {
		go func() {
			for {
				socketMessage := <-socketMessages

				switch socketMessage.category {
				case CONNECT:
					if socketMessage.endpoint != "" {
						log.Printf("Connected to channel '%s'\n", socketMessage.endpoint)
					} else {
						log.Println("Connected to socket")
					}
				case HEARTBEAT:
					fmt.Printf("Heartbeat: %s\n", time.Now())
					err := s.sender.Send(fmt.Sprintf("%d::", HEARTBEAT))
					if err != nil {
						log.Printf("Error sending heartbeat: %s\n", err.Error())
						continue
					}
				case MESSAGE:
					c := s.channels[socketMessage.endpoint]
					if c != nil {
						c.Received(socketMessage.data)
					}
				}
			}
		}()
	}
}
