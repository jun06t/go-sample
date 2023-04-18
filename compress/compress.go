package main

import (
	"compress/flate"
	"compress/gzip"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/andybalholm/brotli"
)

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cw, cleanup := factory(w, r)
		defer cleanup()
		defer cw.Close()

		next.ServeHTTP(cw, r)
	})
}

var brotliPool = sync.Pool{
	New: func() interface{} {
		//return brotli.NewWriterOptions(nil, brotli.WriterOptions{Quality: 11, LGWin: 22})
		return brotli.NewWriter(nil)
	},
}

var gzipPool = sync.Pool{
	New: func() interface{} {
		return gzip.NewWriter(nil)
	},
}

const flateLevel = 5

var flatePool = sync.Pool{
	New: func() interface{} {
		w, _ := flate.NewWriter(nil, flateLevel)
		return w
	},
}

func factory(w http.ResponseWriter, r *http.Request) (CustomResponseWriter, func()) {
	if strings.Contains(r.Header.Get("Accept-Encoding"), "br") {
		br := brotliPool.Get().(*brotli.Writer)
		br.Reset(w)
		cleanup := func() {
			brotliPool.Put(br)
		}
		w.Header().Set("Content-Encoding", "br")
		return &brotliResponseWriter{ResponseWriter: w, br: br}, cleanup
	}
	if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		/*
			gzWriter := gzip.NewWriter(w)
			cleanup := func() {}
		*/

		gzWriter := gzipPool.Get().(*gzip.Writer)
		gzWriter.Reset(w)
		cleanup := func() {
			gzipPool.Put(gzWriter)
		}
		w.Header().Set("Content-Encoding", "gzip")
		return &gzipResponseWriter{ResponseWriter: w, gz: gzWriter}, cleanup
	}
	if strings.Contains(r.Header.Get("Accept-Encoding"), "deflate") {
		flWriter := flatePool.Get().(*flate.Writer)
		flWriter.Reset(w)
		cleanup := func() {
			flatePool.Put(flWriter)
		}
		w.Header().Set("Content-Encoding", "deflate")
		return &deflateResponseWriter{ResponseWriter: w, fl: flWriter}, cleanup
	}
	return &noop{}, func() {}
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
