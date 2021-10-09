package rest

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type MetricProvider interface {
	Handler() http.Handler
	SaveHistogram(handler, method, statusCode string, duration float64)
}

type LimiterProvider interface {
	Allow() bool
	GetDefaultError() (error, int)
}

//MakeMetricsHandlers make url handlers
func MakeMetricsHandlers(r *gin.Engine, service MetricProvider) {
	r.GET("/metrics", func(c *gin.Context) {
		service.Handler().ServeHTTP(c.Writer, c.Request)
	})
}

// HistogramMiddleware for histogram metrics
func HistogramMiddleware(service MetricProvider) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		handler := c.Request.RequestURI
		method := c.Request.Method
		c.Next()
		duration := time.Since(start).Seconds()
		service.SaveHistogram(handler, method, http.StatusText(c.Writer.Status()), duration)
	}
}

// LimiterMiddleware for rate limiter
func LimiterMiddleware(limiter LimiterProvider) gin.HandlerFunc {
	return func(c *gin.Context) {
		if limiter.Allow() {
			c.Next()
			return
		}

		err, status := limiter.GetDefaultError()
		c.Error(err)
		c.AbortWithStatus(status)
	}
}
