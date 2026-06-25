package ratelimiters

import (
	"net/http"
	"sync"
	"time"
)

type windowCounterConfig struct {
	maxRequests uint16 // per second

}

type windowCounter struct {
	reqCount uint16
	gapTime  time.Time
	config   windowCounterConfig
	mu       sync.RWMutex
}

// Non-IP based
// Fixed value per second, anything over is rejected.
func NewWindowCounter(maxRequests uint16) *windowCounter {
	return &windowCounter{
		config: windowCounterConfig{
			maxRequests: maxRequests,
		},
	}
}

func (wc *windowCounter) Allow(_ *http.Request) bool {
	if !wc.checkTimeGap() {
		wc.mu.Lock()
		wc.reqCount++
		wc.mu.Unlock()
		wc.mu.RLock()
		defer wc.mu.RUnlock()
		if wc.reqCount > wc.config.maxRequests {
			return false
		} else {
			return true
		}
	} else{
		wc.mu.Lock()
		defer wc.mu.Unlock()
		wc.gapTime=time.Now()
		return true
	}


}

// returns if 1 second has passed or not since the reqCount reset.
func (wc *windowCounter) checkTimeGap()bool{
	now := time.Now()
	wc.mu.RLock()
	defer wc.mu.RUnlock()
	return now.Sub(wc.gapTime) >= time.Second
}