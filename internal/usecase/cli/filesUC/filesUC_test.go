package filesUC

import (
	"bufio"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"passkeeper/internal/entities/compress"
	"passkeeper/internal/entities/config"
	"passkeeper/internal/entities/structs"
)

func TestFilesUC_UnzipAndDownload(t *testing.T) {
	filename1 := "testfile1"
	file, err := os.CreateTemp(config.TmpDir, filename1)
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

	filename2 := "testfile2"
	file2, err := os.CreateTemp(config.TmpDir, filename2)
	if err != nil {
		t.Errorf("could not create temporary file: %v", err)
		return
	}
	defer os.Remove(file2.Name())

	buf2 := bufio.NewWriter(file2)
	buf2.WriteString(strings.Repeat("222222errerr", 1000))
	err = buf2.Flush()
	if err != nil {
		t.Errorf("could not write to temporary file: %v", err)
		return
	}

	err = file2.Close()
	if err != nil {
		t.Errorf("could not close temporary file: %v", err)
		return
	}

	filename3 := "testfile3"
	file3, err := os.CreateTemp(config.TmpDir, filename3)
	if err != nil {
		t.Errorf("could not create temporary file: %v", err)
		return
	}
	defer os.Remove(file3.Name())

	buf3 := bufio.NewWriter(file3)
	buf3.WriteString(strings.Repeat("3333abc", 1000))
	err = buf3.Flush()
	if err != nil {
		t.Errorf("could not write to temporary file: %v", err)
		return
	}

	err = file3.Close()
	if err != nil {
		t.Errorf("could not close temporary file: %v", err)
		return
	}

	zipData1, err := compress.CompressFile(file.Name())
	if err != nil {
		t.Errorf(err.Error())
	}

	zipData2, err := compress.CompressFile(file2.Name())
	if err != nil {
		t.Errorf(err.Error())
	}

	zipData3, err := compress.CompressFile(file3.Name())
	if err != nil {
		t.Errorf(err.Error())
	}

	var (
		tmpFile1 = &structs.File{
			ID:   "FILE_ID_1",
			Type: structs.BlobFile,
			Name: GenFileShortName(filename1),
			Body: zipData1,
		}
		tmpFile2 = &structs.File{
			ID:   "FILE_ID_2",
			Type: structs.BlobFile,
			Name: GenFileShortName(filename2),
			Body: zipData2,
		}
		tmpFile3 = &structs.File{
			ID:   "FILE_ID_3",
			Type: structs.BlobFile,
			Name: GenFileShortName(filename3),
			Body: zipData3,
		}
	)

	type args struct {
		dirToDownload string
		file          *structs.File
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test Unzip 1: valid",
			args: args{
				dirToDownload: config.TmpDir,
				file:          tmpFile1,
			},
			wantErr: false,
		},
		{
			name: "Test Unzip 2: valid",
			args: args{
				dirToDownload: config.TmpDir,
				file:          tmpFile2,
			},
			wantErr: false,
		},
		{
			name: "Test Unzip 3: valid",
			args: args{
				dirToDownload: config.TmpDir,
				file:          tmpFile3,
			},
			wantErr: false,
		},
		{
			name: "Test Unzip 4: invalid_path",
			args: args{
				dirToDownload: config.TmpDir + "_invalid",
				file:          tmpFile3,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FilesUC{}
			if err := f.UnzipAndDownload(tt.args.dirToDownload, tt.args.file); (err != nil) != tt.wantErr {
				t.Errorf("UnzipAndDownload() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFilesUC_ZipAndUpload(t *testing.T) {

	filename1 := "testfile1"
	file, err := os.CreateTemp(config.TmpDir, filename1)
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

	filename2 := "testfile2"
	file2, err := os.CreateTemp(config.TmpDir, filename2)
	if err != nil {
		t.Errorf("could not create temporary file: %v", err)
		return
	}
	defer os.Remove(file2.Name())

	buf2 := bufio.NewWriter(file2)
	buf2.WriteString(strings.Repeat("222222errerr", 1000))
	err = buf2.Flush()
	if err != nil {
		t.Errorf("could not write to temporary file: %v", err)
		return
	}

	err = file2.Close()
	if err != nil {
		t.Errorf("could not close temporary file: %v", err)
		return
	}

	filename3 := "testfile3"
	file3, err := os.CreateTemp(config.TmpDir, filename3)
	if err != nil {
		t.Errorf("could not create temporary file: %v", err)
		return
	}
	defer os.Remove(file3.Name())

	buf3 := bufio.NewWriter(file3)
	buf3.WriteString(strings.Repeat("3333abc", 1000))
	err = buf3.Flush()
	if err != nil {
		t.Errorf("could not write to temporary file: %v", err)
		return
	}

	err = file3.Close()
	if err != nil {
		t.Errorf("could not close temporary file: %v", err)
		return
	}

	zipData1, err := compress.CompressFile(file.Name())
	if err != nil {
		t.Errorf(err.Error())
	}

	zipData2, err := compress.CompressFile(file2.Name())
	if err != nil {
		t.Errorf(err.Error())
	}

	zipData3, err := compress.CompressFile(file3.Name())
	if err != nil {
		t.Errorf(err.Error())
	}

	var (
		tmpFile1 = &structs.File{
			Body: zipData1,
		}
		tmpFile2 = &structs.File{
			Body: zipData2,
		}
		tmpFile3 = &structs.File{
			Body: zipData3,
		}
	)

	tests := []struct {
		name          string
		uploadingFile string
		wantFile      *structs.File
		wantErr       bool
	}{
		{
			name:          "Test 1: valid tmp1",
			uploadingFile: file.Name(),
			wantFile:      tmpFile1,
			wantErr:       false,
		},
		{
			name:          "Test 2: valid tmp2",
			uploadingFile: file2.Name(),
			wantFile:      tmpFile2,
			wantErr:       false,
		},
		{
			name:          "Test 3: valid tmp3",
			uploadingFile: file3.Name(),
			wantFile:      tmpFile3,
			wantErr:       false,
		},
		{
			name:          "Test 3: invalid file",
			uploadingFile: file3.Name() + "_invalid",
			wantErr:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FilesUC{}
			gotFile, err := f.ZipAndUpload(tt.uploadingFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("ZipAndUpload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			require.Equal(t, tt.wantFile.Body, gotFile.Body)
			require.Equal(t, structs.BlobFile, gotFile.Type)
			require.NotEmpty(t, gotFile.ID)
			require.NotEmpty(t, gotFile.Name)

		})
	}
}

func TestGenFileShortName(t *testing.T) {

	tests := []struct {
		name     string
		filePath string
		want     string
	}{
		{
			name:     "Test 1: simple name",
			filePath: "12345678",
			want:     "12345678",
		},
		{
			name:     "Test 2: large name",
			filePath: strings.Repeat("0123456789-", 4),
			want:     strings.Repeat("0123456789-", 3) + "012" + " ...",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenFileShortName(tt.filePath); got != tt.want {
				t.Errorf("GenFileShortName() = %v, want %v", got, tt.want)
			}
		})
	}
}
