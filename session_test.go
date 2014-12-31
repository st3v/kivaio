package kivaio

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"testing"
	"time"
)

func TestHandshakeURL(t *testing.T) {
	clock = newMockClock(time.Unix(1234567890, 0))
	expectedURL := "http://streams.kiva.org/socket.io/1?t=1234567890"
	actualURL := handshakeURL()
	if actualURL != expectedURL {
		t.Errorf("Unexpected handshake URL. Want: '%s'. Got: '%s'.", expectedURL, actualURL)
	}
}

func TestNewSession(t *testing.T) {
	expectedSocketID := "fake-socket-id"
	expectedHeartbeatTimeout := 12
	expectedCloseTimeout := 34
	expectedTransports := []string{"websocket", "fake-transport-1", "fake-transport-2"}

	server := mockHandshakeServer(
		expectedSocketID,
		expectedHeartbeatTimeout,
		expectedCloseTimeout,
		expectedTransports,
		t,
	)
	defer server.Close()

	session, err := NewSession()

	if err != nil {
		t.Errorf("Unexpected error returned by NewSession: %s", err.Error())
	}

	if session.socketID != expectedSocketID {
		t.Errorf("Unexpected socketID. Want: '%s'. Got: '%s'.", expectedSocketID, session.socketID)
	}

	if session.heartbeatTimeout != time.Duration(expectedHeartbeatTimeout)*time.Second {
		t.Errorf("Unexpected heartbeat timeout in session.")
	}

	if session.closeTimeout != time.Duration(expectedCloseTimeout)*time.Second {
		t.Errorf("Unexpected close timeout in session.")
	}

	if len(session.transports) != len(expectedTransports) {
		t.Errorf(
			"Unexpected transports in sessions. Want: [%s]. Got: [%s]",
			strings.Join(expectedTransports, ","),
			strings.Join(session.transports, ","),
		)
	}

	sort.Strings(expectedTransports)
	sort.Strings(session.transports)

	for i, _ := range session.transports {
		if session.transports[i] != expectedTransports[i] {
			t.Errorf(
				"Unexpected transport in session. Want: '%s'. Got: '%s'",
				expectedTransports[i],
				session.transports[i],
			)
		}
	}
}

func TestNewSessionUnsupportedTransport(t *testing.T) {
	expectedSocketID := "fake-socket-id"
	expectedHeartbeatTimeout := 12
	expectedCloseTimeout := 34
	expectedTransports := []string{"fake-transport-1", "fake-transport-2", "fake-transport-3"}

	server := mockHandshakeServer(
		expectedSocketID,
		expectedHeartbeatTimeout,
		expectedCloseTimeout,
		expectedTransports,
		t,
	)
	defer server.Close()

	_, err := NewSession()

	if err == nil {
		t.Fatal("Expected error but got none.")
	}

	if !strings.HasPrefix(err.Error(), "Transport 'websocket' not supported by server") {
		t.Errorf("Unexpected error: %s", err.Error())
	}
}

func TestSessionConnectNoSlash(t *testing.T) {
	expectedHost := "fake-hostname"
	expectedSocketID := "fake-socket-id"
	expectedProtocol := 123
	expectedCloseTimeout := time.Duration(456) * time.Second
	expectedChannel := "fake-channel"
	expectedChan := make(<-chan string)
	expectedError := errors.New("Expected Error")

	mockSession := mockSocketOpenChannel(
		expectedHost,
		expectedSocketID,
		expectedProtocol,
		expectedCloseTimeout,
		fmt.Sprintf("/%s", expectedChannel),
		expectedChan,
		expectedError,
		t,
	)

	actualChan, actualError := mockSession.Connect(expectedChannel)

	if actualChan != expectedChan {
		t.Errorf("Unexpected chan returned by Connect().")
	}

	if actualError != expectedError {
		t.Errorf(
			"Unexpected error returned by Connect(). Want: '%s'. Got: '%s'.",
			expectedError.Error(),
			actualError.Error(),
		)
	}
}

func TestSessionConnectSlash(t *testing.T) {
	expectedHost := "fake-hostname"
	expectedSocketID := "fake-socket-id"
	expectedProtocol := 123
	expectedCloseTimeout := time.Duration(456) * time.Second
	expectedChannel := "/fake-channel"
	expectedChan := make(<-chan string)
	expectedError := errors.New("Expected Error")

	mockSession := mockSocketOpenChannel(
		expectedHost,
		expectedSocketID,
		expectedProtocol,
		expectedCloseTimeout,
		expectedChannel,
		expectedChan,
		expectedError,
		t,
	)

	actualChan, actualError := mockSession.Connect(expectedChannel)

	if actualChan != expectedChan {
		t.Errorf("Unexpected chan returned by Connect().")
	}

	if actualError != expectedError {
		t.Errorf(
			"Unexpected error returned by Connect(). Want: '%s'. Got: '%s'.",
			expectedError.Error(),
			actualError.Error(),
		)
	}
}

func TestSessionConnectOpenSocketError(t *testing.T) {
	expectedErrorMsg := "fake-open-socket-error"

	openSocket = func(host, socketID string, protocol int, closeTimeout time.Duration) (Socket, error) {
		return nil, errors.New(expectedErrorMsg)
	}

	session := &session{}

	msgChan, err := session.Connect("/fake-channel")

	if !strings.HasPrefix(err.Error(), expectedErrorMsg) {
		t.Errorf(
			"Unexpected error returned by Connect(). Want: '%s'. Got: '%s'.",
			expectedErrorMsg,
			err.Error(),
		)
	}

	if msgChan != nil {
		t.Error("Expected returned channel to be nil. It's not.")
	}

}

func mockSocketOpenChannel(
	expectedHost string,
	expectedSocketID string,
	expectedProtocol int,
	expectedCloseTimeout time.Duration,
	expectedChannelName string,
	expectedChan <-chan string,
	expectedError error,
	t *testing.T,
) *session {
	mockSession := &session{
		socketID:     expectedSocketID,
		host:         expectedHost,
		protocol:     expectedProtocol,
		closeTimeout: expectedCloseTimeout,
	}

	openSocket = func(host, socketID string, protocol int, closeTimeout time.Duration) (Socket, error) {
		if host != expectedHost {
			t.Errorf("Unexpected hostname. Want: '%s'. Got: '%s'.", expectedHost, host)
		}

		if socketID != expectedSocketID {
			t.Errorf("Unexpected scoket id. Want: '%s'. Got: '%s'.", expectedSocketID, socketID)
		}

		if protocol != expectedProtocol {
			t.Errorf("Unexpected protocol. Want: '%d'. Got: '%d'.", expectedProtocol, protocol)
		}

		if closeTimeout != expectedCloseTimeout {
			t.Errorf(
				"Unexpected close timeout. Want: '%d'. Got: '%d'.",
				expectedCloseTimeout.Seconds(),
				closeTimeout.Seconds(),
			)
		}

		return &mockSocket{
			openChannel: func(name string) (<-chan string, error) {
				if name != expectedChannelName {
					t.Errorf(
						"Unexpected channel name passed to socket. Want: '%s'. Got: '%s'.",
						expectedChannelName,
						name,
					)
				}
				return expectedChan, expectedError
			},
		}, nil
	}

	return mockSession
}

func mockHandshakeServer(
	expectedSocketID string,
	expectedHeartbeatTimeout int,
	expectedCloseTimeout int,
	expectedTransports []string,
	t *testing.T,
) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Fatalf("Unexpected request method: %s", r.Method)
		}

		fmt.Fprintf(
			w,
			"%s:%d:%d:%s",
			expectedSocketID,
			expectedHeartbeatTimeout,
			expectedCloseTimeout,
			strings.Join(expectedTransports, ","),
		)

		return
	}))

	handshakeURL = func() string {
		return server.URL
	}

	return server
}
