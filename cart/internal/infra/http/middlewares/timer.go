package middlewares

import (
	"log"
	"net/http"
	"time"
)

type TimerMux struct {
	h http.Handler
}

func NewTimeMux(h http.Handler) http.Handler {
	return &TimerMux{h: h}
}

func (m *TimerMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	m.h.ServeHTTP(w, r)
	log.Printf("handler spent %s", time.Since(now))
}
