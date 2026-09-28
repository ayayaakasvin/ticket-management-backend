package metrics

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	// http metrics
	RequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total HTTP requests.",
		},
		[]string{"method", "route", "status"},
	)

	ClientErrorsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_client_errors_total",
			Help: "Total HTTP client errors (4xx).",
		},
		[]string{"method", "route"},
	)

	ServerErrorsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_server_errors_total",
			Help: "Total HTTP server errors (5xx).",
		},
		[]string{"method", "route"},
	)

	ActiveRequests = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_active_requests",
			Help: "Current active HTTP requests.",
		},
	)

	RequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "route"},
	)

	// database metrics
	DBQueries = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "db_queries_total",
			Help: "Executed database queries.",
		},
	)

	ActiveDBConnections = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "db_active_connections",
			Help: "Current active database connections.",
		},
	)

	// worker pool metrics
	TotalJobsProcessed = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "worker_pool_jobs_processed_total",
			Help: "Total jobs processed by the worker pool.",
		},
	)

	TotalJobErrors = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "worker_pool_job_errors_total",
			Help: "Total job errors in the worker pool.",
		},
	)

	ActiveWorkers = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "worker_pool_active_workers",
			Help: "Current active workers in the worker pool.",
		},
	)

	// Filesystem metrics
	TotalFilesSaved = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "filesystem_files_saved_total",
			Help: "Total files saved in the filesystem.",
		},
	)

	TotalFilesDeleted = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "filesystem_files_deleted_total",
			Help: "Total files deleted from the filesystem.",
		},
	)

	ActiveFiles = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "filesystem_active_files",
			Help: "Current active files in the filesystem.",
		},
	)
)

func ObserveDBCall[T any](op func() (T, error)) (T, error) {
	ActiveDBConnections.Inc()
	DBQueries.Inc()
	defer ActiveDBConnections.Dec()

	return op()
}

func ObserveBackgroundJob(op func() error) error {
	TotalJobsProcessed.Inc()
	ActiveWorkers.Inc()
	defer ActiveWorkers.Dec()

	if err := op(); err != nil {
		TotalJobErrors.Inc()
		return err
	}

	return nil
}

var registerOnce sync.Once

func Register() {
	registerOnce.Do(func() {
		collectors := []prometheus.Collector{
			RequestsTotal,
			RequestDuration,
			ActiveRequests,
			DBQueries,
			ClientErrorsTotal,
			ServerErrorsTotal,
			TotalJobsProcessed,
			TotalJobErrors,
			ActiveWorkers,
			TotalFilesSaved,
			TotalFilesDeleted,
			ActiveFiles,
		}

		for _, collector := range collectors {
			if err := prometheus.Register(collector); err != nil {
				if _, ok := err.(prometheus.AlreadyRegisteredError); !ok {
					panic(err)
				}
			}
		}
	})
}
