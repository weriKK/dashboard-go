package loggermw

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// TODO: Log before and after ServeHTTP, that way errors that happen during
// 		 processing the request are easier to identify. (maybe also use a
// 		 correlation id for all logs that belong to a request?)

// HandlerFunc is a middleware adding request and response logging
func HandlerFunc(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		crw := newCustomResponseWriter(w)
		next.ServeHTTP(crw, r)

		// https://en.wikipedia.org/wiki/Common_Log_Format + response time
		logrus.Infof("%s %s %s [%s] \"%s %s %s\" %d %d %s",
			r.RemoteAddr,
			"-",
			"-",
			start.Format("02/Jan/2006:15:04:05 -0700"),
			r.Method,
			r.RequestURI,
			r.Proto,
			crw.status,
			crw.size,
			time.Since(start),
		)
	})
}

type customResponseWriter struct {
	http.ResponseWriter
	status int
	size   int
}

// ResponseWriter.Header is kept as it is

// Overrides ResponseWriter.WriteHeader
func (c *customResponseWriter) WriteHeader(status int) {
	c.status = status
	c.ResponseWriter.WriteHeader(status)
}

// Overrides ResponseWriter.Write
func (c *customResponseWriter) Write(b []byte) (int, error) {
	size, err := c.ResponseWriter.Write(b)
	c.size += size
	return size, err
}

func (c *customResponseWriter) Flush() {
	if f, ok := c.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

func newCustomResponseWriter(w http.ResponseWriter) *customResponseWriter {
	// Default status is 200 OK, it's a safe assumption when WriteHeader is
	// never called
	return &customResponseWriter{
		ResponseWriter: w,
		status:         200,
	}
}
