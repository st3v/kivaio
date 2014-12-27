package kivaio

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/st3v/tracerr"
)

type Socket interface {
	OpenChannel(name string) (<-chan string, error)
}

type socket struct {
	conn    *websocket.Conn
	handler SocketHandler
}

var socketURL = func(host, socketID string, protocol int) string {
	return fmt.Sprintf("ws://%s/socket.io/%d/websocket/%s", host, protocol, socketID)
}

func (s *socket) OpenChannel(name string) (<-chan string, error) {
	if s.conn == nil {
		return nil, tracerr.Error("SocketHandler not initialized")
	}
	return s.handler.OpenChannel(name)
}

func openSocket(host, socketID string, protocol int, closeTimeout time.Duration) (Socket, error) {
	url := socketURL(host, socketID, protocol)

	conn, resp, err := websocket.DefaultDialer.Dial(url, http.Header{})

	if resp != nil && resp.StatusCode == http.StatusUnauthorized {
		bodyData, _ := ioutil.ReadAll(resp.Body)
		return nil, tracerr.Errorf("Response error: %s\n", string(bodyData))
	}

	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	conn.SetReadLimit(readLimit)

	handler := newSocketHandler(newListener(conn, closeTimeout), newSender(conn))
	handler.Handle()

	return &socket{
		conn:    conn,
		handler: handler,
	}, nil
}
