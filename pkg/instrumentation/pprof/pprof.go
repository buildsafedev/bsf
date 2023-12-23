package pprof

import (
	"net/http"
	"net/http/pprof"
)

// Handler returns the pprof handler
func Handler(mux *http.ServeMux) http.Handler {
	mux.HandleFunc("/pprof/", pprof.Index)
	mux.HandleFunc("/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/pprof/profile", pprof.Profile)
	mux.HandleFunc("/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/pprof/trace", pprof.Trace)
	return mux
}
