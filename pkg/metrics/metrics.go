package metrics

import (
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	RequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "password_validation_requests_total",
			Help: "Total number of password validation requests",
		},
		[]string{"result"},
	)

	ValidationErrorsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "password_validation_errors_total",
			Help: "Total number of validation errors by rule",
		},
		[]string{"rule"},
	)

	RequestDuration = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "password_validation_duration_seconds",
			Help:    "Duration of password validation requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
	)

	InProgress = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "password_validation_in_progress",
			Help: "Number of password validations currently in progress",
		},
	)
)

func RecordValidation(isValid bool, errors []string) {
	if isValid {
		RequestsTotal.WithLabelValues("valid").Inc()
	} else {
		RequestsTotal.WithLabelValues("invalid").Inc()

		for _, err := range errors {
			rule := extractRuleFromError(err)
			ValidationErrorsTotal.WithLabelValues(rule).Inc()
		}
	}
}

func extractRuleFromError(errMsg string) string {
	switch {
	case strings.Contains(errMsg, "characters"):
		return "min_length"
	case strings.Contains(errMsg, "digit"):
		return "digit"
	case strings.Contains(errMsg, "lowercase"):
		return "lowercase"
	case strings.Contains(errMsg, "uppercase"):
		return "uppercase"
	case strings.Contains(errMsg, "special character"):
		return "special_char"
	case strings.Contains(errMsg, "repeated"):
		return "no_duplicates"
	case strings.Contains(errMsg, "whitespace"):
		return "no_whitespace"
	default:
		return "unknown"
	}
}
