package compress

import (
	"bytes"
	"compress/gzip"
	"io"
	"os"
)

// CompressFile compress file and return compressed data
func CompressFile(filePath string) (data []byte, err error) {

	// open file to compress
	file, err := os.Open(filePath)
	if err != nil {

		return nil, err
	}
	defer file.Close()

	// init buffer and gzip writer
	var rw bytes.Buffer
	gw, err := gzip.NewWriterLevel(&rw, gzip.BestCompression)
	if err != nil {
		return nil, err
	}

	// compress file
	_, err = io.Copy(gw, file)
	if err != nil {

		return nil, err
	}

	// finish compression
	err = gw.Close()
	if err != nil {

		return nil, err
	}

	// successful compression
	return rw.Bytes(), nil
}

// DecompressFile uncompress data and save it to new file
func DecompressFile(data []byte, filePath string) (err error) {

	// open file for write, create if not exist
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {

		return err
	}
	defer file.Close()

	// init buffer and gzip reader
	rw := bytes.NewBuffer(data)
	gr, err := gzip.NewReader(rw)
	if err != nil {
		return err
	}

	// uncompress data and save to file
	_, err = io.Copy(file, gr)
	if err != nil {

		return err
	}

	return gr.Close()
}
