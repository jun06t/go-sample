package main

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	mux := chi.NewRouter()
	mux.Use(
		ErrorHandler(&reporter{}),
		//AuthHandler,
	)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("/500", func(w http.ResponseWriter, r *http.Request) {
		err := errors.New("some error")
		sendError(w, r, err)
	})

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Panic(err)
	}
}

type reporter struct{}

func (r *reporter) Report(ctx context.Context, err error) {
	log.Println(err)
}

func sendError(w http.ResponseWriter, r *http.Request, err error) {
	WithError(r, err)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("Hello, World!"))
}
