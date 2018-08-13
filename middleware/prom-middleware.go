package middleware

import (
	"github.com/emicklei/go-restful"
	"github.com/oshankkumar/prometheus-exporter/metrics"
	"net/http"
	"time"
)

func Prometheus(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := restful.NewResponse(w)
		metrics.DefaultReqCounter.Inc(r.Method, r.URL.Path)
		defer func(begin time.Time) {
			metrics.DefaultReqLatency.Observe(r.Method, r.URL.Path, time.Since(begin)/time.Millisecond)
			metrics.DefaultReqSummary.Observe(r.Method, r.URL.Path, time.Since(begin)/time.Millisecond)
			if resp.StatusCode() >= 400 && resp.StatusCode() < 500 {
				metrics.DefaultErrorCounter.Inc4xx(r.Method, r.URL.Path)
			} else if resp.StatusCode() >= 500 && resp.StatusCode() < 600 {
				metrics.DefaultErrorCounter.Inc5xx(r.Method, r.URL.Path)
			}
		}(time.Now())
		next.ServeHTTP(resp, r)
	})
}
