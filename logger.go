package log

import (
	"context"
	"strings"
)

const (
	loggerKey = "logger"
)

// Logger is an interface that defines logging methods with various log levels and formats.
type Logger interface {
	// Info writes an informational message.
	Info(...interface{})
	// Infof writes a formatted informational message.
	Infof(string, ...interface{})
	// Infow writes an informational message with key-value pairs for context.
	Infow(string, ...interface{})
	// Warn writes a warning message.
	Warn(...interface{})
	// Warnf writes a formatted warning message.
	Warnf(string, ...interface{})
	// Warnw writes a warning message with key-value pairs for context.
	Warnw(string, ...interface{})
	// Error writes an error message.
	Error(...interface{})
	// Errorf writes a formatted error message.
	Errorf(string, ...interface{})
	// Errorw writes an error message with key-value pairs for context.
	Errorw(string, ...interface{})
	// Debug writes a debug message.
	Debug(...interface{})
	// Debugf writes a formatted debug message.
	Debugf(string, ...interface{})
	// Debugw writes a debug message with key-value pairs for context.
	Debugw(string, ...interface{})
	// Fatal writes a fatal message and typically triggers application exit.
	Fatal(...interface{})
	// Fatalf writes a formatted fatal message.
	Fatalf(string, ...interface{})
	// With adds fields for structured logging to all subsequent logs.
	With(f ...interface{}) Logger
	// Check returns true if the log level is enabled for the logger instance.
	Check(level LogLevel) bool
	// Print logs a general message without a specific severity.
	Print(v ...interface{})
	// WithField adds a single key-value pair to the Logger instance.
	WithField(key string, value interface{}) Logger
	// WithError attaches an error to the Logger instance for context.
	WithError(err error) Logger
	// SkipCallers skips a specified number of call stack frames for cleaner logs.
	SkipCallers(count int) Logger
}

// LogLevel defines the severity of logs, from Panic (highest) to Trace (lowest).
type LogLevel uint8

var (
	def            Logger          = nil // Global default logger instance
	defaultContext context.Context = nil // Default context with logger settings
)

const (
	// PanicLevel is the highest severity; logs and then panics.
	PanicLevel LogLevel = iota
	// FatalLevel logs and then exits the application.
	FatalLevel
	// ErrorLevel is for errors that require attention.
	ErrorLevel
	// WarnLevel is for non-critical issues that need monitoring.
	WarnLevel
	// InfoLevel is for general operational information.
	InfoLevel
	// DebugLevel is for detailed debugging information.
	DebugLevel
	// TraceLevel is for the most granular level of information.
	TraceLevel
)

// NewLogger creates a new Logger instance based on the provided configuration.
func NewLogger(conf *Config) (Logger, error) {
	return newZap(conf.IsJson, Text2Level(conf.Level))
}

// SetDefaultLogger sets a global Logger instance.
func SetDefaultLogger(l Logger) {
	def = l
}

// SetDefaultContext sets a default context that may include logging configurations.
func SetDefaultContext(ctx context.Context) {
	defaultContext = ctx
}

// GetDefaultLogger returns the global Logger instance or initializes it based on LoggerConfig.
func GetDefaultLogger() Logger {
	if def != nil {
		return def
	}
	if LoggerConfig.Level == "" {
		LoggerConfig.Level = "DEBUG"
	}
	l, err := newZap(LoggerConfig.IsJson, Text2Level(LoggerConfig.Level))
	if err != nil {
		panic(err) // Panic if logger initialization fails
	}
	return l
}

// ToContext attaches a Logger to a given context for retrieval in other parts of the app.
func ToContext(ctx context.Context, l Logger) context.Context {
	return context.WithValue(ctx, loggerKey, l)
}

// FromContext retrieves a Logger from the provided context or falls back to a default logger.
func FromContext(ctx context.Context) Logger {
	var l Logger
	o := ctx.Value(loggerKey)
	if o == nil {
		l = FromDefaultContext()
	} else {
		if loggerFromContext, ok := o.(Logger); ok {
			return loggerFromContext
		} else {
			return nil
		}
	}
	return l
}

// FromDefaultContext returns a Logger instance based on defaultContext settings.
func FromDefaultContext() Logger {
	var l Logger
	// Initialize defaultContext if it is nil
	if defaultContext == nil {
		defaultContext = context.Background()
	}

	// Retrieve Logger from defaultContext or use GetDefaultLogger
	if loggerFromContext, ok := defaultContext.Value(loggerKey).(Logger); ok {
		l = loggerFromContext
	} else {
		l = GetDefaultLogger()
	}

	return l
}

// Text2Level converts a string log level to a LogLevel enum for structured logging.
func Text2Level(level string) LogLevel {
	var logLevel LogLevel
	switch strings.ToUpper(level) {
	case "TRACE":
		logLevel = TraceLevel
	case "DEBUG":
		logLevel = DebugLevel
	case "INFO":
		logLevel = InfoLevel
	case "WARNING":
		logLevel = WarnLevel
	case "ERROR":
		logLevel = ErrorLevel
	case "FATAL":
		logLevel = FatalLevel
	case "PANIC":
		logLevel = PanicLevel
	case "UNKNOWN":
		logLevel = InfoLevel
	}
	return logLevel
}
