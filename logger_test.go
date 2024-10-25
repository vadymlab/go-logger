package log

import (
	"context"
	"testing"
)

// Test Text2Level function to ensure string values are correctly converted to LogLevel
func TestText2Level(t *testing.T) {
	tests := []struct {
		level string
		want  LogLevel
	}{
		{"TRACE", TraceLevel},
		{"DEBUG", DebugLevel},
		{"INFO", InfoLevel},
		{"WARNING", WarnLevel},
		{"ERROR", ErrorLevel},
		{"FATAL", FatalLevel},
		{"PANIC", PanicLevel},
		{"UNKNOWN", InfoLevel}, // assuming default is InfoLevel
	}

	for _, tt := range tests {
		got := Text2Level(tt.level)
		if got != tt.want {
			t.Errorf("Text2Level(%v) = %v; want %v", tt.level, got, tt.want)
		}
	}
}

// Test GetDefaultLogger function to check the initialization of the default logger
func TestGetDefaultLogger(t *testing.T) {
	SetDefaultLogger(nil) // Reset default logger for test
	logger := GetDefaultLogger()
	if logger == nil {
		t.Fatal("Expected default logger, got nil")
	}
}

// Test ToContext and FromContext to check adding and retrieving logger from context
func TestToContextAndFromContext(t *testing.T) {
	ctx := context.Background()
	mockLogger := &MockLogger{} // Mock Logger for testing
	ctx = ToContext(ctx, mockLogger)

	retrievedLogger := FromContext(ctx)
	if retrievedLogger != mockLogger {
		t.Errorf("Expected logger from context, but got different instance")
	}
}

// Test FromDefaultContext to check initialization and retrieval of logger from defaultContext
func TestFromDefaultContext(t *testing.T) {
	SetDefaultContext(nil) // Ensure defaultContext is initialized for the test
	logger := FromDefaultContext()
	if logger == nil {
		t.Fatal("Expected default logger from default context, got nil")
	}
}

// MockLogger to simulate a logger in tests
type MockLogger struct{}

func (m *MockLogger) Info(args ...interface{})                        {}
func (m *MockLogger) Infof(format string, args ...interface{})        {}
func (m *MockLogger) Infow(msg string, keysAndValues ...interface{})  {}
func (m *MockLogger) Warn(args ...interface{})                        {}
func (m *MockLogger) Warnf(format string, args ...interface{})        {}
func (m *MockLogger) Warnw(msg string, keysAndValues ...interface{})  {}
func (m *MockLogger) Error(args ...interface{})                       {}
func (m *MockLogger) Errorf(format string, args ...interface{})       {}
func (m *MockLogger) Errorw(msg string, keysAndValues ...interface{}) {}
func (m *MockLogger) Debug(args ...interface{})                       {}
func (m *MockLogger) Debugf(format string, args ...interface{})       {}
func (m *MockLogger) Debugw(msg string, keysAndValues ...interface{}) {}
func (m *MockLogger) Fatal(args ...interface{})                       {}
func (m *MockLogger) Fatalf(format string, args ...interface{})       {}
func (m *MockLogger) With(f ...interface{}) Logger                    { return m }
func (m *MockLogger) Print(v ...interface{})                          {}
func (m *MockLogger) WithField(key string, value interface{}) Logger  { return m }
func (m *MockLogger) WithError(err error) Logger                      { return m }
func (m *MockLogger) SkipCallers(count int) Logger                    { return m }
func (m *MockLogger) Check(level LogLevel) bool                       { return true }
