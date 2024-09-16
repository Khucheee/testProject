package controller

import (
	"bufio"
	"bytes"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
)

// customWriter релизует интерфейс типа gin.ResponseWriter, нужен, чтобы перехватить
// тело ответа и залогировать его
type customWriter struct {
	originalWriter gin.ResponseWriter
	responseBody   *bytes.Buffer
	statusCode     int
}

func (writer *customWriter) Write(b []byte) (int, error) {
	writer.responseBody.Write(b)
	return writer.originalWriter.Write(b)
}

func (writer *customWriter) WriteHeader(statusCode int) {
	writer.statusCode = statusCode
	writer.originalWriter.WriteHeader(statusCode)
}

func (writer *customWriter) Header() http.Header {
	return writer.originalWriter.Header()
}

func (writer *customWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return writer.originalWriter.Hijack()
}

func (writer *customWriter) Flush() {
	writer.originalWriter.Flush()
}

func (writer *customWriter) CloseNotify() <-chan bool {
	return writer.originalWriter.CloseNotify()
}

func (writer *customWriter) Status() int {
	return writer.originalWriter.Status()
}

func (writer *customWriter) Size() int {
	return writer.originalWriter.Size()
}

func (writer *customWriter) WriteString(s string) (int, error) {
	return writer.originalWriter.WriteString(s)
}

func (writer *customWriter) Written() bool {
	return writer.originalWriter.Written()
}

func (writer *customWriter) WriteHeaderNow() {
	writer.originalWriter.WriteHeaderNow()
}

func (writer *customWriter) Pusher() http.Pusher {
	return writer.originalWriter.Pusher()
}
