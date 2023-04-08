package main

import (
	"io"
	"net/http"
	"os"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Hello, World!\n"))
	})

	mux.HandleFunc("/stream", func(w http.ResponseWriter, r *http.Request) {
		file, err := os.Open("sample.txt")
		if err != nil {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
		defer file.Close()

		w.Header().Set("Content-Type", "text/plain")

		io.Copy(w, file)
	})

	http.ListenAndServe(":8080", middleware(mux))
}
