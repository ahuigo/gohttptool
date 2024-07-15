package httpreq

import (
	"context"
	"net/http"
	"net/textproto"
)

type ContentType string

const (
	ContentTypeNone       ContentType = ""
	ContentTypeFormEncode ContentType = "application/x-www-form-urlencoded"
	ContentTypeFormData   ContentType = "multipart/form-data"
	ContentTypeJson       ContentType = "application/json"
	ContentTypePlain      ContentType = "text/plain"
)

type fileHeader struct {
	Filename string
	Header   textproto.MIMEHeader
	Size     int64
	content  []byte
	// tmpfile   string
	// tmpoff    int64
	// tmpshared bool
}

type request struct {
	rawreq      *http.Request
	url         string
	files       map[string]string     // field -> path
	fileHeaders map[string]fileHeader // field -> contents
	datas       map[string]string     // key -> value
	params      map[string]string     // key -> value
	paramsList  map[string][]string   // key -> value list
}

func R() *request {
	return &request{
		rawreq: &http.Request{
			Method:     "GET",
			Header:     make(http.Header),
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
		},
		files:       make(map[string]string),
		fileHeaders: make(map[string]fileHeader),
		datas:       make(map[string]string),
		params:      make(map[string]string),
		paramsList:  make(map[string][]string),
	}
}

/******************header *************************/
func (r *request) SetHeader(key, value string) *request {
	r.rawreq.Header.Set(key, value)
	return r
}

func (r *request) SetAuthBasic(username, password string) *request {
	r.rawreq.SetBasicAuth(username, password)
	return r
}

func (r *request) SetAuthBearer(token string) *request {
	r.rawreq.Header.Set("Authorization", "Bearer "+token)
	return r
}

func (r *request) AddCookies(cookies []*http.Cookie) *request {
	for _, cookie := range cookies {
		r.rawreq.AddCookie(cookie)
	}
	return r
}
func (r *request) AddCookieKV(name, value string) *request {
	cookie := &http.Cookie{
		Name:  name,
		Value: value,
	}
	r.rawreq.AddCookie(cookie)
	return r
}

/************** file **********************/
func (r *request) AddFile(fieldname, path string) *request {
	r.files[fieldname] = path
	return r
}

func (r *request) AddFileHeader(fieldname, filename string, content []byte) *request {
	r.fileHeaders[fieldname] = fileHeader{
		Filename: filename,
		content:  content,
		Size:     int64(len(content)),
	}
	return r
}

func (r *request) SetUrl(url string) *request {
	r.url = url
	return r
}

func (r *request) SetReq(method string, url string) *request {
	r.rawreq.Method = method
	r.url = url
	return r
}

func (r *request) SetParams(params map[string]string) *request {
	r.params = params
	return r
}

func (r *request) SetData(data map[string]string) *request {
	r.datas = data
	return r
}

func (r *request) GetRawreq() *http.Request {
	return r.rawreq
}

func (r *request) SetCtx(ctx context.Context) *request {
	r.rawreq = r.rawreq.WithContext(ctx)
	return r
}

func (r *request) EnableTrace(ctx context.Context) *request {
	trace := clientTraceNew(r.rawreq.Context())
	r.rawreq = r.rawreq.WithContext(trace.ctx)
	return r
}
