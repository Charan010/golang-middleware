package ratelimiter

import (
	"fmt"
	"log"
	"math"
	"net"
	"net/http"
	"sync"
	"time"
)

type TokenBucket struct {
	Tokens         float64
	MaxTokens      float64
	RefillRate     float64
	LastRefillTime time.Time  // Timestamp of the last token refill
	mu             sync.Mutex //to ensure thread-safe access when multi-threading.
}

var Users = make(map[string]*TokenBucket) //hash map mapping ip's to TokenBucket struct
var usersMu sync.Mutex

// Creating bucket for every new user
func NewTokenBucket(maxTokens, refillRate float64) *TokenBucket {
	return &TokenBucket{
		Tokens:         maxTokens,
		MaxTokens:      maxTokens,
		RefillRate:     refillRate,
		LastRefillTime: time.Now(),
	}
}

// Refilling Tokens in bucket
func (tb *TokenBucket) refill() {

	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	duration := now.Sub(tb.LastRefillTime)
	tokensToAdd := tb.RefillRate * duration.Seconds()
	tb.Tokens = math.Min(math.Round(tb.Tokens+tokensToAdd), tb.MaxTokens)

	tb.LastRefillTime = now
}

// Checks whether tokens are present or not.
func (tb *TokenBucket) Request(tokens float64) bool {

	tb.refill()

	tb.mu.Lock()
	defer tb.mu.Unlock()

	if tokens <= tb.Tokens {
		tb.Tokens -= tokens
		return true
	}
	//if user exhausted all of the tokens, then return false.
	return false
}

// Wrapper for your function to implement rate limiting
func RateLimitMiddleware(maxTokens, refillRate float64, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		userIP, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.Error(w, "Unable to parse user IP", http.StatusInternalServerError)
			return
		}

		usersMu.Lock()
		tb, exists := Users[userIP]
		if !exists {
			tb = NewTokenBucket(maxTokens, refillRate)
			Users[userIP] = tb
		}
		usersMu.Unlock()

		if tb.Request(1) {
			remainingRequests := int(tb.Tokens)

			//Custom Header to send user how many request are still left.
			w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", remainingRequests))

			log.Printf("Request from IP: %s | Remaining Tokens: %.2f\n", userIP, tb.Tokens)

			next.ServeHTTP(w, r)

		} else {

			//reject the request from user when tokens are exhausted
			w.Header().Set("X-RateLimit-Remaining", "0")
			w.WriteHeader(http.StatusTooManyRequests)
			fmt.Fprintln(w, "Rate limit exceeded or you are blocked. Please try again later.")
		}
	})
}
