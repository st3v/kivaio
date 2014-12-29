package kivaio

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/st3v/tracerr"
)

const (
	endpoint  = "streams.kiva.org"
	protocol  = 1
	transport = "websocket"
	readLimit = 65536
)

type Session interface {
	Connect(channel string) (<-chan string, error)
}

type session struct {
	socket           Socket
	socketID         string
	heartbeatTimeout time.Duration
	closeTimeout     time.Duration
	transports       []string
	host             string
	protocol         int
}

var handshakeURL = func() string {
	return fmt.Sprintf("http://%s/socket.io/%d?t=%d", endpoint, protocol, clock.Now().Unix())
}

func NewSession() (*session, error) {
	resp, err := http.Get(handshakeURL())
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	result, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	parts := strings.Split(string(result), ":")

	transports := strings.Split(parts[3], ",")
	err = assertTransport(transports)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	session := &session{
		socketID:         parts[0],
		heartbeatTimeout: parseDuration(parts[1]),
		closeTimeout:     parseDuration(parts[2]),
		transports:       transports,
		host:             endpoint,
		protocol:         protocol,
	}

	return session, nil
}

func (s *session) Connect(name string) (<-chan string, error) {
	if !strings.HasPrefix(name, "/") {
		name = fmt.Sprintf("/%s", name)
	}

	if s.socket == nil {
		socket, err := openSocket(s.host, s.socketID, s.protocol, s.closeTimeout)
		if err != nil {
			return nil, tracerr.Wrap(err)
		}
		s.socket = socket
	}

	return s.socket.OpenChannel(name)
}

func assertTransport(transports []string) error {
	for _, t := range transports {
		if t == transport {
			return nil
		}
	}
	return tracerr.Errorf("Transport '%s' not supported by server.", transport)
}

func parseDuration(str string) time.Duration {
	i, _ := strconv.Atoi(str)
	return time.Duration(i) * time.Second
}
