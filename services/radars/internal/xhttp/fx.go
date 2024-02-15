package xhttp

import (
	"go.uber.org/fx"
)

var Module = fx.Module("xhttp", fx.Invoke(RegisterRoutes))
