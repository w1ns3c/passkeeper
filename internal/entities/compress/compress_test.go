package compress

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/w1ns3c/go-examples/crypto"

	"passkeeper/internal/entities/config"
)

func TestCompressFile(t *testing.T) {

	file, err := os.CreateTemp(config.TmpDir, "testfile")
	if err != nil {
		t.Errorf("could not create temporary file: %v", err)
		return
	}
	defer os.Remove(file.Name())

	buf := bufio.NewWriter(file)
	buf.WriteString(strings.Repeat("123123", 1000))
	err = buf.Flush()
	if err != nil {
		t.Errorf("could not write to temporary file: %v", err)
		return
	}

	err = file.Close()
	if err != nil {
		t.Errorf("could not close temporary file: %v", err)
		return
	}

	info, err := os.Stat(file.Name())
	if err != nil {
		t.Errorf("could not get temporary file info: %v", err)
		return
	}

	tests := []struct {
		name     string
		filepath string
		startLen int64
		wantLen  int
		wantErr  bool
	}{
		{
			name:     "Test compress1: tmpFile 1",
			filepath: file.Name(),
			startLen: info.Size(),
			wantLen:  48,
			wantErr:  false,
		},
		{
			name:     "Test compress 2: folder not exist",
			filepath: file.Name() + "tmptmtp",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotData, err := CompressFile(tt.filepath)
			if (err != nil) != tt.wantErr {
				t.Errorf("CompressFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			if int64(len(gotData)) >= tt.startLen {
				t.Errorf("CompressFile() compressed file too large = %v, wantErr %v", tt.startLen, len(gotData))
				return
			}

			if len(gotData) != tt.wantLen {
				t.Errorf("CompressFile() len(gotData) = %v, want %v", len(gotData), tt.wantLen)
			}
		})
	}
}

func TestDecompressFile(t *testing.T) {
	type args struct {
		data     []byte
		filePath string
	}

	randSuffix, _ := crypto.GenRandStr(5)
	filePathIn := filepath.Join(config.TmpDir, "test_file_in_"+randSuffix)
	filePathOut := filepath.Join(config.TmpDir, "test_file_out_"+randSuffix)
	filePathNotExist := filepath.Join(config.TmpDir, "test_file_in_", randSuffix)

	file, err := os.Create(filePathIn)
	if err != nil {
		t.Errorf("could not create temporary file: %v", err)
		return
	}
	defer os.Remove(file.Name())

	buf := bufio.NewWriter(file)
	buf.WriteString(strings.Repeat("123123", 1000))
	err = buf.Flush()
	if err != nil {
		t.Errorf("could not write to temporary file: %v", err)
		return
	}

	err = file.Close()
	if err != nil {
		t.Errorf("could not close temporary file: %v", err)
		return
	}

	data, err := CompressFile(filePathIn)
	if err != nil {
		t.Errorf("could not compress temp file: %v", err)
	}

	defer os.Remove(filePathIn)
	defer os.Remove(filePathOut)

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test decompress: tmpFile 1",
			args: args{
				data:     data,
				filePath: filePathOut,
			},
			wantErr: false,
		},
		{
			name: "Test decompress: folder not exist",
			args: args{
				data:     data,
				filePath: filePathNotExist,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DecompressFile(tt.args.data, tt.args.filePath); (err != nil) != tt.wantErr {
				t.Errorf("DecompressFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
