package main

import (
	"context"
	"net/http"
)

var userIDKey struct{}

func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func AuthHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := WithUserID(r.Context(), "user001")
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
