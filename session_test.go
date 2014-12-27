package kivaio

import (
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
