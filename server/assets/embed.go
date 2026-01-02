package assets

import "embed"

// Virtual filesystem
//
//go:embed *
var AssetsFS embed.FS
