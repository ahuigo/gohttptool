package httpreq

import (
	"strings"
	"testing"

	"github.com/ahuigo/gohttptool/httpreq"
)

type MapString = map[string]string

// POST: curl -X POST -d ” 'http://local/post?p=1'
func TestPostParams(t *testing.T) {
	curl, err := httpreq.R().
		SetQueryParams(map[string]string{"p": "1"}).
		SetReq("POST", "http://local/post").
		GenCurlCommand()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(curl, "curl ") ||
		!strings.Contains(curl, "?p=1") {
		t.Fatal("bad curl: ", curl)
	} else {
		t.Log(curl)
	}
}

// POST: curl -X POST -H 'Content-Type: application/json' -d '{"name":"ahuigo"}' http://localhost/path
func TestPostJson(t *testing.T) {
	curl, err := httpreq.R().
		SetJson(map[string]string{"name": "ahuigo"}).
		SetReq("POST", "/path").
		GenCurlCommand()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(curl, "curl ") ||
		!strings.Contains(curl, `'Content-Type: application/json`) ||
		!strings.Contains(curl, `{"name":"ahuigo"}`) {
		t.Fatal("bad curl: ", curl)
	} else {
		t.Log(curl)
	}
}

// Post Data: curl -H 'Content-Type: application/x-www-form-urlencoded' http://local/post -d 'name=Alex'
func TestPostFormUrlEncode(t *testing.T) {
	curl, err := httpreq.R().
		SetFormData(map[string]string{"name": "Alex"}).
		SetReq("POST", "http://local/post").
		GenCurlCommand()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(curl, "curl ") ||
		!strings.Contains(curl, `'Content-Type: application/x-www-form-urlencode`) ||
		!strings.Contains(curl, "name=Alex") {
		t.Fatal("bad curl: ", curl)
	} else {
		t.Log(curl)
	}

}

// POST FormData: multipart/form-data; boundary=....
// curl http://local/post -F 'name=Alex'
func TestPostFormMultipart(t *testing.T) {
	curl, err := httpreq.R().
		SetIsMultiPart(true).
		SetFormData(map[string]string{"name": "Alex"}).
		SetReq("POST", "http://local/post").
		GenCurlCommand()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(curl, "curl ") ||
		!strings.Contains(curl, `'Content-Type: multipart/form-data`) ||
		!strings.Contains(curl, `name="name"`) {
		t.Fatal("bad curl: ", curl)
	} else {
		println(curl)
	}
}

// POST: curl -X POST -H 'Content-Type: text/plain' -d 'hello!' http://local/post
func TestPostPlainData(t *testing.T) {
	curl, err := httpreq.R().
		SetContentType(httpreq.ContentTypePlain).
		SetBody([]byte("hello!")).
		SetReq("POST", "http://local/post").
		GenCurlCommand()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(curl, "curl ") ||
		!strings.Contains(curl, `'Content-Type: text/plain`) ||
		!strings.Contains(curl, `-d 'hello!'`) {
		t.Fatal("bad curl: ", curl)
	} else {
		println(curl)
	}
}
