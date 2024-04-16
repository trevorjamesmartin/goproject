// Package middleware ...
package middleware

import (
	"log"
	"net/http"
	"time"
)

// informedWriter ... http.ResponseWriter + statusCode
type informedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (i *informedWriter) WriteHeader(statusCode int) {
	i.ResponseWriter.WriteHeader(statusCode)
	i.statusCode = statusCode // inform
}

// Wrapper ...
type Wrapper struct {
	Condition   func(*http.Request) bool
	HandlerFunc http.HandlerFunc
	Alternate   map[string]http.Handler
}

func (w *Wrapper) writer(respWriter http.ResponseWriter) *informedWriter {
	iw := &informedWriter{
		ResponseWriter: respWriter,
		statusCode:     http.StatusOK,
	}
	return iw
}

// Handler ...
func (w *Wrapper) Handler() http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		i := w.writer(responseWriter)

		if w.Condition(request) {
			w.HandlerFunc.ServeHTTP(i.ResponseWriter, request)
		} else {
			// did we specified an alternate route ?
			if w.Alternate[request.Method] != nil {
				w.Alternate[request.Method].ServeHTTP(i.ResponseWriter, request)
			} else {
				http.NotFound(i.ResponseWriter, request)
			}
		}
	})
}

// Logging ... http.Request logger
func Logging(nextHandler http.Handler) http.Handler {
	mw := Wrapper{Condition: func(_ *http.Request) bool { return true }}
	mw.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		writer := mw.writer(w)
		nextHandler.ServeHTTP(writer, r)
		log.Println(writer.statusCode, r.Method, r.URL.Path, time.Since(start))
	}
	return mw.Handler()
}
