package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"ratelimiter/RL"
	"ratelimiter/store"
)

var rl RL.RateLimiter

func main() {
	algo := flag.String("algo", "", "rate limiter algorithm to use")
	flag.Parse()
	store.InitRedisClient()
	switch *algo {
	case "TokenBucket":
		rl = RL.CreateTokenBucket(5, 3, 1)
	default:
		fmt.Println("Not supported algo")
		os.Exit(1)

	}
	http.HandleFunc("/test", rateLimiterMiddleware(testHandler))
	http.ListenAndServe(":8080", nil)
}

// rateLimiterMiddleware checks if a request can proceed based on the RateLimiter.
func rateLimiterMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if rl.Consume() {
			next(w, r)
		} else {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
		}
	}
}

// testHandler is a basic handler function to respond when rate limiting allows.
func testHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello! You're not rate limited.")
}
