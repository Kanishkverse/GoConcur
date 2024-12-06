// main_test.go
package main

import (
	"sync"
	"testing"
	"time"
)

func TestRateLimiter(t *testing.T) {
	limiter := NewRateLimiter(3, 1) // 3 requests per second

	// Test basic acquisition
	if !limiter.TryAcquire() {
		t.Error("First acquisition should succeed")
	}
	if !limiter.TryAcquire() {
		t.Error("Second acquisition should succeed")
	}
	if !limiter.TryAcquire() {
		t.Error("Third acquisition should succeed")
	}
	if limiter.TryAcquire() {
		t.Error("Fourth acquisition should fail")
	}

	// Test release
	limiter.Release()
	if !limiter.TryAcquire() {
		t.Error("Acquisition after release should succeed")
	}
}

func TestRateLimiterConcurrent(t *testing.T) {
	limiter := NewRateLimiter(5, 1)
	var wg sync.WaitGroup
	successCount := 0
	var mu sync.Mutex

	// Launch 10 goroutines trying to acquire simultaneously
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if limiter.TryAcquire() {
				mu.Lock()
				successCount++
				mu.Unlock()
				time.Sleep(100 * time.Millisecond)
				limiter.Release()
			}
		}()
	}

	wg.Wait()

	if successCount != 5 {
		t.Errorf("Expected 5 successful acquisitions, got %d", successCount)
	}
}

func TestResourceInitialization(t *testing.T) {
	resource := NewResource("TestResource", 3, 1)
	var wg sync.WaitGroup
	initCount := 0
	var mu sync.Mutex

	// Launch multiple goroutines to test sync.Once behavior
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			// Override initialize for testing
			resource.initOnce.Do(func() {
				mu.Lock()
				initCount++
				mu.Unlock()
			})
			_ = resource.Use(id)
		}(i)
	}

	wg.Wait()

	if initCount != 1 {
		t.Errorf("Expected initialization to happen exactly once, got %d times", initCount)
	}
}

func TestResourceRateLimiting(t *testing.T) {
	resource := NewResource("TestResource", 2, 1) // 2 requests per second
	var wg sync.WaitGroup
	errorCount := 0
	var mu sync.Mutex

	// Launch 5 concurrent requests
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			if err := resource.Use(id); err != nil {
				mu.Lock()
				errorCount++
				mu.Unlock()
			}
		}(i)
	}

	wg.Wait()

	if errorCount != 3 { // 5 requests - 2 allowed = 3 errors
		t.Errorf("Expected 3 rate limit errors, got %d", errorCount)
	}
}
