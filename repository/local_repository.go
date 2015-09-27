package repository

import (
	"fmt"
	"github.com/krrrr38/maven-runner/utils"
	"github.com/mitchellh/go-homedir"
	"os"
)

const (
	// DefaultLocalRepositoryDirectory is defaut local repository directory path.
	DefaultLocalRepositoryDirectory = "%s/.m2/repository"
)

// LocalRepository struct
type LocalRepository struct {
	// repository directory
	Dir string
}

// NewLocalRepository is constructor of LocalRepository struct.
func NewLocalRepository() (*LocalRepository, error) {
	homeDir, err := homedir.Dir()
	if err != nil {
		return nil, err
	}

	localRepositoryDir := fmt.Sprintf(DefaultLocalRepositoryDirectory, homeDir)
	return &LocalRepository{
		Dir: localRepositoryDir,
	}, nil
}

// FindJarArtifact get artifact file path in local storage
func (r *LocalRepository) FindJarArtifact(artifact *Artifact) (string, error) {
	path := fmt.Sprintf("%s/%s/%s/%s/%s", r.Dir, artifact.GroupIDPath(), artifact.Name, artifact.Version, artifact.JarFileName())
	utils.Log("info", fmt.Sprintf("Try to find jar in local repository: %s", path))
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return "", fmt.Errorf("Could not find local jar")
	}
	utils.Log("info", fmt.Sprintf("Found jar in local repository: %s", path))
	return path, nil
}
