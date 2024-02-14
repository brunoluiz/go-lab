package xgin

import "go.uber.org/fx"

var Module = fx.Module("xgin", fx.Provide(New))
