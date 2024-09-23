package filesUC

import (
	"path/filepath"

	"passkeeper/internal/entities/compress"
	"passkeeper/internal/entities/hashes"
	"passkeeper/internal/entities/structs"
)

// FilesUsecaseInf describe some actions under file blobs
type FilesUsecaseInf interface {
	ZipAndUpload(uploadingFile string) (file *structs.File, err error)
	UnzipAndDownload(dirToDownload string, file *structs.File) error
}

// FilesUC implement FilesUsecaseInf
type FilesUC struct {
}

// ZipAndUpload zip local file and return file blob
func (f *FilesUC) ZipAndUpload(uploadingFile string) (file *structs.File, err error) {
	zipData, err := compress.CompressFile(uploadingFile)
	if err != nil {
		return nil, err
	}

	name := filepath.Base(uploadingFile)

	file = &structs.File{
		ID:   hashes.GeneratePassID(),
		Type: structs.BlobFile,
		Name: name,
		Body: zipData,
	}

	return file, nil
}

// UnzipAndDownload unzip file blob and save it
func (f *FilesUC) UnzipAndDownload(dirToDownload string, file *structs.File) error {
	fileName := filepath.Join(dirToDownload, file.Name)

	err := compress.DecompressFile(file.Body, fileName)
	if err != nil {

		return err
	}

	return nil
}
