package kivaio

import (
	"errors"

	"github.com/gorilla/websocket"
	"github.com/st3v/tracerr"
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
		err := errors.New("Socket not initialized")
		return tracerr.Wrap(err)
	}

	writer, err := s.socket.NextWriter(1)
	if err != nil {
		return tracerr.Wrap(err)
	}

	writer.Write([]byte(message))

	writer.Close()

	return nil
}
