package filesUC

import (
	"path/filepath"

	"passkeeper/internal/entities"
	"passkeeper/internal/entities/compress"
	"passkeeper/internal/entities/hashes"
)

// FilesUsecaseInf describe some actions under file blobs
type FilesUsecaseInf interface {
	ZipAndUpload(uploadingFile string) (file *entities.File, err error)
	UnzipAndDownload(dirToDownload string, file *entities.File) error
}

// FilesUC implement FilesUsecaseInf
type FilesUC struct {
}

// ZipAndUpload zip local file and return file blob
func (f *FilesUC) ZipAndUpload(uploadingFile string) (file *entities.File, err error) {
	zipData, err := compress.CompressFile(uploadingFile)
	if err != nil {
		return nil, err
	}

	name := filepath.Base(uploadingFile)

	file = &entities.File{
		ID:   hashes.GeneratePassID2(),
		Type: entities.BlobFile,
		Name: name,
		Body: zipData,
	}

	return file, nil
}

// UnzipAndDownload unzip file blob and save it
func (f *FilesUC) UnzipAndDownload(dirToDownload string, file *entities.File) error {
	fileName := filepath.Join(dirToDownload, file.Name)

	err := compress.DecompressFile(file.Body, fileName)
	if err != nil {

		return err
	}

	return nil
}
