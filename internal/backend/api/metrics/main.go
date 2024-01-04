package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

const (
	labelNameHandler                            = "handler"
	labelNameCode                               = "code"
	labelNameMethod                             = "method"
	metricNameInFlightRequests                  = "in_flight_requests"
	metricDescriptionInFlightRequests           = "A gauge of requests currently being served"
	metricNameHTTPRequestsTotal                 = "http_requests_total"
	metricDescriptionHTTPRequestsTotal          = "The total number of incoming http requests"
	metricNameHTTPRequestDurationSeconds        = "http_request_duration_seconds"
	metricDescriptionHTTPRequestDurationSeconds = "The http request duration in seconds"
)

// A Middleware for prometheus metrics
type Middleware struct {
	inflightRequests prometheus.Gauge
	requestsTotal    *prometheus.CounterVec
	requestsDuration *prometheus.HistogramVec
}

// NewMiddleware creates a new metrics middleware
func NewMiddleware() *Middleware {
	return &Middleware{
		inflightRequests: promauto.NewGauge(prometheus.GaugeOpts{
			Name: metricNameInFlightRequests,
			Help: metricDescriptionInFlightRequests,
		}),
		requestsTotal: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: metricNameHTTPRequestsTotal,
			Help: metricDescriptionHTTPRequestsTotal,
		}, []string{labelNameCode, labelNameMethod}),
		requestsDuration: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name: metricNameHTTPRequestDurationSeconds,
			Help: metricDescriptionHTTPRequestDurationSeconds,
		}, []string{labelNameHandler, labelNameMethod}),
	}
}

// Decorate the given handlerFunction to include request metrics
func (mmw *Middleware) Decorate(handlerName string, handlerFunction http.HandlerFunc) http.Handler {

	chain := promhttp.InstrumentHandlerInFlight(mmw.inflightRequests,
		promhttp.InstrumentHandlerDuration(mmw.requestsDuration.MustCurryWith(prometheus.Labels{labelNameHandler: handlerName}),
			promhttp.InstrumentHandlerCounter(mmw.requestsTotal, handlerFunction)))
	return chain
}
