package job

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"os"
)

func Sha1OfFile(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}

	hasher := sha1.New()
	_, err = io.Copy(hasher, f)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}
