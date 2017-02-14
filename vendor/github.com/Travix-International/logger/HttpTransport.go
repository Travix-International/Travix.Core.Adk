package logger

import (
	"bytes"
	"errors"
	"net/http"
	"time"
)

type httpTransport struct {
	client *http.Client
	url    string
	method string
}

var ErrNilClient = errors.New("Nil HttpClient!")

func (h *httpTransport) Write(p []byte) (n int, err error) {
	if h == nil {
		return 0, ErrNilClient
	}

	req, err := http.NewRequest(h.method, h.url, bytes.NewReader(p))

	if err != nil {
		return 0, err
	}

	_, err = h.client.Do(req)

	if err != nil {
		return 0, err
	}

	return len(p), nil
}

func NewHttpTransport(url string, format Formatter) *Transport {
	ht := &httpTransport{
		client: &http.Client{
			Timeout: time.Second * 10,
		},
		url:    url,
		method: "POST",
	}

	return NewTransport(ht, format)
}
