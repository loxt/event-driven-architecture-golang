package queue

import (
	"net/http"
)

type ResponseWriter struct {
	body       []byte
	statusCode int
	header     http.Header
}

func NewQueueResponseWriter() *ResponseWriter {
	return &ResponseWriter{
		header: http.Header{},
	}
}

func (w *ResponseWriter) Header() http.Header {
	return w.header
}

func (w *ResponseWriter) Write(b []byte) (int, error) {
	w.body = b
	return 0, nil
}

func (w *ResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
}

var okFn = func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
