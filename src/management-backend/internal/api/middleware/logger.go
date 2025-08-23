package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func Logger(log *zap.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			// Execute handler
			err := next(c)

			stop := time.Now()
			latency := stop.Sub(start)

			// Log with zap fields
			fields := []zap.Field{
				zap.String("request_id", c.Response().Header().Get(echo.HeaderXRequestID)), // needs middleware.RequestID
				zap.String("remote_ip", c.RealIP()),
				zap.String("host", c.Request().Host),
				zap.String("method", c.Request().Method),
				zap.String("uri", c.Request().RequestURI),
				zap.String("user_agent", c.Request().UserAgent()),
				zap.Int("status", c.Response().Status),
				zap.Int64("latency", latency.Nanoseconds()),
				zap.String("latency_human", latency.String()),
				zap.Int64("bytes_in", c.Request().ContentLength),
				zap.Int64("bytes_out", c.Response().Size),
			}

			if err != nil {
				fields = append(fields, zap.String("error", err.Error()))
			} else {
				fields = append(fields, zap.String("error", ""))
			}

			log.Info("request", fields...)

			return err
		}
	}
}
