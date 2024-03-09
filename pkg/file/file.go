package file

import (
	"fmt"
	"io"
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
	CopyFile(src, dst string) error
	// CreatedAt(path string) (time.Time, error)
}

type File struct{}

func New() *File {
	return &File{}
}

func (f *File) GetAllFiles(path string) ([]PhotoFile, error) {
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
			createdAt, err := createdAt(path)
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

func (f *File) CopyFile(src, dst string) error {
	// Open the source file for reading
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// Create the destination file for writing
	destinationFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	// Copy the contents of the source file to the destination file
	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return err
	}

	// Ensure that any writes to the destination file are committed
	err = destinationFile.Sync()
	return err
}

func createdAt(path string) (time.Time, error) {
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
