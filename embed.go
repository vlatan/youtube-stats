package resources

import "embed"

//go:embed web/static web/templates
var Files embed.FS
