package httpreq

import (
	"strings"
	"testing"

	"github.com/ahuigo/gohttptool/httpreq"
)

func TestPostFile(t *testing.T) {
	curl, err := httpreq.R().
		SetQueryParams(map[string]string{"p": "1"}).
		SetFormData(map[string]string{"key": "xx"}).
		SetAuthBasic("user", "pass").
		SetHeader("header1", "value1").
		AddCookieKV("count", "1").
		AddFileHeader("file", "test.txt", []byte("hello world")).
		AddFile("file2", getTestDataPath("a.txt")).
		SetReq("GET", "/path").
		ToCurl()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(curl, "curl ") {
		t.Fatal("bad curl: ", curl)
	} else {
		t.Log(curl)
	}
}
