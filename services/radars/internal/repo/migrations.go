package repo

import "embed"

//go:embed migrations
var MigrationsFS embed.FS
