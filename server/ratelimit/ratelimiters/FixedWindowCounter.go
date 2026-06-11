package ratelimiters

type FixedWindowCounterConfig struct { 
	MaxRequests uint16 // per sec
}

func Sig(){

}

type FixedWindowCounter struct {
	Config  FixedWindowCounterConfig
	counter uint16
}