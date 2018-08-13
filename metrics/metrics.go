package metrics

import (
	"github.com/oshankkumar/prometheus-exporter"
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

var (
	DefaultReqCounter   = NewRequestCountMetrics()
	DefaultReqLatency   = NewRequestHistogram()
	DefaultReqSummary   = NewRequestSummary()
	DefaultErrorCounter = NewErrorCounterMetrics()
)


func NewRequestCountMetrics() *RequestCounter {
	reqCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Subsystem: prometheus_exporter.ServiceName(),
			Name:      "http_request_count",
			Help:      "total number of http request received",
		}, []string{"method", "path"},
	)
	prometheus.MustRegister(reqCounter)
	return &RequestCounter{reqCounter}
}

type RequestCounter struct {
	counter *prometheus.CounterVec
}

func (rc *RequestCounter) Inc(method string, path string) {
	rc.counter.WithLabelValues(method, path).Inc()
}

func NewRequestHistogram() *RequestHistogram {
	histogram := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Subsystem: prometheus_exporter.ServiceName(),
			Name:      "http_request_latency",
			Help:      "request latency in millisecond ",
		}, []string{"method", "path"},
	)
	prometheus.MustRegister(histogram)
	return &RequestHistogram{histogram}
}

type RequestHistogram struct {
	histogram *prometheus.HistogramVec
}

func (h *RequestHistogram) Observe(method, path string, dur time.Duration) {
	h.histogram.WithLabelValues(method, path).Observe(float64(dur))
}

func NewErrorCounterMetrics() *ErrorCounter {
	err4xxCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Subsystem: prometheus_exporter.ServiceName(),
			Name:      "error_4xx_count",
			Help:      "total number of 4xx error encountered",
		}, []string{"method", "path"},
	)
	err5xxCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Subsystem: prometheus_exporter.ServiceName(),
			Name:      "error_5xx_count",
			Help:      "total number of 5xx error encountered",
		}, []string{"method", "path"},
	)
	prometheus.MustRegister(err4xxCounter, err5xxCounter)
	return &ErrorCounter{err4xxCounter, err5xxCounter}
}

type ErrorCounter struct {
	error4XX *prometheus.CounterVec
	error5XX *prometheus.CounterVec
}

func (e *ErrorCounter) Inc4xx(method string, path string) {
	e.error4XX.WithLabelValues(method, path).Inc()
}

func (e *ErrorCounter) Inc5xx(method string, path string) {
	e.error5XX.WithLabelValues(method, path).Inc()
}

func NewRequestSummary() *RequestSummary {
	summary := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Subsystem: prometheus_exporter.ServiceName(),
			Name:      "http_request_percentile",
			Help:      "request percentile in millisecond ",
		}, []string{"method", "path"},
	)
	prometheus.MustRegister(summary)
	return &RequestSummary{summary: summary}
}

type RequestSummary struct {
	summary *prometheus.SummaryVec
}

func (s *RequestSummary) Observe(method, path string, dur time.Duration) {
	s.summary.WithLabelValues(method, path).Observe(float64(dur))
}
