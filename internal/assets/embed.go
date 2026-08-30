package assets

import "embed"

// Embedded contains the built-in Aurora theme and compiled administration UI.
//
//go:embed all:aurora all:admin
var Embedded embed.FS
