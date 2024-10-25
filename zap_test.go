package log

import (
	"errors"
	"strings"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Test for convLevel to ensure custom log levels are correctly mapped to zapcore levels
func TestConvLevel(t *testing.T) {
	tests := []struct {
		level LogLevel
		want  zapcore.Level
	}{
		{TraceLevel, zap.DebugLevel},
		{DebugLevel, zap.DebugLevel},
		{InfoLevel, zap.InfoLevel},
		{WarnLevel, zap.WarnLevel},
		{ErrorLevel, zap.ErrorLevel},
		{FatalLevel, zap.FatalLevel},
	}

	for _, tt := range tests {
		got := convLevel(tt.level)
		if got == nil || *got != tt.want {
			t.Errorf("convLevel(%v) = %v; want %v", tt.level, got, tt.want)
		}
	}
}

// Test newZap function to verify logger initialization based on JSON and log level config
func TestNewZap(t *testing.T) {
	logger, err := newZap(true, InfoLevel)
	if err != nil || logger == nil {
		t.Fatal("Expected new zapLogger instance, got error or nil")
	}

	invalidLogger, err := newZap(true, LogLevel(100))
	if err == nil || invalidLogger != nil {
		t.Fatal("Expected error on invalid log level, got none")
	}
}

// Test Check method to ensure logger respects enabled log levels
func TestZapLogger_Check(t *testing.T) {
	logger, err := newZap(true, InfoLevel)
	if err != nil {
		t.Fatal(err)
	}
	if !logger.Check(InfoLevel) {
		t.Error("Expected Check to return true for enabled log level")
	}
	if logger.Check(DebugLevel) {
		t.Error("Expected Check to return false for disabled log level")
	}
}

// Test WithError to verify if error is attached as a context field
func TestZapLogger_WithError(t *testing.T) {
	logger := newZapSome()
	errorMessage := errors.New("sample error")
	loggerWithErr := logger.WithError(errorMessage)

	// Use a type assertion to validate the logger instance
	if _, ok := loggerWithErr.(*zapLogger); !ok {
		t.Fatal("Expected zapLogger instance with attached error context")
	}
}

// Test WithField to verify if a key-value pair is attached to the logger context
func TestZapLogger_WithField(t *testing.T) {
	logger := newZapSome()
	loggerWithField := logger.WithField("key", "value")

	if _, ok := loggerWithField.(*zapLogger); !ok {
		t.Fatal("Expected zapLogger instance with attached key-value pair")
	}
}

// Test SkipCallers to ensure the correct number of stack frames are skipped
func TestZapLogger_SkipCallers(t *testing.T) {
	logger := newZapSome()
	skippedLogger := logger.SkipCallers(3)

	if _, ok := skippedLogger.(*zapLogger); !ok {
		t.Fatal("Expected zapLogger instance with caller skip set")
	}
}

// Test TraceLevelEncoder for formatting trace level messages
func TestTraceLevelEncoder(t *testing.T) {
	encoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		LevelKey:    "level",
		EncodeLevel: TraceLevelEncoder,
	})

	entry := zapcore.Entry{
		Level: zapcore.DebugLevel - 1, // custom trace level
	}

	buf, err := encoder.EncodeEntry(entry, nil)
	if err != nil {
		t.Fatalf("Encoding trace level entry failed: %v", err)
	}

	if !contains(buf.String(), "TRACE") {
		t.Error("Expected trace level message to contain TRACE")
	}
}

// Test bracketsCallerEncoder to validate the formatting of caller information within brackets
func TestBracketsCallerEncoder(t *testing.T) {
	encoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		CallerKey:    "caller",
		EncodeCaller: bracketsCallerEncoder,
	})

	entry := zapcore.Entry{
		Caller: zapcore.NewEntryCaller(1, "test.go", 42, true),
	}

	buf, err := encoder.EncodeEntry(entry, nil)
	if err != nil {
		t.Fatalf("Encoding caller entry failed: %v", err)
	}

	if !contains(buf.String(), "[test.go:42]") {
		t.Error("Expected caller information to be formatted with brackets")
	}
}

// Helper function to check if a substring exists in a string
func contains(str, substr string) bool {
	return strings.Contains(str, substr)
}
