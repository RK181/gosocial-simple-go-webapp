package middlewares

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type gzipWrappedWriter struct {
	http.ResponseWriter
	io.Writer
}

func (gzw *gzipWrappedWriter) Write(b []byte) (int, error) {
	return gzw.Writer.Write(b)
}

func CompressGzip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}
		// Set the appropriate gzip header
		w.Header().Set("Content-Encoding", "gzip")
		// Compress the response
		gw, _ := gzip.NewWriterLevel(w, gzip.BestSpeed)
		defer gw.Close()

		gww := &gzipWrappedWriter{
			ResponseWriter: w,
			Writer:         gw,
		}

		next.ServeHTTP(gww, r)
		//log.Println("Compressed response using gzip")
	})
}
