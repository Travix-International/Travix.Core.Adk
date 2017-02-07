package logger

import (
	"io"
)

type FilterFunc func(e *Entry) bool

type Transport struct {
	destination io.Writer
	formatter   Formatter
	filter      FilterFunc
}

type Formatter interface {
	Format(entry *Entry) (string, error)
}

func NewTransport(w io.Writer, f Formatter) *Transport {
	return &Transport{
		w,
		f,
		FilterAllowAll(),
	}
}

// SetFilter chnages the message filtering to the given filter function.
func (t *Transport) SetFilter(fn FilterFunc) {
	t.filter = fn
}

func (t *Transport) log(e *Entry) (err error) {
	whatToWrite, err := t.formatter.Format(e)

	if err != nil {
		return err
	}

	_, err = t.destination.Write([]byte(whatToWrite))

	return err
}

// FilterAllowAll returns a function that will allow any entry, regardless of its contents
func FilterAllowAll() FilterFunc {
	return func(e *Entry) bool {
		return true
	}
}
