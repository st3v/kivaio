package kivaio

type mockSocketHandler struct {
	openChannel func(name string) (<-chan string, error)
	handle      func()
}

func newMockSocketHandler() *mockSocketHandler {
	return &mockSocketHandler{
		openChannel: func(name string) (<-chan string, error) {
			return nil, nil
		},
		handle: func() {},
	}
}

func (m *mockSocketHandler) OpenChannel(name string) (<-chan string, error) {
	return m.openChannel(name)
}

func (m *mockSocketHandler) Handle() {}
