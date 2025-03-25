//go:build unit
// +build unit

package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	io_prometheus_client "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"

	"github.com/dexlabsio/garlic/pkg/monitoring"
)

func TestTrafficMonitoring(t *testing.T) {
	monitoring.TrafficMetric.Reset()
	prometheus.Unregister(monitoring.TrafficMetric)
	prometheus.MustRegister(monitoring.TrafficMetric)

	// Create a test HTTP request and response recorder
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	// Create a test HTTP handler that will be wrapped by the middleware
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("OK")); err != nil {
			panic(err)
		}
	})

	// Wrap the testHandler with the MetricsMonitor middleware
	handler := MetricsMonitor(testHandler)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rec.Code)
	}

	collector, err := monitoring.TrafficMetric.GetMetricWithLabelValues("GET", "unknown", "200")
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, 1.0, testutil.ToFloat64(collector))
}

func TestActiveRequestsMonitoring(t *testing.T) {
	monitoring.ActiveRequests.Reset()
	prometheus.Unregister(monitoring.ActiveRequests)
	prometheus.MustRegister(monitoring.ActiveRequests)

	// Create a test HTTP request and response recorder
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	// Create a test HTTP handler that will be wrapped by the middleware
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		activeRequestsGauge, err := monitoring.ActiveRequests.GetMetricWithLabelValues("GET", "unknown")
		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, 1.0, testutil.ToFloat64(activeRequestsGauge))

		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("OK")); err != nil {
			panic(err)
		}
	})

	// Wrap the testHandler with the MetricsMonitor middleware
	handler := MetricsMonitor(testHandler)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rec.Code)
	}

	labeledCollector, err := monitoring.ActiveRequests.GetMetricWithLabelValues("GET", "unknown")
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, 0.0, testutil.ToFloat64(labeledCollector))
}

func TestLatencyMonitoring(t *testing.T) {
	monitoring.LatencyMetric.Reset()
	prometheus.Unregister(monitoring.LatencyMetric)
	prometheus.MustRegister(monitoring.LatencyMetric)

	// Create a test HTTP request and response recorder
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	// Create a test HTTP handler that will be wrapped by the middleware
	requestRunTime := 1.1
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Duration(requestRunTime) * time.Second)
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("OK")); err != nil {
			panic(err)
		}
	})

	histBeforeRequest := extractHistogramFromGatherer(t, "http_request_duration_seconds")

	// Wrap the testHandler with the MetricsMonitor middleware
	handler := MetricsMonitor(testHandler)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rec.Code)
	}

	histAfterRequest := extractHistogramFromGatherer(t, "http_request_duration_seconds")

	histDiff := subtractMaps(histAfterRequest, histBeforeRequest)

	// Check the histogram values against expected results
	expectedHistDiff := map[float64]uint64{
		0.005: 0,
		0.01:  0,
		0.025: 0,
		0.05:  0,
		0.1:   0,
		0.25:  0,
		0.5:   0,
		1:     0,
		2.5:   1, // Request runtime falls in this bucket and on
		5:     1,
		10:    1, // +Inf bucket will always have the count of requests
	}

	assert.Equal(t, expectedHistDiff, histDiff)
}

// extractHistogramFromGatherer extracts the histrogram of a target metric in the default gatherer into a map
func extractHistogramFromGatherer(t *testing.T, target string) map[float64]uint64 {
	metricFamilies, err := prometheus.DefaultGatherer.Gather()
	if err != nil {
		t.Error(err)
	}

	latencyHistogramValues := parseHistogramCollection(metricFamilies, target)
	return latencyHistogramValues
}

// subtractMaps subtracts map b from map a values iteractively and returns the result.
func subtractMaps(a map[float64]uint64, b map[float64]uint64) map[float64]uint64 {
	if b == nil {
		return a
	}

	result := make(map[float64]uint64)

	for key, valA := range a {
		valB, exists := b[key]

		if exists {
			result[key] = valA - valB
		} else {
			result[key] = valA
		}
	}

	return result
}

// parseHistogramCollection extracts histogram values from the collected metric families
func parseHistogramCollection(metricFamilies []*io_prometheus_client.MetricFamily, target string) map[float64]uint64 {
	histogram := make(map[float64]uint64)

	for _, metricFamily := range metricFamilies {
		if *metricFamily.Name == target {

			for _, metric := range metricFamily.Metric {
				for _, bucket := range metric.GetHistogram().Bucket {

					if count, ok := histogram[*bucket.UpperBound]; ok {
						histogram[*bucket.UpperBound] = count + *bucket.CumulativeCount
					} else {
						histogram[*bucket.UpperBound] = *bucket.CumulativeCount
					}
				}
			}
		}
	}

	return histogram
}
