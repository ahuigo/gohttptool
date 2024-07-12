package demo

import (
	"bytes"
	"io"
	"testing"

	"github.com/ahuigo/gohttptool/filetool"
)

func TestCreateFileHeader(t *testing.T) {
	content := []byte("hello world")
	fd, err := filetool.CreateFileHeaderFromBytes("test.txt", content)
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
