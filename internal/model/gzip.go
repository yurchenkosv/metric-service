package model

import (
	"io"
	"net/http"
)

// GzipWriter is custom type to realize writer interface of ResponseWriter and Writer
type GzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w GzipWriter) Write(b []byte) (int, error) {
	// w.Writer responsible for compression
	return w.Writer.Write(b)
}
