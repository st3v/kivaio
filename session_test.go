package kivaio

import (
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
	expectedTransports := []string{"fake-transport-1", "fake-transport-2", "fake-transport-3"}

	handshakeRequested := false

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handshakeRequested = true

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
	defer server.Close()

	handshakeURL = func() string {
		return server.URL
	}

	session, err := NewSession()

	if err != nil {
		t.Errorf("Unexpected error returned by NewSession: %s", err.Error())
	}

	if !handshakeRequested {
		t.Error("Handshake expected but did not happen.")
	}

	if session.socketID != expectedSocketID {
		t.Errorf("Unexpected socketID. Want: '%s'. Got: '%s'.", expectedSocketID, session.socketID)
	}

	if session.heartbeatTimeout != time.Duration(expectedHeartbeatTimeout)*time.Second {
		t.Errorf("Unexpected hartbeat timeout in session.")
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
