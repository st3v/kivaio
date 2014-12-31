package kivaio

import (
	"errors"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
)

func TestSocketURL(t *testing.T) {
	expectedURL := "ws://fake-hostname/socket.io/123/websocket/fake-socket-id"
	actualURL := socketURL("fake-hostname", "fake-socket-id", 123)
	if actualURL != expectedURL {
		t.Errorf("Unexpected socket URL. Want: '%s'. Got: '%s'.", expectedURL, actualURL)
	}
}

func TestOpenChannelWithNilConnection(t *testing.T) {
	socket := &socket{}
	msgChan, err := socket.OpenChannel("fake-channel")

	if err == nil {
		t.Error("Expected error returned by OpenChannel but got none.")
	}

	expectedErrorMsg := "Socket connection not initialized"

	if !strings.HasPrefix(err.Error(), expectedErrorMsg) {
		t.Errorf(
			"Unexpected error returned by OpenChannel(). Want: '%s'. Got: '%s'.",
			expectedErrorMsg,
			err.Error(),
		)
	}

	if msgChan != nil {
		t.Error("Expected channel returned by OpenChannel() to be nil. It's not.")
	}
}

func TestOpenChannel(t *testing.T) {
	expectedChannelName := "/fakeChannel"
	expectedChan := make(<-chan string)
	expectedError := errors.New("fake-error")

	mockSocketHandler := newMockSocketHandler()
	mockSocketHandler.openChannel = func(name string) (<-chan string, error) {
		if name != expectedChannelName {
			t.Errorf(
				"Unexpected error returned by OpenChannel(). Want: '%s'. Got: '%s'.",
				expectedChannelName,
				name,
			)
		}

		return expectedChan, expectedError
	}

	socket := &socket{
		conn:    &websocket.Conn{},
		handler: mockSocketHandler,
	}

	actualChan, actualError := socket.OpenChannel(expectedChannelName)

	if actualChan != expectedChan {
		t.Errorf("Unexpected chan returned by OpenChannel().")
	}

	if actualError != expectedError {
		t.Errorf(
			"Unexpected error returned by OpenChannel(). Want: '%s'. Got: '%s'.",
			expectedError.Error(),
			actualError.Error(),
		)
	}
}

type mockSocket struct {
	openChannel func(name string) (<-chan string, error)
}

func (m *mockSocket) OpenChannel(name string) (<-chan string, error) {
	return m.openChannel(name)
}
