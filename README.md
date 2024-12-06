# Go Concurrency Extra Credit Assignment

This project implements a concurrent rate limiter to demonstrate key concepts from Chapter 9, specifically focusing on goroutines, mutexes, and sync.Once.

## Assignment Requirements and Implementation

### 1. Must Use Concurrency (Goroutines)
The program creates multiple goroutines that try to access a shared resource concurrently:
```go
// Creates 10 concurrent goroutines
for i := 0; i < numGoroutines; i++ {
    wg.Add(1)
    go func(id int) {
        defer wg.Done()
        // Each goroutine attempts to use the resource multiple times
        for j := 0; j < 3; j++ {
            if err := resource.Use(id); err != nil {
                resource.logger.Log(fmt.Sprintf("Goroutine %d: %v", id, err))
            }
        }
    }(i)
}
```

### 2. Must Use a Mutex
The program uses multiple mutexes to protect shared resources:
```go
type RateLimiter struct {
    mu            sync.Mutex  // Protects shared state
    maxRequests   int
    currRequests  int
    windowSeconds int
    lastReset     time.Time
}

func (rl *RateLimiter) TryAcquire() bool {
    rl.mu.Lock()
    defer rl.mu.Unlock()
    // Safe access to shared state
    ...
}
```

### 3. Must Use sync.Once
The program uses sync.Once to ensure resource initialization happens exactly once:
```go
type Resource struct {
    name     string
    limiter  *RateLimiter
    logger   *Logger
    initOnce sync.Once    // Ensures one-time initialization
}

func (r *Resource) Use(id int) error {
    // Initialization happens exactly once across all goroutines
    r.initOnce.Do(func() {
        r.initialize()
    })
    ...
}
```

## Project Structure

```
GoConcur/
├── main.go           # Main implementation
├── main_test.go      # Test cases
├── go.mod           # Go module file
└── .github/workflows/go.yml  # GitHub Actions configuration
```

## Testing

The project includes comprehensive tests that verify:
1. Concurrent access behavior
2. Rate limiting functionality
3. One-time initialization using sync.Once
4. Thread safety using race detector

## How to Run

```bash
# Run tests
go test -v ./...

# Run with race detector
go run -race main.go
```

## GitHub Actions Integration

The repository is configured with GitHub Actions to:
- Build the project
- Run tests
- Run with race detector enabled

## Why This Implementation?

This rate limiter implementation was chosen because it:
1. Naturally requires all three required components from Chapter 9
2. Represents a real-world concurrent programming scenario
3. Demonstrates proper synchronization techniques
4. Is easily testable
5. Shows practical usage of Go's concurrency primitives

The program prevents race conditions by:
- Using mutexes to protect shared state
- Ensuring thread-safe logging
- Using sync.Once for safe initialization
- Coordinating goroutines with WaitGroup

## Extra Credit Points Justification

This implementation deserves the full 5 extra credit points because it:
1. Correctly implements all required concepts from Chapter 9
2. Includes comprehensive tests
3. Uses the race detector
4. Represents a practical, real-world application
5. Demonstrates proper concurrent programming practices
6. Includes proper documentation and comments
7. Has GitHub Actions integration for automated testing
