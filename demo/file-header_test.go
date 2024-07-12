package demo

import (
	"bytes"
	"io"
	"testing"

	"github.com/ahuigo/gohttptool/file"
)

func TestCreateFileHeader(t *testing.T) {
	content := []byte("hello world")
	fd, err := file.CreateFileHeaderFromBytes("test.txt", content)
	if err != nil {
		t.Fatal(err)
	}
	fh, err := fd.Open()
	if err != nil {
		t.Fatal(err)
	}
	r, _ := io.ReadAll(fh)
	if !bytes.Equal(r, content)  {
		t.Fatalf("content not match: %s, %s\n", string(r), string(content))
	}
}
