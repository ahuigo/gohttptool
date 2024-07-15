package filetool

import (
	"bytes"
	"errors"
	"io"
	"mime"
	"mime/multipart"
	"os"
	"path/filepath"
)

func CreateFileHeaderFromFile(filePath string) (*multipart.FileHeader, error) {
	filename := filepath.Base(filePath)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	if content, err := io.ReadAll(file); err == nil {
		return CreateFileHeaderFromBytes(filename, content)
	} else {
		return nil, err
	}
}

func CreateFileHeaderFromBytes(filename string, content []byte) (*multipart.FileHeader, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file1", filename)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(part, bytes.NewReader(content))
	if err != nil {
		return nil, err
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	_, params, err := mime.ParseMediaType(writer.FormDataContentType())
	if err != nil {
		return nil, err
	}

	boundary, ok := params["boundary"]
	if !ok {
		return nil, errors.New("no boundary")
	}

	reader := multipart.NewReader(body, boundary)
	mf, _ := reader.ReadForm(1 << 8)
	fileHeader := mf.File["file1"][0]

	return fileHeader, nil
}
