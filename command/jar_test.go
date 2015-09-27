package command

import (
	"testing"

	"github.com/mitchellh/cli"
)

func TestJarCommand_implement(t *testing.T) {
	var _ cli.Command = &JarCommand{}
}
