package middlewarechain

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestChainAlt(t *testing.T) {

	t.Run("No middleware", func(t *testing.T) {

		const (
			expectedStatusCode = http.StatusOK
			expectedMessage    = "Handler"
		)

		handlerFn := func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(expectedMessage))
		}
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)

		Chain(handlerFn).ServeHTTP(w, r)

		if sc := w.Result().StatusCode; sc != expectedStatusCode {
			t.Errorf("Got %v instead of %v", sc, expectedStatusCode)
		}

		resp, err := io.ReadAll(w.Result().Body)
		if err != nil {
			t.Fatal(err)

		} else if str := string(resp); str != expectedMessage {
			t.Errorf("Got %v instead %v", str, expectedMessage)
		}

	})

	t.Run("Multiple middlewares", func(t *testing.T) {
		type customTypeForContextPassing string

		const (
			message = "Handler"

			middlewareKey1 customTypeForContextPassing = "M1"
			middlewareKey2 customTypeForContextPassing = "M2"

			Value1 = string(middlewareKey1) + " >> "
			Value2 = string(middlewareKey2) + " >> "

			expectedMessage = Value1 + Value2 + message

			expectedStatusCode = http.StatusOK
		)

		middlewareFn1 := func(h http.HandlerFunc) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				ctx := context.WithValue(r.Context(), middlewareKey1, Value1)
				h.ServeHTTP(w, r.WithContext(ctx))
			}
		}

		middlewareFn2 := func(h http.HandlerFunc) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				ctx := context.WithValue(r.Context(), middlewareKey2, Value2)
				h.ServeHTTP(w, r.WithContext(ctx))
			}
		}

		handlerFn := func(w http.ResponseWriter, r *http.Request) {
			val1 := r.Context().Value(middlewareKey1).(string)
			val2 := r.Context().Value(middlewareKey2).(string)

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(val1 + val2 + message))
		}

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)

		Chain(handlerFn, middlewareFn1, middlewareFn2).ServeHTTP(w, r)

		if sc := w.Result().StatusCode; sc != expectedStatusCode {
			t.Errorf("Got %v instead of %v", sc, expectedStatusCode)
		}

		resp, err := io.ReadAll(w.Result().Body)
		if err != nil {
			t.Fatal(err)

		} else if str := string(resp); str != expectedMessage {
			t.Errorf("Got %v instead %v", str, expectedMessage)
		}

	})

}
