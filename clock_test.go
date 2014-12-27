package kivaio

import (
	"testing"
	"time"
)

func TestClockNow(t *testing.T) {
	delta := clock.Now().Unix() - time.Now().Unix()
	if 0 < delta || delta > 1 {
		t.Errorf("Unexpected delta %d between clock.Now() and time.Now()", delta)
	}
}

type mockClock struct {
	now time.Time
}

func newMockClock(now time.Time) Clock {
	return &mockClock{
		now: now,
	}
}

func (m *mockClock) Now() time.Time {
	return m.now
}
