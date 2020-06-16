package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Foo struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type wrapper struct {
	http.ResponseWriter
	mw io.Writer
}

func (w *wrapper) Write(b []byte) (int, error) {
	return w.mw.Write(b)
}

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf := bytes.NewBuffer(nil)
		wrap := &wrapper{w, io.MultiWriter(w, buf)}
		next.ServeHTTP(wrap, r)
		fmt.Println("[DEBUG]", buf.String())
	})
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", middleware(fooHandler()))
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}

func fooHandler() http.HandlerFunc {
	body := Foo{"Alice", 20}

	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(body)
	}
}
