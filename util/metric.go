package util

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

//service implements Service interface
type service struct {
	httpRequestHistogram *prometheus.HistogramVec
}

type HistogramData struct {
	Handler    string
	Method     string
	StatusCode string
	Duration   float64
}

//NewPrometheusService create a new prometheus service
func NewPrometheusService() (*service, error) {
	http := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "http",
		Name:      "request_duration_seconds",
		Help:      "The latency of the HTTP requests.",
		Buckets:   prometheus.DefBuckets,
	}, []string{"handler", "method", "code"})

	s := &service{
		httpRequestHistogram: http,
	}
	err := prometheus.Register(s.httpRequestHistogram)
	if err != nil && err.Error() != "duplicate metrics collector registration attempted" {
		return nil, err
	}
	return s, nil
}

func (s *service) Handler() http.Handler {
	return promhttp.Handler()
}

//SaveHTTP send metrics to server
func (s *service) SaveHistogram(handler, method, statusCode string, duration float64) {
	s.httpRequestHistogram.WithLabelValues(handler, method, statusCode).Observe(duration)
}
