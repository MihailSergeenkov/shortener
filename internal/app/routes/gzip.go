package routes

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/MihailSergeenkov/shortener/internal/app/logger"
	"go.uber.org/zap"
)

type compressWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}

func (c *compressWriter) Write(p []byte) (int, error) {
	return c.zw.Write(p) //nolint:wrapcheck // Нужно обернуть, но возврат должен остаться оригинальным
}

func (c *compressWriter) WriteHeader(statusCode int) {
	if statusCode < http.StatusMultipleChoices {
		c.w.Header().Set("Content-Encoding", "gzip")
	}
	c.w.WriteHeader(statusCode)
}

func (c *compressWriter) Close() error {
	return c.zw.Close() //nolint:wrapcheck // Нужно обернуть, но возврат должен остаться оригинальным
}

type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, fmt.Errorf("failed to init gzip reader: %w", err)
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

func (c compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p) //nolint:wrapcheck // Нужно обернуть, но возврат должен остаться оригинальным
}

func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return fmt.Errorf("failed to close base reader: %w", err)
	}
	return c.zr.Close() //nolint:wrapcheck // Нужно обернуть, но возврат должен остаться оригинальным
}

func gzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ow := w

		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		if supportsGzip {
			cw := newCompressWriter(w)
			ow = cw
			defer func(cw *compressWriter) {
				err := cw.Close()

				if err != nil {
					log.Printf("failed to close compress writer: %v", err)
				}
			}(cw)
		}

		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			contentType := r.Header.Get("Content-Type")
			if !(strings.Contains(contentType, "application/json") || strings.Contains(contentType, "text/html")) {
				logger.Log.Warn("content encoding for bad content type", zap.String("content_type", contentType))
			}

			cr, err := newCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Printf("failed to create compress reader: %v", err)
				return
			}

			r.Body = cr
			defer func(cr *compressReader) {
				err := cr.Close()

				if err != nil {
					log.Printf("failed to close compress reader: %v", err)
				}
			}(cr)
		}

		next.ServeHTTP(ow, r)
	})
}
