package inner_test

import(
	"github.com/ayayaakasvin/oneflick-ticket/internal/models/inner"
	"golang.org/x/time/rate"
	"testing"
)

func TestRateLimiter_GetLimiter(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		userID uint
		want   *rate.Limiter
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rl := inner.NewRateLimiter()
			got := rl.GetLimiter(tt.userID)
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("GetLimiter() = %v, want %v", got, tt.want)
			}
		})
	}
}
