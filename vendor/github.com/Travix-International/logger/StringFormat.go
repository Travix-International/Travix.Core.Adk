package logger

import (
	"bytes"
	"fmt"
)

type StringFormat struct {
	levelFormat       string
	eventFormat       string
	messageFormat     string
	metaKeyFormat     string
	metaMessageFormat string
	lineSuffix        string
}

func NewStringFormat(levelFormat, eventFormat, messageFormat, metaKeyFormat, metaMessageFormat, lineSuffix string) *StringFormat {
	if len(levelFormat) == 0 {
		return DefaultStringFormat
	}

	return &StringFormat{
		levelFormat,
		eventFormat,
		messageFormat,
		metaKeyFormat,
		metaMessageFormat,
		lineSuffix,
	}
}

var DefaultStringFormat = &StringFormat{}

func (s *StringFormat) Format(e *Entry) (out string, err error) {
	if e == nil {
		return "", nil
	}

	if s == nil || len(s.levelFormat) == 0 {
		return e.String(), nil
	}

	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf(s.levelFormat, e.Level))
	buf.WriteString(fmt.Sprintf(s.eventFormat, e.Event))
	buf.WriteString(fmt.Sprintf(s.messageFormat, e.Message))

	if len(e.Meta) > 0 {
		for k, v := range e.Meta {
			buf.WriteString(fmt.Sprintf(s.metaKeyFormat, k))
			buf.WriteString(fmt.Sprintf(s.metaMessageFormat, v))
		}
	}

	buf.Write([]byte(s.lineSuffix))

	return buf.String(), nil
}
