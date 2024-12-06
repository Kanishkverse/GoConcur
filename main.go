// main.go
package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// RateLimiter manages resource access with configurable limits
type RateLimiter struct {
	mu            sync.Mutex
	maxRequests   int
	currRequests  int
	windowSeconds int
	lastReset     time.Time
}

// Resource represents a shared resource that needs rate limiting
type Resource struct {
	name     string
	limiter  *RateLimiter
	logger   *Logger
	initOnce sync.Once
}

// Logger provides thread-safe logging
type Logger struct {
	mu sync.Mutex
}

func (l *Logger) Log(message string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	log.Printf("%s: %s\n", time.Now().Format("15:04:05"), message)
}

// NewRateLimiter creates a new rate limiter with specified limits
func NewRateLimiter(maxRequests, windowSeconds int) *RateLimiter {
	return &RateLimiter{
		maxRequests:   maxRequests,
		windowSeconds: windowSeconds,
		lastReset:     time.Now(),
	}
}

// TryAcquire attempts to acquire a rate limit token
func (rl *RateLimiter) TryAcquire() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	if now.Sub(rl.lastReset) >= time.Duration(rl.windowSeconds)*time.Second {
		rl.currRequests = 0
		rl.lastReset = now
	}

	if rl.currRequests >= rl.maxRequests {
		return false
	}

	rl.currRequests++
	return true
}

// Release releases a rate limit token
func (rl *RateLimiter) Release() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.currRequests > 0 {
		rl.currRequests--
	}
}

// NewResource creates a new resource with rate limiting
func NewResource(name string, maxRequests, windowSeconds int) *Resource {
	return &Resource{
		name:    name,
		limiter: NewRateLimiter(maxRequests, windowSeconds),
		logger:  &Logger{},
	}
}

// initialize performs one-time initialization of the resource
func (r *Resource) initialize() {
	r.logger.Log(fmt.Sprintf("Initializing resource: %s", r.name))
	// Simulate some initialization work
	time.Sleep(100 * time.Millisecond)
}

// Use attempts to use the resource with rate limiting
func (r *Resource) Use(id int) error {
	// Ensure initialization happens exactly once
	r.initOnce.Do(func() {
		r.initialize()
	})

	if !r.limiter.TryAcquire() {
		return fmt.Errorf("rate limit exceeded for resource %s", r.name)
	}
	defer r.limiter.Release()

	r.logger.Log(fmt.Sprintf("Goroutine %d using resource: %s", id, r.name))
	// Simulate some work
	time.Sleep(200 * time.Millisecond)
	return nil
}

func main() {
	// Create a shared resource with rate limiting
	resource := NewResource("DatabaseConnection", 3, 1) // max 3 requests per second

	// Create multiple goroutines trying to access the resource
	var wg sync.WaitGroup
	numGoroutines := 10

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// Each goroutine tries to use the resource multiple times
			for j := 0; j < 3; j++ {
				if err := resource.Use(id); err != nil {
					resource.logger.Log(fmt.Sprintf("Goroutine %d: %v", id, err))
				}
				// Random delay between attempts
				time.Sleep(time.Duration(100+id*50) * time.Millisecond)
			}
		}(i)
	}

	wg.Wait()
	resource.logger.Log("All goroutines completed")
}
