package repository

import (
	"fmt"
	"strings"
)

// Artifact struct
type Artifact struct {
	// groupId
	GroupID string

	// artifact name
	Name string

	// version
	Version string
}

// GroupIDPath return "/" separated group ID string
func (a *Artifact) GroupIDPath() string {
	return strings.Replace(a.GroupID, ".", "/", -1)
}

// JarFileName jar filename
func (a *Artifact) JarFileName() string {
	return fmt.Sprintf("%s-%s.jar", a.Name, a.Version)
}
