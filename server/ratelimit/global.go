package ratelimit

type FixedWindowCounterConfig struct { 
	MaxRequests uint16 // per sec
}

type FixedWindowCounter struct {
	Config  FixedWindowCounterConfig
	counter uint16
}