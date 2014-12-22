package socketio

import (
	"errors"
	"fmt"

	"github.com/gorilla/websocket"
)

type Sender interface {
	Send(message string) error
}

type sender struct {
	socket *websocket.Conn
}

func newSender(socket *websocket.Conn) Sender {
	return &sender{
		socket: socket,
	}
}

func (s *sender) Send(message string) error {
	if s.socket == nil {
		err := errors.New("Socket Not Initialized")
		fmt.Printf("Error sending message: %s\n", err.Error())
		return err
	}

	writer, err := s.socket.NextWriter(1)
	if err != nil {
		fmt.Printf("Error obtaining writer from socket: %s\n", err.Error())
		return err
	}

	writer.Write([]byte(message))

	writer.Close()

	return nil
}
