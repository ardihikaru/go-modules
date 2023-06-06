package loglevel

import "strings"

// maps log level
const (
	Debug string = "debug"
	Info  string = "info"
	Warn  string = "warn"
	Error string = "error"
	Fatal string = "fatal"
	Panic string = "panic"
)

// GetLogLevelMap returns a boolean value to verify if the log level exists
func GetLogLevelMap() map[string]bool {
	return map[string]bool{
		Debug: true,
		Info:  true,
		Warn:  true,
		Error: true,
		Fatal: true,
		Panic: true,
	}
}

func DebugMode(logLevel string) bool {
	return strings.ToLower(logLevel) == Debug
}
