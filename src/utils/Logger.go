package utils

import (
	"io"

	log "github.com/sirupsen/logrus"
)

// LoggerGenerator is a tool function to generate *log.Entry
func LoggerGenerator(formatter log.Formatter, output io.Writer, level log.Level) *log.Entry {
	entry := log.NewEntry(&log.Logger{})

	entry.Logger.SetLevel(level)
	entry.Logger.SetFormatter(formatter)
	entry.Logger.SetOutput(output)
	return entry
}
