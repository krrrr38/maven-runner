package main

import (
	"github.com/krrrr38/maven-runner/command"
	"github.com/mitchellh/cli"
)

// Commands are collections of maven-runner commands.
func Commands(meta *command.Meta) map[string]cli.CommandFactory {
	return map[string]cli.CommandFactory{
		"jar": func() (cli.Command, error) {
			return &command.JarCommand{
				Meta: *meta,
			}, nil
		},

		"version": func() (cli.Command, error) {
			return &command.VersionCommand{
				Meta:     *meta,
				Version:  Version,
				Revision: GitCommit,
				Name:     Name,
			}, nil
		},
	}
}
