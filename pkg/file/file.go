package file

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

type PhotoFile struct {
	Path      string
	Ext       string
	CreatedAt time.Time
}

//go:generate mockery --name FileProvider
type FileProvider interface {
	GetAllFiles(path string) ([]PhotoFile, error)
	CreatedAt(path string) (time.Time, error)
}

func GetAllFiles(path string) ([]PhotoFile, error) {
	var photoFiles []PhotoFile

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
			createdAt, err := CreatedAt(path)
			if err != nil {
				return fmt.Errorf("failed to get the creation date of the file %s: %w", path, err)
			}

			file := PhotoFile{
				Path:      path,
				Ext:       strings.ToLower(filepath.Ext(path)),
				CreatedAt: createdAt,
			}

			photoFiles = append(photoFiles, file)
		}

		return nil
	})

	return photoFiles, err
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

	return time.Unix(stat.Birthtimespec.Sec, stat.Birthtimespec.Nsec), nil
}
