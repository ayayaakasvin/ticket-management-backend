package inner

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type RateLimiter struct {
	mu       sync.Mutex
	limiters map[uint]*rate.Limiter // userID → limiter
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{limiters: make(map[uint]*rate.Limiter)}
}

func (rl *RateLimiter) GetLimiter(userID uint) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	lim, exists := rl.limiters[userID]
	if !exists {
		// allow 1 upload every 3s, burst of 2
		lim = rate.NewLimiter(rate.Every(3*time.Second), 2)
		rl.limiters[userID] = lim
	}
	return lim
}
