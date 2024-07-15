package httpreq

import (
	"strings"
	"testing"

	"github.com/ahuigo/gohttptool/httpreq"
)

func TestCurl(t *testing.T) {
	curl, err := httpreq.R().
		SetParams(map[string]string{"p": "1"}).
		AddCookieKV("count", "1").
		AddFileHeader("file", "test.txt", []byte("hello world")).
		ToCurl()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(curl, "curl ") {
		t.Fatal("bad curl: ", curl)
	}else{
		t.Log("curl: ", curl)
	}
}
