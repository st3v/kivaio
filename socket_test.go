package kivaio

import "testing"

func TestSocketURL(t *testing.T) {
	expectedURL := "ws://fake-hostname/socket.io/123/websocket/fake-socket-id"
	actualURL := socketURL("fake-hostname", "fake-socket-id", 123)
	if actualURL != expectedURL {
		t.Errorf("Unexpected socket URL. Want: '%s'. Got: '%s'.", expectedURL, actualURL)
	}
}
