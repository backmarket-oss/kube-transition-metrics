package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type responseLogger struct {
	responseWriter http.ResponseWriter
	statusCode     int
	bodyBytesSent  int64
}

func (rl *responseLogger) WriteHeader(code int) {
	rl.statusCode = code
	rl.responseWriter.WriteHeader(code)
}

func (rl *responseLogger) Write(data []byte) (int, error) {
	length, err := rl.responseWriter.Write(data)
	rl.bodyBytesSent += int64(length)

	if err != nil {
		err = fmt.Errorf("failed to write response: %w", err)
	}

	return length, err
}

func (rl *responseLogger) Header() http.Header {
	return rl.responseWriter.Header()
}

// ZerologHTTPHandler is a custom request logger middleware.
type ZerologHTTPHandler struct {
	handler http.Handler
}

// NewZerologHTTPHandler creates a new ZerologHTTPHandler middleware.
func NewZerologHTTPHandler(handler http.Handler) *ZerologHTTPHandler {
	return &ZerologHTTPHandler{handler: handler}
}

func (c *ZerologHTTPHandler) logger() *zerolog.Logger {
	logger := log.With().
		Str("subsystem", "http").
		Logger()

	return &logger
}

func (c *ZerologHTTPHandler) ServeHTTP(
	writer http.ResponseWriter,
	req *http.Request,
) {
	start_time := time.Now()
	logger := c.logger()

	response_logger := &responseLogger{responseWriter: writer}
	// Call the next handler in the chain
	c.handler.ServeHTTP(response_logger, req)

	// Log request information
	duration := time.Since(start_time)

	// log_format combined '$remote_addr - $remote_user [$time_local] '
	//                    '"$request" $status $body_bytes_sent '
	//                    '"$http_referer" "$http_user_agent"';
	nanoseconds := float64(time.Nanosecond) / float64(time.Second)
	logger.Info().
		Dict("http", zerolog.Dict().
			Int64("bytes_sent", response_logger.bodyBytesSent).
			Float64("duration", float64(duration.Nanoseconds())*nanoseconds).
			Str("host", req.Host).
			Str("http_referer", req.Referer()).
			Str("method", req.Method).
			Str("path", req.URL.Path).
			Str("remote_address", req.RemoteAddr).
			Str("remote_user", "-").
			Str("request_uri", req.RequestURI).
			Float64("start_time", float64(start_time.UnixNano())*nanoseconds).
			Int("status_code", response_logger.statusCode).
			Str("user_agent", req.UserAgent())).
		Msgf(
			"%s - - [%s] %#v %d %d %#v %#v",
			req.RemoteAddr,
			time.Now().Format("2006-01-02 15:04:05"),
			req.Method+" "+req.RequestURI,
			response_logger.statusCode,
			response_logger.bodyBytesSent,
			req.Referer(),
			req.UserAgent(),
		)
}
