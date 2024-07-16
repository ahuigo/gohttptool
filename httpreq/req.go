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

type RequestBuilder struct {
	rawreq      *http.Request
	url         string
	files       map[string]string     // field -> path
	fileHeaders map[string]fileHeader // field -> contents
	datas       map[string]string     // key -> value
	params      map[string]string     // key -> value
	paramsList  map[string][]string   // key -> value list
}

func R() *RequestBuilder {
	return &RequestBuilder{
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
func (r *RequestBuilder) SetHeader(key, value string) *RequestBuilder {
	r.rawreq.Header.Set(key, value)
	return r
}

func (r *RequestBuilder) SetAuthBasic(username, password string) *RequestBuilder {
	r.rawreq.SetBasicAuth(username, password)
	return r
}

func (r *RequestBuilder) SetAuthBearer(token string) *RequestBuilder {
	r.rawreq.Header.Set("Authorization", "Bearer "+token)
	return r
}

func (r *RequestBuilder) AddCookies(cookies []*http.Cookie) *RequestBuilder {
	for _, cookie := range cookies {
		r.rawreq.AddCookie(cookie)
	}
	return r
}
func (r *RequestBuilder) AddCookieKV(name, value string) *RequestBuilder {
	cookie := &http.Cookie{
		Name:  name,
		Value: value,
	}
	r.rawreq.AddCookie(cookie)
	return r
}

/************** file **********************/
func (r *RequestBuilder) AddFile(fieldname, path string) *RequestBuilder {
	r.files[fieldname] = path
	return r
}

func (r *RequestBuilder) AddFileHeader(fieldname, filename string, content []byte) *RequestBuilder {
	r.fileHeaders[fieldname] = fileHeader{
		Filename: filename,
		content:  content,
		Size:     int64(len(content)),
	}
	return r
}

func (r *RequestBuilder) SetUrl(url string) *RequestBuilder {
	r.url = url
	return r
}

func (r *RequestBuilder) SetReq(method string, url string) *RequestBuilder {
	r.rawreq.Method = method
	r.url = url
	return r
}

func (r *RequestBuilder) SetParams(params map[string]string) *RequestBuilder {
	r.params = params
	return r
}

func (r *RequestBuilder) SetData(data map[string]string) *RequestBuilder {
	r.datas = data
	return r
}

func (r *RequestBuilder) GetRawreq() *http.Request {
	return r.rawreq
}

func (r *RequestBuilder) SetCtx(ctx context.Context) *RequestBuilder {
	r.rawreq = r.rawreq.WithContext(ctx)
	return r
}

func (r *RequestBuilder) EnableTrace(ctx context.Context) *RequestBuilder {
	trace := clientTraceNew(r.rawreq.Context())
	r.rawreq = r.rawreq.WithContext(trace.ctx)
	return r
}
