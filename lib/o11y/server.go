package o11y

import (
	"log/slog"
	"net/http"
	"net/http/pprof"

	"github.com/brunoluiz/go-lab/lib/httpx"
)

type Option func(*options)

type options struct {
	addr string
}

func WithAddr(addr string) Option {
	return func(o *options) {
		o.addr = addr
	}
}

func New(logger *slog.Logger, opts ...Option) *httpx.Server {
	o := &options{
		addr: "0.0.0.0:9090",
	}
	for _, opt := range opts {
		opt(o)
	}

	mux := http.NewServeMux()

	// 1. Health Endpoints
	// mux.HandleFunc("/healthz", healthzHandler) // Liveness
	// mux.HandleFunc("/ready", readyHandler)     // Readiness

	// 2. Metrics Endpoint
	// mux.Handle("/metrics", promhttp.Handler())

	// 3. Profiling/Debugging Endpoints (pprof)
	// Registers standard pprof handlers under /debug/pprof/
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	return httpx.New(o.addr, mux, httpx.WithName("o11y"), httpx.WithLogger(logger))
}
