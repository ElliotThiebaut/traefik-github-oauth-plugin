package traefik_github_oauth_server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/MuXiu1997/traefik-github-oauth-plugin/internal/app/traefik-github-oauth-server/model"
	"github.com/MuXiu1997/traefik-github-oauth-plugin/internal/pkg/constant"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// NewApiSecretKeyMiddleware returns a middleware that checks the api secret key.
func NewApiSecretKeyMiddleware(apiSecretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(apiSecretKey) == 0 {
			c.Next()
			return
		}
		reqSecretKey := c.GetHeader(constant.HTTP_HEADER_AUTHORIZATION)
		if reqSecretKey != fmt.Sprintf("%s %s", constant.AUTHORIZATION_PREFIX_TOKEN, apiSecretKey) {
			c.JSON(http.StatusUnauthorized, model.ResponseError{
				Message: "Unauthorized",
			})
		}
		c.Next()
	}
}

// NewLoggerMiddleware returns a middleware that logs the request.
func NewLoggerMiddleware(logger *zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Stop timer
		end := time.Now()

		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		if raw != "" {
			path = path + "?" + raw
		}

		logger.Debug().
			Str("module", "gin").
			Time("time_start", start).
			Int("status", statusCode).
			TimeDiff("took", end, start).
			Str("client_ip", clientIP).
			Str("method", method).
			Str("path", path).
			Func(func(e *zerolog.Event) {
				if errorMessage != "" {
					e.Str("error", errorMessage)
				}
			}).
			Msg("")
	}
}
