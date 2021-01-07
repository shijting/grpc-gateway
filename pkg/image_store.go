package pkg

import (
	"bytes"
	"fmt"
	"github.com/google/uuid"
	"os"
)

// ImageStore is an interface to store  images
type ImageStore interface {
	Save(imageType string, imageData bytes.Buffer) (string, error)
}

// DiskImageStore stores image on disk
type DiskImageStore struct {
	imageFolder string
}

// ImageInfo contains information of the  image
type ImageInfo struct {
	Type     string
	Path     string
}

// NewDiskImageStore returns a new DiskImageStore
func NewDiskImageStore(imageFolder string) *DiskImageStore {
	return &DiskImageStore{
		imageFolder: imageFolder,
	}
}

func (store *DiskImageStore) Save(
	imageType string,
	imageData bytes.Buffer,
) (string, error) {
	imageID, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("cannot generate image id: %w", err)
	}

	imagePath := fmt.Sprintf("%s/%s%s", store.imageFolder, imageID, imageType)

	file, err := os.Create(imagePath)
	if err != nil {
		return "", fmt.Errorf("cannot create image file: %w", err)
	}
	_, err = imageData.WriteTo(file)
	if err != nil {
		return "", fmt.Errorf("cannot write image to file: %w", err)
	}

	return imagePath, nil
}
