package socketio

import (
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

type Listener interface {
	Listen() (<-chan *message, <-chan bool)
}

type listener struct {
	socket       *websocket.Conn
	closeTimeout time.Duration
}

func newListener(socket *websocket.Conn, closeTimeout time.Duration) Listener {
	return &listener{
		socket:       socket,
		closeTimeout: closeTimeout,
	}
}

func (l *listener) Listen() (<-chan *message, <-chan bool) {
	parser := newParser()

	messages := make(chan *message)
	ready := make(chan bool)

	l.resetCloseTimeout()

	go func() {
		for {
			_, p, err := l.socket.ReadMessage()
			if err != nil {
				fmt.Printf("Error reading message: %s\n", err.Error())
				continue
			}
			message, err := parser.Parse(string(p))
			if err != nil {
				fmt.Printf("Error parsing message: %s\n", err.Error())
				continue
			}

			if message.category != DISCONNECT {
				l.resetCloseTimeout()
			}

			if message.endpoint == "" {
				if message.category == CONNECT {
					ready <- true
				} else if message.category == DISCONNECT {
					ready <- false
				}
			}

			messages <- message
		}
	}()

	return messages, ready
}

func (l *listener) resetCloseTimeout() {
	l.socket.SetReadDeadline(time.Now().Add(l.closeTimeout))
}
