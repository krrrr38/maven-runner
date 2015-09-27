package repository

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/krrrr38/maven-runner/utils"
	"io/ioutil"
	"net/http"
)

// ValidateDownload try to download file and validate checksum
func ValidateDownload(url string) ([]byte, error) {
	utils.Log("info", fmt.Sprintf("Validate checksum: %s", url))
	contents, err := download(url)
	if err != nil {
		return []byte{}, err
	}
	md5URL := url + ".md5"
	if md5Response, err := download(md5URL); err == nil {
		md5Checksum := md5.Sum(contents)
		if hex.EncodeToString(md5Checksum[:]) == string(md5Response) {
			return contents, nil
		}
		utils.Log("debug", fmt.Sprintf("Failed to validate md5 checksum: actual=%s, expected=%s", hex.EncodeToString(md5Checksum[:]), string(md5Response)))
	}
	sha1URL := url + ".sha1"
	if sha1Response, err := download(sha1URL); err == nil {
		sha1Checksum := sha1.Sum(contents)
		if hex.EncodeToString(sha1Checksum[:]) == string(sha1Response) {
			return contents, nil
		}
		utils.Log("debug", fmt.Sprintf("Failed to validate sha1 checksum: actual=%s, expected=%s", hex.EncodeToString(sha1Checksum[:]), string(sha1Response)))
	}
	return []byte{}, fmt.Errorf("Checksum validation failed for %s", url)
}

func download(url string) ([]byte, error) {
	utils.Log("debug", fmt.Sprintf("Try to download file: %s", url))
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
