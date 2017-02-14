package logger

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
)

type jsonEntry struct {
	Level   string            `json:"level"`
	Event   string            `json:"event"`
	Message string            `json:"message"`
	Meta    map[string]string `json:"meta,omitempty"`

	levelKey   string
	eventKey   string
	messageKey string
}

var (
	ErrNoLevel   = errors.New("No Level!")
	ErrNoEvent   = errors.New("No Event!")
	ErrNoMessage = errors.New("No Message!")
)

func (j jsonEntry) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer

	if len(j.Level) == 0 {
		return nil, ErrNoLevel
	}

	if len(j.Event) == 0 {
		return nil, ErrNoEvent
	}

	if len(j.Message) == 0 {
		return nil, ErrNoMessage
	}

	_, err := buf.WriteString(fmt.Sprintf("{\"%s\":\"%s\", \"%s\":\"%s\", \"%s\":\"%s\"", j.levelKey, j.Level, j.eventKey, j.Event, j.messageKey, j.Message))

	if err != nil {
		return nil, err
	}

	if len(j.Meta) > 0 {
		_, err = buf.WriteString(",")

		if err != nil {
			return nil, err
		}

		i := len(j.Meta)
		for k, v := range j.Meta {
			_, err = buf.WriteString(fmt.Sprintf(" \"%s\":\"%s\"", k, v))

			if err != nil {
				return nil, err
			}

			i = i - 1

			if i > 0 {
				_, err = buf.WriteString(",")

				if err != nil {
					return nil, err
				}
			}
		}
	}

	_, err = buf.WriteString("}")

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type JSONFormat struct {
	levelKey   string
	eventKey   string
	messageKey string
}

func NewJSONFormat() *JSONFormat {
	f := &JSONFormat{}

	f.levelKey = "level"
	f.eventKey = "event"
	f.messageKey = "message"

	return f
}

var DefaultJSONFormat = &JSONFormat{}

func (j *JSONFormat) SetLevelKey(levelKey string) *JSONFormat {
	j.levelKey = levelKey

	return j
}

func (j *JSONFormat) SetEventKey(eventKey string) *JSONFormat {
	j.eventKey = eventKey

	return j
}

func (j *JSONFormat) SetMessageKey(messageKey string) *JSONFormat {
	j.messageKey = messageKey

	return j
}

func (j *JSONFormat) Format(e *Entry) (string, error) {
	if e == nil {
		return "", nil
	}

	if j == nil {
		return "JSONFormat not initialized", nil
	}

	item := jsonEntry{
		Level:   e.Level,
		Event:   e.Event,
		Message: e.Message,
		Meta:    e.Meta,

		levelKey:   j.levelKey,
		eventKey:   j.eventKey,
		messageKey: j.messageKey,
	}

	out, err := json.Marshal(item)

	return string(out), err
}
