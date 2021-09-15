package util

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

//MetricService implements MetricService interface
type MetricService struct {
	httpReq *prometheus.HistogramVec
}

//NewPrometheusService create a new prometheus service
func NewPrometheusService() (*MetricService, error) {
	http := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "http",
		Name:      "request_duration_seconds",
		Help:      "The latency of the HTTP requests.",
		Buckets:   prometheus.DefBuckets,
	}, []string{"handler", "method", "code"})

	s := &MetricService{
		httpReq: http,
	}
	err := prometheus.Register(s.httpReq)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *MetricService) Handler() http.Handler {
	return promhttp.Handler()
}

func (s *MetricService) SaveHistogram(handler, method, statusCode string, duration float64) {
	s.httpReq.WithLabelValues(handler, method, statusCode).Observe(duration)
}
