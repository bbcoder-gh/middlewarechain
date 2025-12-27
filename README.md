# MiddlewareChain

A lightweight, zero-dependency Go library for elegantly chaining HTTP middleware functions with the standard library.

## Features
- 🚀 Zero dependencies - uses only Go's standard library
- 🔗 Chain multiple middleware functions with clean, readable syntax
- 📦 Tiny footprint - single file implementation
- ✅ Well-tested and production-ready
- 🎯 Works seamlessly with `net/http`

## Installation
```bash
go get github.com/bbcoder-gh/middlewarechain
```

## Quick Start
```go
package main

import (
    "log"
    "net/http"
    "time"
    
    "github.com/bbcoder-gh/middlewarechain"
)

func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        log.Printf("[%s] %s", r.Method, r.URL.Path)
        next(w, r)
        log.Printf("Request took %v", time.Since(start))
    }
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if token == "" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        next(w, r)
    }
}

func handleHello(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Hello, World!"))
}

func main() {
    // Chain middlewares: Logging -> Auth -> Handler
    http.HandleFunc("/api/hello", 
        middlewarechain.Chain(handleHello, LoggingMiddleware, AuthMiddleware))
    
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

## Usage

### Basic Chaining
```go
import "github.com/bbcoder-gh/middlewarechain"

// Apply middleware to a single handler
handler := middlewarechain.Chain(
    myHandler,
    LoggingMiddleware,
    AuthMiddleware,
    RateLimitMiddleware,
)

http.HandleFunc("/api/endpoint", handler)
```

### Middleware Execution Order

Middlewares execute in the order specified (left to right):
```go
// Execution flow:
// Request → RateLimitMiddleware → AuthMiddleware → LoggingMiddleware → Handler → Response

handler := middlewarechain.Chain(
    handleLogin,
    LoggingMiddleware,    // Executes third
    AuthMiddleware,       // Executes second  
    RateLimitMiddleware,  // Executes first
)
```

### Reusable Middleware Stacks
```go
// Create reusable middleware combinations
func protectedRoute(h http.HandlerFunc) http.HandlerFunc {
    return middlewarechain.Chain(h, LoggingMiddleware, AuthMiddleware)
}

func publicRoute(h http.HandlerFunc) http.HandlerFunc {
    return middlewarechain.Chain(h, LoggingMiddleware, CORSMiddleware)
}

// Use across multiple routes
http.HandleFunc("/api/users", protectedRoute(handleUsers))
http.HandleFunc("/api/profile", protectedRoute(handleProfile))
http.HandleFunc("/api/health", publicRoute(handleHealth))
```

### Writing Custom Middleware
```go
// Middleware signature
type Middleware func(http.HandlerFunc) http.HandlerFunc

// Example: Rate limiting middleware
func RateLimitMiddleware(next http.HandlerFunc) http.HandlerFunc {
    limiter := rate.NewLimiter(10, 20) // 10 req/sec, burst of 20
    
    return func(w http.ResponseWriter, r *http.Request) {
        if !limiter.Allow() {
            http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
            return
        }
        next(w, r)
    }
}

// Example: Context-based middleware
func TenantMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        tenantID := r.Header.Get("X-Tenant-ID")
        ctx := context.WithValue(r.Context(), "tenantID", tenantID)
        next(w, r.WithContext(ctx))
    }
}
```

## Common Middleware Examples

### CORS Middleware
```go
func CORSMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        next(w, r)
    }
}
```

### Request ID Middleware
```go
func RequestIDMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        requestID := uuid.New().String()
        ctx := context.WithValue(r.Context(), "requestID", requestID)
        w.Header().Set("X-Request-ID", requestID)
        next(w, r.WithContext(ctx))
    }
}
```

### Recovery Middleware
```go
func RecoveryMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                log.Printf("Panic recovered: %v", err)
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            }
        }()
        next(w, r)
    }
}
```

## API Reference

### Chain
```go
func Chain(h http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc
```

Chains multiple middleware functions together and returns the final handler.

**Parameters:**
- `h` - The final handler function to be executed
- `middlewares` - Variadic middleware functions to apply (executed in order)

**Returns:**
- The composed handler with all middlewares applied

**Example:**
```go
handler := middlewarechain.Chain(
    handleRequest,
    Middleware1,
    Middleware2,
    Middleware3,
)
```

## Testing

Run the test suite:
```bash
go test
```

Run tests with coverage:
```bash
go test -cover
```

Run tests with verbose output:
```bash
go test -v
```

## Real-World Example
```go
package main

import (
    "context"
    "encoding/json"
    "log"
    "net/http"
    "time"
    
    "github.com/bbcoder-gh/middlewarechain"
)

// Middleware stack for authenticated API routes
func apiMiddleware(h http.HandlerFunc) http.HandlerFunc {
    return middlewarechain.Chain(h,
        LoggingMiddleware,
        RecoveryMiddleware,
        CORSMiddleware,
        TenantMiddleware,
        AuthMiddleware,
        RateLimitMiddleware,
    )
}

// Public routes (no auth)
func publicMiddleware(h http.HandlerFunc) http.HandlerFunc {
    return middlewarechain.Chain(h,
        LoggingMiddleware,
        RecoveryMiddleware,
        CORSMiddleware,
    )
}

func main() {
    // Public routes
    http.HandleFunc("/api/health", publicMiddleware(handleHealth))
    http.HandleFunc("/api/auth/login", publicMiddleware(handleLogin))
    
    // Protected routes
    http.HandleFunc("/api/users", apiMiddleware(handleUsers))
    http.HandleFunc("/api/messages", apiMiddleware(handleMessages))
    http.HandleFunc("/api/settings", apiMiddleware(handleSettings))


    //Or just simplest
    http.HandleFunc("/api/hello", middlewarechain.Chain(handleHello, LoggingMiddleware))
    
    log.Println("Server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

## Why MiddlewareChain?

- **Simple**: Just one function, easy to understand
- **Standard**: Works with any `http.HandlerFunc`
- **Composable**: Build complex middleware stacks from simple pieces
- **Testable**: Each middleware can be tested independently
- **Performant**: Minimal overhead, compiles to efficient code

## Comparison with Alternatives

Unlike larger frameworks, MiddlewareChain:
- Has zero dependencies
- Works with the standard library
- Doesn't force a specific router or framework
- Is just 10 lines of code you can understand completely

## Contributing

See [CONTRIBUTING.md](./CONTRIBUTING.md) for guidelines on:
- Reporting bugs
- Suggesting features
- Submitting pull requests
- Code style and testing requirements

## License

See [LICENSE](./LICENSE) for details.

## Support

- 📖 [Documentation](https://github.com/bbcoder-gh/middlewarechain)
- 🐛 [Issue Tracker](https://github.com/bbcoder-gh/middlewarechain/issues)
- 💬 [Discussions](https://github.com/bbcoder-gh/middlewarechain/discussions)

---

**Maintained by [bbcoder-gh](https://github.com/bbcoder-gh)**

⭐ Star this repo if you find it useful!