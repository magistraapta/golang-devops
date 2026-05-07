package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/magistraapta/golang-devops/internal/metrics"
)

func HTTPMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		duration := time.Since(start).Seconds()
		method := c.Request.Method
		endpoint := c.FullPath() // ✅ FullPath() preferred over URL.Path for route patterns
		if endpoint == "" {
			endpoint = c.Request.URL.Path // fallback for unmatched routes
		}
		statusCode := strconv.Itoa(c.Writer.Status())

		metrics.HTTPRequestsTotal.WithLabelValues(method, endpoint, statusCode).Inc()
		metrics.HTTPRequestDuration.WithLabelValues(method, endpoint, statusCode).Observe(duration)
	}
}
