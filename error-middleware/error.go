package main

import (
	"context"
	"net/http"
)

type Reporter interface {
	Report(ctx context.Context, err error)
}

// ErrorHandler reports errors to the reporter.
func ErrorHandler(reporter Reporter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sw := &statusResponseWriter{ResponseWriter: w}
			next.ServeHTTP(sw, r)

			if sw.status >= http.StatusInternalServerError {
				err := errorFromContext(r.Context())
				if err == nil {
					return
				}
				// notify
				reporter.Report(r.Context(), err)
			}
		})
	}
}

type statusResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusResponseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

var errorContextKey struct{}

// WithError attaches an error to the request's context.
// It modifies the request in place to ensure that later middleware and handlers
// using the original pointer to the request can access the updated context.
func WithError(r *http.Request, err error) {
	if err == nil {
		return
	}
	r2 := r.WithContext(withError(r.Context(), err))
	*r = *r2
}

// withError sets an error to the context.
func withError(ctx context.Context, err error) context.Context {
	if err == nil {
		return ctx
	}
	return context.WithValue(ctx, errorContextKey, err)
}

// errorFromContext returns an error from the context.
func errorFromContext(ctx context.Context) error {
	if err, ok := ctx.Value(errorContextKey).(error); ok {
		return err
	}
	return nil
}
