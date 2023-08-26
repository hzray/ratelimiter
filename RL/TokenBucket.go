package RL

import (
	"fmt"
	"github.com/go-redis/redis"
	"math"
	"ratelimiter/store"
	"strconv"
	"time"
)

type TokenBucket struct {
	tokens     float64
	capacity   float64
	refillRate float64
	userid     int
}

func CreateTokenBucket(capacity, refillRate float64, userid int) *TokenBucket {
	tb := &TokenBucket{
		tokens:     capacity,
		capacity:   capacity,
		refillRate: refillRate,
		userid:     userid,
	}
	fmt.Println("Token Bucket rate limiter created")
	return tb
}

func (tb *TokenBucket) Refill() {
	now := time.Now()
	value, err := store.RDB.Get(strconv.Itoa(tb.userid) + "_last_refill_time").Result()
	if err == redis.Nil {
		store.RDB.Set(strconv.Itoa(tb.userid)+"_tokens", tb.capacity, 0)
		store.RDB.Set(strconv.Itoa(tb.userid)+"_last_refill_time", now.Unix(), 0)
		return
	}

	timeInUnix, _ := strconv.ParseInt(value, 10, 64)
	lastRefillTime := time.Unix(timeInUnix, 0)
	duration := now.Sub(lastRefillTime)

	value, _ = store.RDB.Get(strconv.Itoa(tb.userid) + "_tokens").Result()
	tokens, _ := strconv.ParseFloat(value, 64)
	tokensToAdd := tb.refillRate * duration.Minutes()
	store.RDB.Set(strconv.Itoa(tb.userid)+"_tokens", math.Min(tb.capacity, tokens+tokensToAdd), 0)
	store.RDB.Set(strconv.Itoa(tb.userid)+"_last_refill_time", now.Unix(), 0)
}

func (tb *TokenBucket) Consume() bool {
	tb.Refill()
	value, _ := store.RDB.Get(strconv.Itoa(tb.userid) + "_tokens").Result()
	tokens, _ := strconv.ParseFloat(value, 64)
	if tokens < 1 {
		return false
	}
	store.RDB.Set(strconv.Itoa(tb.userid)+"_tokens", tokens-1, 0)
	return true
}
