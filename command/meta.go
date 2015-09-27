package command

import "github.com/mitchellh/cli"

// ExitCodes
const (
	ExitCodeOK     int = 0
	ExitCodeFailed int = 1
)

// Meta contain the meta-option that nealy all subcommand inherited.
type Meta struct {
	UI cli.Ui
}
