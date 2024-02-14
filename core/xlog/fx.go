package xlog

import "go.uber.org/fx"

var Module = fx.Module("log", fx.Provide(New))
