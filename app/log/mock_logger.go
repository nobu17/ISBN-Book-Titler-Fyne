package log

type mockLogger struct {
}

func (m *mockLogger) Info(message string) {
}

func (m *mockLogger) Warn(message string) {
}

func (z *mockLogger) Error(message string, err error) {
}

func NewMockLogger() AppLogger {
	return &mockLogger{}
}