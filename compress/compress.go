package main

import (
	"compress/flate"
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/andybalholm/brotli"
)

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cw := factory(w, r)
		defer cw.Close()

		next.ServeHTTP(cw, r)
	})
}

func factory(w http.ResponseWriter, r *http.Request) CustomResponseWriter {
	if strings.Contains(r.Header.Get("Accept-Encoding"), "br") {
		br := brotli.NewWriter(w)
		w.Header().Set("Content-Encoding", "br")
		return &brotliResponseWriter{ResponseWriter: w, br: br}
	}
	if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		gzWriter := gzip.NewWriter(w)
		w.Header().Set("Content-Encoding", "gzip")
		return &gzipResponseWriter{ResponseWriter: w, gz: gzWriter}
	}
	if strings.Contains(r.Header.Get("Accept-Encoding"), "deflate") {
		flWriter, _ := flate.NewWriter(w, 5)
		w.Header().Set("Content-Encoding", "deflate")
		return &deflateResponseWriter{ResponseWriter: w, fl: flWriter}
	}
	return &noop{}
}

type CustomResponseWriter interface {
	http.ResponseWriter
	io.Closer
}

type noop struct {
	http.ResponseWriter
}

func (n *noop) Close() error {
	return nil
}

type deflateResponseWriter struct {
	http.ResponseWriter
	fl *flate.Writer
}

func (rw *deflateResponseWriter) Write(b []byte) (int, error) {
	return rw.fl.Write(b)
}

func (rw *deflateResponseWriter) Close() error {
	return rw.fl.Close()
}

type brotliResponseWriter struct {
	http.ResponseWriter
	br *brotli.Writer
}

func (rw *brotliResponseWriter) Write(b []byte) (int, error) {
	return rw.br.Write(b)
}

func (rw *brotliResponseWriter) Close() error {
	return rw.br.Close()
}

type gzipResponseWriter struct {
	http.ResponseWriter
	gz *gzip.Writer
}

func (g *gzipResponseWriter) Write(b []byte) (int, error) {
	return g.gz.Write(b)
}

func (g *gzipResponseWriter) Close() error {
	return g.gz.Close()
}
