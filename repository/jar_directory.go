package repository

import (
	"fmt"
	"github.com/krrrr38/maven-runner/utils"
	"io/ioutil"
	"os"
)

// JarDirectory for maven-runner
type JarDirectory struct {
	// repository directory
	Dir string
}

// FindJarArtifact get artifact file path in local storage
func (d *JarDirectory) FindJarArtifact(artifact *Artifact) (string, error) {
	utils.Log("info", fmt.Sprintf("Try to find jar in jar directory: %+v", artifact))
	path := fmt.Sprintf("%s/%s", d.Dir, artifact.JarFileName())
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return "", fmt.Errorf("Could not find jar in jar directory")
	}
	utils.Log("info", fmt.Sprintf("Found jar in jar directory: %s", path))
	return path, nil
}

// WriteFile write data in jar directory file
func (d *JarDirectory) WriteFile(data []byte, filename string) (string, error) {
	path := fmt.Sprintf("%s/%s", d.Dir, filename)
	utils.Log("info", fmt.Sprintf("Write file in jar directory: %s", path))
	if err := ioutil.WriteFile(path, data, 0644); err != nil {
		return "", err
	}
	return path, nil
}

// CopyFile copy path file to jar directory file
func (d *JarDirectory) CopyFile(path, filename string) (string, error) {
	utils.Log("info", fmt.Sprintf("Copy file in jar directory: from %s", path))
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return d.WriteFile(data, filename)
}
