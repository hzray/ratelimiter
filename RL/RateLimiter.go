package RL

type RateLimiter interface {
	Refill()
	Consume() bool
}
