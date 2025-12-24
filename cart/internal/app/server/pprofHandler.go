package server

import (
	"net/http"
	_ "net/http/pprof"
)

func (s *Server) PprofHandler(w http.ResponseWriter, r *http.Request) {
	http.DefaultServeMux.ServeHTTP(w, r)
}
