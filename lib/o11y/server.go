package o11y

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/pprof"

	"github.com/brunoluiz/go-lab/lib/closer"
	"github.com/brunoluiz/go-lab/lib/httpx"
	"github.com/hellofresh/health-go/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

func Run(ctx context.Context, logger *slog.Logger, healthz *health.Health, opts ...Option) error {
	o := &options{
		addr: "0.0.0.0:9090",
	}
	for _, opt := range opts {
		opt(o)
	}

	mux := http.NewServeMux()

	mux.Handle("/metrics", promhttp.Handler())
	mux.Handle("/healthz", healthz.Handler())
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	srv := httpx.New(o.addr, mux, httpx.WithName("o11y"), httpx.WithLogger(logger))
	defer closer.WithLogContext(ctx, logger, "failed to shutdown o11y server", srv.Close)

	return srv.Run(ctx)
}
