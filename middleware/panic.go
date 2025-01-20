package middleware

import (
	"net/http"

	"github.com/gleb-korostelev/GophKeeper/tools/logger"
	"go.uber.org/zap"
)

// PanicMid is a middleware that recovers from panics in HTTP handlers,
// logs the error, and returns a 500 Internal Server Error response to the client.
//
// Workflow:
// 1. Wraps the next HTTP handler in a `defer` function to catch panics.
// 2. Logs the error, along with the URL and method of the request, if a panic occurs.
// 3. Sends an HTTP 500 Internal Server Error response to the client.
//
// Parameters:
// - next: The HTTP handler to be wrapped.
//
// Returns:
// - http.Handler: A wrapped HTTP handler with panic recovery.
//
// Example usage:
//
//	router := mux.NewRouter()
//	router.Use(middleware.PanicMid)
//	router.HandleFunc("/example", exampleHandler)
//
// Behavior:
// - If no panic occurs, the request proceeds normally.
// - If a panic occurs, the error is logged, and the client receives a 500 response.
func PanicMid(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			// Recover from any panic and handle the error.
			err := recover()
			if err != nil {
				// Log the error along with the request URL and method.
				logger.Error(err.(error).Error(), zap.String("url", r.RequestURI), zap.String("method", r.Method))

				// Return an HTTP 500 Internal Server Error response to the client.
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		// Call the next handler in the chain.
		next.ServeHTTP(w, r)
	})
}
