package repository

import (
	"fmt"
	"github.com/krrrr38/maven-runner/utils"
)

// JarResolver struct
type JarResolver struct {
	JarDir    JarDirectory
	LocalRepo LocalRepository
	MavenRepo MavenRepository
}

// ResolveJar get jar from 1. jar directory, 2. local repository, 3. maven repository
// then set jar into jar directory and return jar path
func (r *JarResolver) ResolveJar(artifact *Artifact) (string, error) {
	utils.Log("info", fmt.Sprintf("Resolve jar: %+v", artifact))
	// jar directory
	jarPath, jarDirErr := r.JarDir.FindJarArtifact(artifact)
	if jarDirErr == nil {
		return jarPath, nil
	}
	utils.Log("debug", fmt.Sprintf("Could not find jar in jar directory: %+v", r.JarDir))

	// local repository
	path, localRepoErr := r.LocalRepo.FindJarArtifact(artifact)
	if localRepoErr == nil {
		// copy from local repo to jar directory
		if jarPath, err := r.JarDir.CopyFile(path, artifact.JarFileName()); err == nil {
			return jarPath, nil
		}
	}
	utils.Log("debug", fmt.Sprintf("Could not find jar in local repository: %+v", r.LocalRepo))

	// maven repository
	content, mavenRepoErr := r.MavenRepo.DownloadJarArtifact(artifact)
	if mavenRepoErr == nil {
		return r.JarDir.WriteFile(content, artifact.JarFileName())
	}
	utils.Log("debug", fmt.Sprintf("Could not find jar in maven repository: %+v", r.MavenRepo))

	return "", fmt.Errorf("Failed to resolve jar : [%v,%v,%v]", jarDirErr, localRepoErr, mavenRepoErr)
}
