package kivaio

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/st3v/tracerr"
)

const (
	protocol  = 1
	transport = "websocket"
	readLimit = 65536
)

type Session interface {
	Connect(channel string) (<-chan string, error)
}

type session struct {
	socket           *websocket.Conn
	handler          Handler
	socketId         string
	heartbeatTimeout time.Duration
	closeTimeout     time.Duration
	transports       []string
	host             string
	protocol         int
}

func NewSession(host string) (Session, error) {
	url := fmt.Sprintf("http://%s/socket.io/%d?t=%d", host, protocol, time.Now().Unix())
	fmt.Println(url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	result, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	fmt.Println(string(result))

	parts := strings.Split(string(result), ":")

	session := &session{
		socketId:         parts[0],
		heartbeatTimeout: parseDuration(parts[1]),
		closeTimeout:     parseDuration(parts[2]),
		transports:       strings.Split(parts[3], ","),
		host:             host,
		protocol:         protocol,
	}

	return session, nil
}

func (s *session) Connect(name string) (<-chan string, error) {
	name = fmt.Sprintf("/%s", name)

	err := s.openSocket()
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	return s.handler.AddChannel(name)
}

func (s *session) openSocket() error {
	if s.socket != nil {
		return nil
	}

	url := fmt.Sprintf("ws://%s/socket.io/%d/%s/%s", s.host, s.protocol, transport, s.socketId)
	fmt.Println(url)

	socket, resp, err := websocket.DefaultDialer.Dial(url, http.Header{})

	if resp != nil && resp.StatusCode == http.StatusUnauthorized {
		bodyData, _ := ioutil.ReadAll(resp.Body)
		return tracerr.Wrap(fmt.Errorf("Response error: %s\n", string(bodyData)))
	}

	if err != nil {
		return tracerr.Wrap(err)
	}

	s.socket = socket
	s.socket.SetReadLimit(readLimit)

	s.handler = newHandler(newListener(s.socket, s.closeTimeout), newSender(s.socket))
	s.handler.Handle()

	return nil
}

func parseDuration(str string) time.Duration {
	i, _ := strconv.Atoi(str)
	return time.Duration(i) * time.Second
}
