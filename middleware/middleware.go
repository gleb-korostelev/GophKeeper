package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type ctxKey uint8

const (
	ctxKeyAddress ctxKey = iota
	ctxKeyRole
)

type CoreMW struct {
}

func NewCoreMW() *CoreMW {
	return &CoreMW{}
}

func (a *CoreMW) Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a.contextUpdate(next).ServeHTTP(w, r)
	}
}

func (a *CoreMW) contextUpdate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GzipCompressMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(&gzipResponseWriter{ResponseWriter: w, Writer: w}, r)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func GzipDecompressMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Encoding") == "gzip" {
			gzReader, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, "Failed to create gzip reader", http.StatusBadRequest)
				return
			}
			defer gzReader.Close()
			r.Body = gzReader
			r.Header.Del("Content-Encoding")
		}
		next.ServeHTTP(w, r)
	})
}

type responseWriter struct {
	http.ResponseWriter
	status int
	size   int
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *responseWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	size, err := w.ResponseWriter.Write(b)
	w.size += size
	return size, err
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
	hasWrittenHeader bool
}

func (w *gzipResponseWriter) WriteHeader(statusCode int) {
	if w.hasWrittenHeader {
		return
	}
	w.hasWrittenHeader = true
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	if !w.hasWrittenHeader {
		contentType := w.Header().Get("Content-Type")
		if strings.HasPrefix(contentType, "application/json") || strings.HasPrefix(contentType, "text/html") {
			w.Header().Set("Content-Encoding", "gzip")
			w.Writer = gzip.NewWriter(w.ResponseWriter)
			defer func() {
				if w.Writer != nil {
					w.Writer.(*gzip.Writer).Close()
				}
			}()
		}
		w.WriteHeader(http.StatusOK)
	}
	return w.Writer.Write(b)
}
