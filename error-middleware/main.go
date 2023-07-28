package main

import (
	"context"
	"errors"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("/500", func(w http.ResponseWriter, r *http.Request) {
		err := errors.New("some error")
		WithError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	})

	http.ListenAndServe(":8080", ErrorHandler(&reporter{})(mux))
}

type reporter struct{}

func (r *reporter) Report(ctx context.Context, err error) {
	log.Println(err)
}
