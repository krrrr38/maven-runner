package utils

import (
	"github.com/motemen/go-colorine"
)

var logger = colorine.NewLogger(
	colorine.Prefixes{
		"verbose": colorine.Verbose,
		"notice":  colorine.Notice,
		"info":    colorine.Info,
		"warn":    colorine.Warn,
		"error":   colorine.Error,
		"debug": colorine.TextStyle{
			colorine.White,
			colorine.None,
		},
	},
	colorine.Info,
)

// Log outputs `message` with `prefix` by go-colorine
func Log(prefix, message string) {
	logger.Log(prefix, message)
}
