package log

// Config defines the logging configuration structure.
// Level sets the logging level (e.g., "info", "debug", "error").
// IsJson toggles between JSON format (true) or plain text format (false) for log output.
type Config struct {
	Level  string // Level defines the logging severity (e.g., "info", "debug").
	IsJson bool   // IsJson determines if the log output should be in JSON format.
}

// LoggerConfig holds the global logging configuration instance.
// This can be modified to set the desired logging settings across the application.
var LoggerConfig = Config{}
