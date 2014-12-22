package socketio

import (
	"fmt"
	"time"
)

type Handler interface {
	Handle()
	AddChannel(name string) (<-chan string, error)
}

type handler struct {
	listener       Listener
	sender         Sender
	channels       map[string]Channel
	socketMessages <-chan message
}

func newHandler(listener Listener, sender Sender) *handler {
	return &handler{
		listener: listener,
		sender:   sender,
		channels: make(map[string]Channel),
	}
}

func (h *handler) AddChannel(name string) (<-chan string, error) {
	if h.channels[name] == nil {
		channel, err := newChannel(name, h.sender)
		if err != nil {
			fmt.Printf("Error instantiating new channel: %s\n", err.Error())
			return nil, err
		}
		h.channels[name] = channel
	}
	return h.channels[name].Messages(), nil
}

func (h *handler) Handle() {
	socketMessages, ready := h.listener.Listen()

	if <-ready == true {
		go func() {
			for {
				socketMessage := <-socketMessages

				switch socketMessage.category {
				case CONNECT:
					if socketMessage.endpoint != "" {
						fmt.Printf("Connected to channel '%s'\n", socketMessage.endpoint)
					} else {
						fmt.Println("Connected to socket")
					}
				case HEARTBEAT:
					fmt.Printf("Heartbeat: %s\n", time.Now())
					err := h.sender.Send(fmt.Sprintf("%d::", HEARTBEAT))
					if err != nil {
						continue
					}
				case MESSAGE:
					c := h.channels[socketMessage.endpoint]
					if c != nil {
						c.Received(socketMessage.data)
					}
				}
			}
		}()
	}
}
