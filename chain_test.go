package middlewarechain

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Test helpers
func executeRequest(handler http.Handler) string {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	handler.ServeHTTP(w, r)

	body, _ := io.ReadAll(w.Result().Body)
	return string(body)
}

func prefixMiddleware(prefix string) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(prefix))
			next.ServeHTTP(w, r)
		}
	}
}

func TestChain(t *testing.T) {
	// Handler writes "Handler" to trace execution
	handler := func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("Handler"))
	}

	tests := []struct {
		name       string
		middleware []Middleware
		want       string
	}{
		{
			name:       "no middleware",
			middleware: nil,
			want:       "Handler",
		},
		{
			name: "one middleware",
			middleware: []Middleware{
				prefixMiddleware("Middleware --> "),
			},
			want: "Middleware --> Handler",
		},
		{
			name: "two middlewares",
			middleware: []Middleware{
				prefixMiddleware("First Middleware -->"),
				prefixMiddleware("Second Middleware -->"),
			},
			want: "First Middleware -->Second Middleware -->Handler",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := executeRequest(Chain(handler, tt.middleware...))

			if got != tt.want {
				t.Errorf("Chain() = %q, want %q", got, tt.want)
			}
		})
	}
}
