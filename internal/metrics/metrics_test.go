package metrics

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestRegisterIsIdempotent(t *testing.T) {
	Register()
	Register()
}

func TestObserveDBCallTracksQueryAndConnectionMetrics(t *testing.T) {
	before := testutil.ToFloat64(DBQueries)
	beforeActive := testutil.ToFloat64(ActiveDBConnections)

	value, err := ObserveDBCall(func() (int, error) {
		return 42, nil
	})
	if err != nil {
		t.Fatalf("ObserveDBCall returned unexpected error: %v", err)
	}
	if value != 42 {
		t.Fatalf("ObserveDBCall returned %d, want 42", value)
	}

	if got := testutil.ToFloat64(DBQueries); got != before+1 {
		t.Fatalf("db queries metric = %v, want %v", got, before+1)
	}
	if got := testutil.ToFloat64(ActiveDBConnections); got != beforeActive {
		t.Fatalf("active db connections metric = %v, want %v", got, beforeActive)
	}
}

func TestObserveBackgroundJobTracksWorkerPoolMetrics(t *testing.T) {
	beforeProcessed := testutil.ToFloat64(TotalJobsProcessed)
	beforeErrors := testutil.ToFloat64(TotalJobErrors)
	beforeActive := testutil.ToFloat64(ActiveWorkers)

	if err := ObserveBackgroundJob(func() error { return nil }); err != nil {
		t.Fatalf("ObserveBackgroundJob returned unexpected error: %v", err)
	}

	if got := testutil.ToFloat64(TotalJobsProcessed); got != beforeProcessed+1 {
		t.Fatalf("processed jobs metric = %v, want %v", got, beforeProcessed+1)
	}
	if got := testutil.ToFloat64(TotalJobErrors); got != beforeErrors {
		t.Fatalf("job errors metric = %v, want %v", got, beforeErrors)
	}
	if got := testutil.ToFloat64(ActiveWorkers); got != beforeActive {
		t.Fatalf("active workers metric = %v, want %v", got, beforeActive)
	}
}
