package log

import (
	"errors"
	"fmt"
	"runtime"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// zapLogger is a struct that encapsulates zap's SugaredLogger and custom trace level handling.
type zapLogger struct {
	log        zap.SugaredLogger // The main logger instance for logging.
	traceLevel bool              // Indicates if trace-level logging is enabled.
}

// skipCallers defines the number of stack frames to skip when retrieving caller information.
var skipCallers = 1

// options defines global zap options to set up the logger's behavior, such as caller information.
var options = []zap.Option{
	zap.Development(),
	zap.AddCaller(),
	zap.AddCallerSkip(skipCallers),
}

// convLevel converts a custom LogLevel to a corresponding zapcore.Level.
// Returns nil if the LogLevel is invalid.
func convLevel(level LogLevel) *zapcore.Level {
	var lvl zapcore.Level

	switch level {
	case TraceLevel:
		lvl = zap.DebugLevel
	case DebugLevel:
		lvl = zap.DebugLevel
	case InfoLevel:
		lvl = zap.InfoLevel
	case WarnLevel:
		lvl = zap.WarnLevel
	case ErrorLevel:
		lvl = zap.ErrorLevel
	case FatalLevel:
		lvl = zap.FatalLevel
	default:
		return nil
	}

	return &lvl
}

// newZap creates a new zapLogger instance based on the provided configuration.
// Accepts a boolean for JSON formatting and a LogLevel for severity.
// Returns an error if the LogLevel is invalid.
func newZap(json bool, level LogLevel) (Logger, error) {
	lvl := convLevel(level)

	if lvl == nil {
		return nil, errors.New("wrong logging level")
	}

	config := zap.Config{
		Level:       zap.NewAtomicLevelAt(*lvl),
		Development: true,
		Encoding:    "json",
		OutputPaths: []string{"stdout"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:   "message",
			LevelKey:     "severity",
			TimeKey:      "timestamp",
			CallerKey:    "module",
			EncodeLevel:  zapcore.LowercaseLevelEncoder,
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}

	// Configure logger for console output if JSON formatting is disabled.
	if !json {
		config.Encoding = "console"
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		config.EncoderConfig.TimeKey = ""
		config.EncoderConfig.EncodeCaller = bracketsCallerEncoder
	}

	// Custom handling for TraceLevel logs.
	if level == TraceLevel {
		config.EncoderConfig.EncodeLevel = TraceLevelEncoder
	}

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}
	return &zapLogger{*logger.Sugar(), TraceLevel == level}, nil
}

// TraceLevelEncoder formats trace-level messages distinctly for higher visibility.
func TraceLevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	if l == zapcore.DebugLevel-1 {
		enc.AppendString("TRACE")
		return
	}
	zapcore.CapitalColorLevelEncoder(l, enc)
}

// bracketsCallerEncoder formats the caller path within brackets for enhanced readability.
func bracketsCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + caller.TrimmedPath() + "]:")
}

// newZapSome initializes an unconfigured zapLogger for development purposes.
func newZapSome() *zapLogger {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.StacktraceKey = ""
	config.EncoderConfig.TimeKey = ""
	l, _ := config.Build()
	return &zapLogger{*l.Named("<unconfigured logger>").Sugar(), false}
}

// trace logs a custom trace-level message, with adjustments for caller information.
func trace(l *zapLogger, msg string) {
	skipLogger := l.log.WithOptions(options...)
	const callerSkipOffset = 2
	ce := &zapcore.CheckedEntry{}
	ce = ce.AddCore(zapcore.Entry{}, skipLogger.Desugar().Core())
	if ce != nil {
		ce.Entry.Caller = zapcore.NewEntryCaller(runtime.Caller(callerSkipOffset))
		ce.Entry.Message = msg
		ce.Entry.Level = zapcore.DebugLevel - 1
		ce.Write()
	}
}

// The following methods implement different log levels and formats for zapLogger.
func (l *zapLogger) Info(i ...interface{}) {
	skipLogger := l.log.WithOptions(options...)
	skipLogger.Info(i...)
}

func (l *zapLogger) Infof(s string, i ...interface{}) {
	skipLogger := l.log.WithOptions(options...)
	skipLogger.Infof(s, i...)
}

func (l *zapLogger) Infow(s string, i ...interface{}) {
	skipLogger := l.log.WithOptions(options...)
	skipLogger.Infow(s, i...)
}

func (l *zapLogger) Warn(i ...interface{}) {
	skipLogger := l.log.WithOptions(options...)
	skipLogger.Warn(i...)
}

func (l *zapLogger) Warnf(s string, i ...interface{}) {
	skipLogger := l.log.WithOptions(options...)
	skipLogger.Warnf(s, i...)
}

func (l *zapLogger) Warnw(s string, i ...interface{}) {
	skipLogger := l.log.WithOptions(options...)
	skipLogger.Warnw(s, i...)
}

func (l *zapLogger) Error(i ...interface{}) {
	skipLogger := l.log.WithOptions(options...)
	skipLogger.Error(i...)
}

func (l *zapLogger) Errorf(s string, i ...interface{}) {
	skipLogger := l.log.WithOptions(options...)
	skipLogger.Errorf(s, i...)
}

func (l *zapLogger) Errorw(s string, i ...interface{}) {
	skipLogger := l.log.WithOptions(options...)
	skipLogger.Errorw(s, i...)
}

func (l *zapLogger) Debug(i ...interface{}) {
	skipLogger := l.log.WithOptions(options...)
	skipLogger.Debug(i...)
}

func (l *zapLogger) Debugf(s string, i ...interface{}) {
	skipLogger := l.log.WithOptions(options...)
	skipLogger.Debugf(s, i...)
}

func (l *zapLogger) Trace(s string, i ...interface{}) {
	skipLogger := l.log.WithOptions(options...)
	skipLogger.Debugf(s, i...)
}

func (l *zapLogger) Tracef(s string, i ...interface{}) {
	skipLogger := l.log.WithOptions(options...)
	skipLogger.Debugf(s, i...)
}

func (l *zapLogger) Debugw(s string, i ...interface{}) {
	skipLogger := l.log.WithOptions(options...)
	skipLogger.Debugw(s, i...)
}

func (l *zapLogger) Fatal(i ...interface{}) {
	skipLogger := l.log.WithOptions(options...)
	skipLogger.Fatal(i...)
}

func (l *zapLogger) Fatalf(s string, i ...interface{}) {
	skipLogger := l.log.WithOptions(options...)
	skipLogger.Fatalf(s, i...)
}

func (l *zapLogger) Print(args ...interface{}) {
	trace(l, fmt.Sprint(args...))
}

func (l *zapLogger) Printf(arg string, int ...interface{}) {
	trace(l, fmt.Sprintf(arg, int...))
}

func (l *zapLogger) Panic(args ...interface{}) {
	skipLogger := l.log.WithOptions(options...)
	skipLogger.Panic(args)
}

func (l *zapLogger) Panicf(msg string, args ...interface{}) {
	skipLogger := l.log.WithOptions(options...)
	skipLogger.Panicf(msg, args)
}

// WithError attaches an error message as a context field to the logger.
func (l *zapLogger) WithError(err error) Logger {
	return &zapLogger{*l.log.With("error", err), l.traceLevel}
}

// WithField attaches a key-value pair as a context field to the logger.
func (l *zapLogger) WithField(key string, value interface{}) Logger {
	return &zapLogger{*l.log.With(key, value), l.traceLevel}
}

// SkipCallers configures the logger to skip a specified number of caller stack frames.
func (l *zapLogger) SkipCallers(count int) Logger {
	return &zapLogger{*l.log.Desugar().WithOptions(zap.AddCallerSkip(count)).Sugar(), l.traceLevel}
}

// With adds multiple context fields for structured logging.
func (l *zapLogger) With(f ...interface{}) Logger {
	return &zapLogger{log: *l.log.With(f)}
}

// Check determines if logging should proceed at the specified LogLevel.
func (l *zapLogger) Check(level LogLevel) bool {
	if level == TraceLevel {
		return l.traceLevel
	}

	lvl := convLevel(level)
	// Invalid logging level
	if lvl == nil {
		return false
	}

	return l.log.Desugar().Check(*lvl, "") != nil
}
