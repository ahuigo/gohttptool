package httpreq

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/textproto"
	"net/url"
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
	rawreq *http.Request
	url    string

	queryParam  url.Values
	formData    url.Values
	isMultiPart bool
	json        any
	files       map[string]string     // field -> path
	fileHeaders map[string]fileHeader // field -> contents
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
		queryParam: make(url.Values),
		formData:   make(map[string][]string),
		// paramsList:  make(map[string][]string),
		files:       make(map[string]string),
		fileHeaders: make(map[string]fileHeader),
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

func (r *RequestBuilder) SetContentType(ct ContentType) *RequestBuilder {
	r.rawreq.Header.Set("Content-Type", string(ct))
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

/************** params **********************/
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

/************** params **********************/
func (r *RequestBuilder) SetQueryParams(params map[string]string) *RequestBuilder {
	for p, v := range params {
		r.SetQueryParam(p, v)
	}
	return r
}
func (r *RequestBuilder) SetQueryParam(param, value string) *RequestBuilder {
	r.queryParam.Set(param, value)
	return r
}

func (r *RequestBuilder) SetQueryParamsFromValues(params url.Values) *RequestBuilder {
	for p, v := range params {
		for _, pv := range v {
			r.queryParam.Add(p, pv)
		}
	}
	return r
}

/************** body(bytes) **********************/
func (r *RequestBuilder) SetBody(body []byte) *RequestBuilder {
	r.rawreq.Body = io.NopCloser(bytes.NewReader(body))
	return r
}

/************** body(form) **********************/
// Set Form data(encode or multipart)
func (r *RequestBuilder) SetIsMultiPart(b bool) *RequestBuilder {
	r.isMultiPart = b
	return r
}
func (r *RequestBuilder) SetFormData(data map[string]string) *RequestBuilder {
	for k, v := range data {
		r.formData.Set(k, v)
	}
	return r
}

// SetFormDataFromValues method appends multiple form parameters with multi-value
//
//	SetFormDataFromValues(url.Values{"words": []string{"book", "glass", "pencil"},})
func (r *RequestBuilder) SetFormDataFromValues(data url.Values) *RequestBuilder {
	for k, v := range data {
		for _, kv := range v {
			r.formData.Add(k, kv)
		}
	}
	return r
}

/************** body(json) **********************/
func (r *RequestBuilder) SetJson(data any) *RequestBuilder {
	r.json = data
	return r
}

/************** body(plain) **********************/

/************** utils **********************/
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
