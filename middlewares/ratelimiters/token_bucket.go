package ratelimiters

import (
	"net/http"
	"strings"
	"sync"
	"time"
)

type GetIp func(r *http.Request) string
type tokenBucketConfig struct {
	maxCapacity           uint16
	rate                  uint16 // per sec
	GetIp                 GetIp  // Logic to get Ips from requests. if using a proxy service then put your own logic in this function. Do not mutate after starting
	IpStickySession       uint16 // seconds, cleans the idle ips above this threshold
	IpCleanTickerDuration uint16 // seconds, checks every T amount of seconds to remove the idle ips
}

type tokenBucket struct {
	buckets []*user
	Config  tokenBucketConfig
	mu      sync.RWMutex
}

type user struct {
	ip           string
	tokens       uint16
	lastRefillAt time.Time
	lastSeenAt   time.Time
	mu           sync.RWMutex
}

func DefaultGetIp(splitPort bool) GetIp {
	return func(r *http.Request) string {
		if splitPort {
			return strings.Split(r.RemoteAddr, ":")[0]
		} else {
			return r.RemoteAddr
		}
	}
}

// IP based ratelimiting
// logic for getting ip is in tb.Config.GetIp, change the func to your own if you face issues
func NewTokenBucket(maxCapacity, refillRate uint16) *tokenBucket {
	tb:=tokenBucket{
		Config: tokenBucketConfig{
			maxCapacity: maxCapacity,
			rate:        refillRate,
			GetIp:       DefaultGetIp(true),
		},
	}
	tb.cleanupRunner()
	return &tb
}

func (tb *tokenBucket) Allow(r *http.Request) bool {
	ip := tb.Config.GetIp(r)
	allow, ok := tb.findBucket(ip) // if it finds it then auto refills and clears it.
	if ok {
		return allow
	}
	// if ip doesn't exist
	tb.addNewBucket(ip)
	return true
}

func (tb *tokenBucket) addNewBucket(ip string){
	tb.mu.Lock()
	tb.buckets = append(tb.buckets, &user{
		ip:           ip,
		tokens:       tb.Config.maxCapacity,
		lastRefillAt: time.Now(),
		lastSeenAt:   time.Now(),
	})
	tb.mu.Unlock()
}

func (tb *tokenBucket) findBucket(ip string) (bool, bool) {
	tb.mu.RLock()
	defer tb.mu.RUnlock()
	for _, b := range tb.buckets {
		if ip == b.ip {
			tb.refill(b)
			b.mu.Lock()
			defer b.mu.Unlock()
			if b.tokens >= 1 {
				b.tokens--
				return true, true
			}
			return false, true
		}
	}
	return false, false
}
func (tb *tokenBucket) refill(u *user) {
	now := time.Now()
	elapsed := uint16(now.Sub(u.lastRefillAt).Seconds()) // trust me nobody's visiting the api after ~45 days
	u.mu.Lock()
	u.tokens = min(tb.Config.maxCapacity, u.tokens+elapsed*tb.Config.rate) // even if they do, it caps out at the 16 bit limit.
	u.lastRefillAt = now
	u.lastSeenAt = now
	u.mu.Unlock()
}


func (tb *tokenBucket) cleanupRunner() {
	ticker := time.NewTicker(time.Duration(tb.Config.IpCleanTickerDuration) * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		var newBuckets []*user
		tb.mu.Lock()
		for i, b := range tb.buckets {
			idle := time.Since(b.lastSeenAt)
			if idle > (time.Duration(tb.Config.IpStickySession) * time.Second) {
				newBuckets=append(tb.buckets[:i], tb.buckets[i+1:]... )	
			}
		}
		if newBuckets!=nil{
			tb.buckets=newBuckets
		}
		tb.mu.Unlock()
	}
}
