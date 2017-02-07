package logger

import (
	"bytes"
	"fmt"
)

type Entry struct {
	Level   string
	Event   string
	Message string
	Meta    map[string]string
}

func (e *Entry) String() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("[%s] %s: %s", e.Level, e.Event, e.Message))

	if len(e.Meta) > 0 {
		for k, v := range e.Meta {
			buf.WriteString(fmt.Sprintf("\n%s: %s", k, v))
		}
	}

	return buf.String()
}

func NewEntry(level string, event string, message string, meta map[string]string) *Entry {
	entry := &Entry{
		Level:   level,
		Event:   event,
		Message: message,
		Meta:    meta,
	}

	return entry
}
