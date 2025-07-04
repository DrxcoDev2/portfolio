package utils

import (
	"os"
	"path/filepath"
	"time"
)

func GetLatestModTime(dirs ...string) time.Time {
	var latest time.Time

	for _, dir := range dirs {
		filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return nil
			}
			if modTime := info.ModTime(); modTime.After(latest) {
				latest = modTime
			}
			return nil
		})
	}

	return latest
}
