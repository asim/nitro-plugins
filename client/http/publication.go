package http

import (
	"github.com/micro/go-micro/client"
)

type httpPublication struct {
	topic       string
	contentType string
	message     interface{}
}

func newHTTPPublication(topic string, message interface{}, contentType string) client.Publication {
	return &httpPublication{
		message:     message,
		topic:       topic,
		contentType: contentType,
	}
}

func (h *httpPublication) ContentType() string {
	return h.contentType
}

func (h *httpPublication) Topic() string {
	return h.topic
}

func (h *httpPublication) Message() interface{} {
	return h.message
}
