package xmiddleware

import (
	"errors"

	"github.com/brunoluiz/go-lab/core/app"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
)

func ErrorHandler(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) == 0 {
			c.Next()
			return
		}

		l := log.
			With("path", c.FullPath()).
			With("client_ip", c.ClientIP())

		for _, gerr := range c.Errors {
			var appErr app.Err
			if ok := errors.As(gerr, &appErr); ok {
				switch appErr.Code() {
				case app.ErrCodeNotFound:
					l.Warn(appErr.Error())
					c.JSON(404, map[string]string{"message": appErr.Error()})
					return
				}
			}
		}

		l.ErrorContext(c.Request.Context(),
			"Unknown unmapped error",
			"error", c.Errors,
		)
		c.JSON(500, map[string]string{"message": "Internal error"})
	}
}
