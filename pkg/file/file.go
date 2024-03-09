package file

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

//go:generate mockery --name FileProvider
type FileProvider interface {
	GetAllFiles(path string) ([]string, error)
	MoveFile(source string, destination string) error
	CreatedAt(path string) (time.Time, error)
}

func GetAllFiles(path string) ([]string, error) {
	var files []string

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err // Return the error to stop the walk
		}
		if info.IsDir() {
			return nil // Skip directories
		}

		// Check the file extension (case insensitive)
		switch strings.ToLower(filepath.Ext(path)) {
		case ".png", ".jpg", ".jpeg", ".raw":
			files = append(files, path)
		}

		return nil
	})

	return files, err
}

func CreatedAt(path string) (time.Time, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return time.Time{}, err
	}

	stat, ok := fileInfo.Sys().(*syscall.Stat_t)
	if !ok {
		return time.Time{}, fmt.Errorf("failed to get raw syscall.Stat_t data")
	}

	return time.Unix(stat.Ctimespec.Sec, stat.Ctimespec.Nsec), nil
}
