package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const RequestIDKey = "X-Request-ID"

func Logger(base *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		lgr := base.WithGroup("http")

		reqID := c.GetString(RequestIDKey)
		if reqID == "" {
			reqID = c.GetHeader(RequestIDKey)
		}
		if reqID == "" {
			reqID = newRequestID()
		}

		c.Set(RequestIDKey, reqID)
		c.Writer.Header().Set(RequestIDKey, reqID)

		lgr = lgr.With(
			slog.String("path", c.Request.URL.Path),
			slog.String("request_id", reqID),
			slog.String("method", c.Request.Method),
			slog.String("client_ip", c.ClientIP()),
		)

		start := time.Now()
		lgr = lgr.With(slog.Time("start_time", start))

		lgr.Info("request started")

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		lgr = lgr.With(
			slog.Int("status", status),
			slog.Duration("latency", latency),
		)

		ginErrors := c.Errors.ByType(gin.ErrorTypeAny).Errors()
		if len(ginErrors) > 0 {
			msgs := strings.Join(ginErrors, "; ")
			lgr = lgr.With(slog.String("error_messages", msgs))
		}

		ctx := c.Request.Context()

		switch {
		case status >= http.StatusInternalServerError:
			lgr.ErrorContext(ctx, "request failed")
		case status >= http.StatusBadRequest:
			lgr.WarnContext(ctx, "request failed")
		default:
			lgr.InfoContext(ctx, "request completed")
		}
	}
}

func newRequestID() string {
	var b [16]byte
	_, err := rand.Read(b[:])
	if err != nil {
		now := time.Now().Unix()
		return fmt.Sprintf("fallback-%d", now)
	}
	return hex.EncodeToString(b[:])
}
