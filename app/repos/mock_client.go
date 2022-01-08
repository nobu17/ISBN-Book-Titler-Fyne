package repos

type mockClient struct {
	getfunc func() ([]byte, error)
}

func (m *mockClient) Get(url string, params map[string]string) ([]byte, error) {
	return m.getfunc()
}

func NewMockClient(getfunc func() ([]byte, error)) Client {
	return &mockClient{getfunc}
}
