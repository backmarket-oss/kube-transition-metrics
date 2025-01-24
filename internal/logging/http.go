package logging

import (
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type httpResponseLogger struct {
	responseWriter http.ResponseWriter
	statusCode     int
	bodyBytesSent  int64
}

func (rl *httpResponseLogger) WriteHeader(code int) {
	rl.statusCode = code
	rl.responseWriter.WriteHeader(code)
}

func (rl *httpResponseLogger) Write(data []byte) (int, error) {
	length, err := rl.responseWriter.Write(data)
	rl.bodyBytesSent += int64(length)

	if err != nil {
		err = fmt.Errorf("failed to write response: %w", err)
	}

	return length, err
}

func (rl *httpResponseLogger) Header() http.Header {
	return rl.responseWriter.Header()
}

// HTTPHandler is a custom request logger middleware.
type HTTPHandler struct {
	handler http.Handler
}

// NewHTTPHandler creates a new Handler middleware.
func NewHTTPHandler(handler http.Handler) *HTTPHandler {
	return &HTTPHandler{handler: handler}
}

func (c *HTTPHandler) logger() *zerolog.Logger {
	logger := log.With().
		Str("subsystem", "http").
		Logger()

	return &logger
}

func (c *HTTPHandler) ServeHTTP(
	writer http.ResponseWriter,
	req *http.Request,
) {
	startTime := time.Now()
	logger := c.logger()

	responseLogger := &httpResponseLogger{responseWriter: writer}
	// Call the next handler in the chain
	c.handler.ServeHTTP(responseLogger, req)

	// Log request information
	duration := time.Since(startTime)

	nanoseconds := float64(time.Nanosecond) / float64(time.Second)
	logger.Info().
		Dict("http", zerolog.Dict().
			Int64("bytes_sent", responseLogger.bodyBytesSent).
			Float64("duration", float64(duration.Nanoseconds())*nanoseconds).
			Str("host", req.Host).
			Str("http_referer", req.Referer()).
			Str("method", req.Method).
			Str("path", req.URL.Path).
			Str("remote_address", req.RemoteAddr).
			Str("remote_user", "-").
			Str("request_uri", req.RequestURI).
			Float64("start_time", float64(startTime.UnixNano())*nanoseconds).
			Int("status_code", responseLogger.statusCode).
			Str("user_agent", req.UserAgent())).
		Msgf(
			"%s - - [%s] %#v %d %d %#v %#v",
			req.RemoteAddr,
			time.Now().Format("2006-01-02 15:04:05"),
			req.Method+" "+req.RequestURI,
			responseLogger.statusCode,
			responseLogger.bodyBytesSent,
			req.Referer(),
			req.UserAgent(),
		)
}
