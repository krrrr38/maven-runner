package repository

import (
	"bytes"
	"fmt"
	"github.com/krrrr38/maven-runner/utils"
	"launchpad.net/xmlpath"
)

// See Rest API in http://search.maven.org/#api
const (
	// DefaultMavenScheme const
	DefaultMavenScheme = "http"
	// DefaultMavenHost const
	DefaultMavenHost = "repo1.maven.org"
)

// MavenRepository struct
type MavenRepository struct {
	scheme string
	host   string
}

// NewMavenRepository is constructor of MavenRepository struct.
func NewMavenRepository() *MavenRepository {
	return &MavenRepository{
		scheme: DefaultMavenScheme,
		host:   DefaultMavenHost,
	}
}

// LatestVersion find latest version of the artifact
func (r *MavenRepository) LatestVersion(artifact *Artifact) (version string, err error) {
	metadataURL := fmt.Sprintf("%s/maven-metadata.xml", r.artifactURLBase(artifact))
	utils.Log("info", fmt.Sprintf("Search latest version to %s: %+v", metadataURL, artifact))
	content, err := ValidateDownload(metadataURL)
	if err != nil {
		return "", err
	}
	xmlPath := xmlpath.MustCompile(`/metadata/versioning/latest`)
	root, err := xmlpath.Parse(bytes.NewReader(content))
	if err != nil {
		return "", err
	}
	if latestVersion, ok := xmlPath.String(root); ok {
		utils.Log("info", fmt.Sprintf("Found latest version: %s", latestVersion))
		return latestVersion, nil
	}
	return "", fmt.Errorf("Could not find latest version in %s", metadataURL)
}

// DownloadJarArtifact try to download jar artifact into outputDir
func (r *MavenRepository) DownloadJarArtifact(artifact *Artifact) ([]byte, error) {
	url := fmt.Sprintf("%s/%s/%s", r.artifactURLBase(artifact), artifact.Version, artifact.JarFileName())
	utils.Log("info", fmt.Sprintf("Try to download jar: %s", url))
	return ValidateDownload(url)
}

func (r *MavenRepository) artifactURLBase(artifact *Artifact) string {
	return fmt.Sprintf("%s://%s/maven2/%s/%s", r.scheme, r.host, artifact.GroupIDPath(), artifact.Name)
}
