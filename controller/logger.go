package controller

import (
	"bytes"
	"encoding/json"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// Logger is a wrapper around zap.SugaredLogger
type Logger struct {
	logger *zap.SugaredLogger
}

// NewLogger initializes a new Logger instance
func NewLogger(logger *zap.SugaredLogger) *Logger {
	return &Logger{logger: logger}
}

// responseWriter wraps echo.Response to capture the response body
type responseWriter struct {
	echo.Response
	body *bytes.Buffer
}

// Write captures the response body while continuing to write to the original response
func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)            // Buffer the response body
	return w.Response.Write(b) // Write the response to the client
}

// LoggerMiddleware intercepts the response to log any errors present in the response body
func (rc *Logger) LoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Wrap the original response writer to intercept the response body
		var res = c.Response()
		var bodyBuffer = new(bytes.Buffer)
		var writer = &responseWriter{
			Response: *res,
			body:     bodyBuffer,
		}
		c.Response().Writer = writer

		// Proceed with request handling
		var err = next(c)

		// After the handler, check if there was an error in the response
		var failureResponse FailureResponse
		if json.Unmarshal(writer.body.Bytes(), &failureResponse) == nil {
			// If the response contains an error, log it
			if failureResponse.Error != "" {
				rc.logger.Errorf("Error logged: %v", failureResponse.Error)
			}
		}

		return err
	}
}
