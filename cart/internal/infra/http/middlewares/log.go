package middlewares

import (
	"log"
	"net/http"
)

type LogMux struct {
	h http.Handler
}

func NewLogMux(h http.Handler) http.Handler {
	return &LogMux{h: h}
}

func (m *LogMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("request got. Method: %s; URL; %s", r.Method, r.URL.Path)

	m.h.ServeHTTP(w, r)

	log.Printf("request processed. Method: %s; URL; %s", r.Method, r.URL.Path)
}
