package socketio

const (
	DISCONNECT = 0
	CONNECT    = 1
	HEARTBEAT  = 2
	MESSAGE    = 3
	JSON       = 4
	EVENT      = 5
	ACK        = 6
	ERROR      = 7
	NOOP       = 8
)

type message struct {
	id       string
	category int
	data     string
	endpoint string
	ack      string
}
