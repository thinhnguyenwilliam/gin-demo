// gin-demo/internal/middleware/logger.go
package middleware

import (
	"bytes"
	"io"
	"mime/multipart"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thinhnguyenwilliam/gin-demo/internal/logger"
)

type bodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b) // capture response
	return w.ResponseWriter.Write(b)
}

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {

		start := time.Now()

		// Wrap response writer
		bw := &bodyWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
		}
		c.Writer = bw

		var bodyLog interface{}

		// ---- READ BODY SAFELY ----
		if c.Request.Body != nil {
			contentType := c.GetHeader("Content-Type")

			switch {
			// JSON
			case contentType == "application/json":
				bodyBytes, _ := io.ReadAll(c.Request.Body)
				bodyLog = string(bodyBytes)

				// Restore body for next handlers
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			// Form
			case contentType == "application/x-www-form-urlencoded":
				c.Request.ParseForm()
				bodyLog = c.Request.Form

			// Multipart
			case contentType == "multipart/form-data":
				err := c.Request.ParseMultipartForm(10 << 20) // 10MB limit
				if err == nil && c.Request.MultipartForm != nil {
					bodyLog = summarizeMultipart(c.Request.MultipartForm)
				}
			}
		}

		// Process request
		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		// Choose log level
		event := logger.Log.Info()
		if status >= 500 {
			event = logger.Log.Error()
		} else if status >= 400 {
			event = logger.Log.Warn()
		}

		event.
			Str("client_ip", c.ClientIP()).
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Int("status", status).
			Dur("latency", latency).
			Interface("request_body", bodyLog).
			Str("response_body", bw.body.String()).
			Msg("request completed")
	}
}

// ---- Helper: avoid logging full files ----
func summarizeMultipart(form *multipart.Form) map[string]interface{} {
	result := make(map[string]interface{})

	result["values"] = form.Value

	files := make(map[string][]string)
	for key, headers := range form.File {
		for _, h := range headers {
			files[key] = append(files[key], h.Filename)
		}
	}

	result["files"] = files

	return result
}
