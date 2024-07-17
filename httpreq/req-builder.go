package httpreq

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/pkg/errors"
)

func (rb *RequestBuilder) ToRequest() (*http.Request, error) {
	var dataType = ContentType(rb.rawreq.Header.Get("Content-Type"))
	var origurl = rb.url
	if rb.isMultiPart || len(rb.files) > 0 || len(rb.fileHeaders) > 0 {
		dataType = ContentTypeFormData
	} else if rb.json != nil {
		dataType = ContentTypeJson
	} else if len(rb.formData) > 0 {
		dataType = ContentTypeFormEncode
	}
	if dataType != "" {
		rb.rawreq.Header.Set("Content-Type", string(dataType))
	}

	URL, err := rb.buildURLParams(origurl)
	if err != nil {
		return nil, err
	}
	if URL.Scheme == "" || URL.Host == "" {
		err = &url.Error{Op: "parse", URL: origurl, Err: fmt.Errorf("failed")}
		return nil, err
	}

	switch dataType {
	case ContentTypeJson:
		rb.setBodyJson()
	case ContentTypeFormEncode:
		if len(rb.formData) > 0 {
			rb.setBodyFormEncode(rb.formData)
		}
	case ContentTypeFormData:
		// multipart/form-data
		rb.buildFilesAndForms()
	default:
	}

	if rb.rawreq.Body == nil && rb.rawreq.Method != "GET" {
		rb.rawreq.Body = http.NoBody
	}

	rb.rawreq.URL = URL

	return rb.rawreq, nil
}

// set form urlencode
func (rb *RequestBuilder) setBodyFormEncode(formData url.Values) {
	data := formData.Encode()
	rb.rawreq.Body = io.NopCloser(strings.NewReader(data))
	rb.rawreq.ContentLength = int64(len(data))
}

func (rb *RequestBuilder) setBodyJson() {
	bodyBuf, err := noescapeJSONMarshalIndent(&rb.json)
	if err == nil {
		prtBodyBytes := bodyBuf.Bytes()
		plen := len(prtBodyBytes)
		if plen > 0 && prtBodyBytes[plen-1] == '\n' {
			prtBodyBytes = prtBodyBytes[:plen-1]
		}
		rb.rawreq.Body = io.NopCloser(bytes.NewReader(prtBodyBytes))
	}
}

func (rb *RequestBuilder) buildURLParams(userURL string) (*url.URL, error) {
	if strings.HasPrefix(userURL, "/") {
		userURL = "http://localhost" + userURL
	} else if userURL == "" {
		userURL = "http://unknown"
	}
	parsedURL, err := url.Parse(userURL)

	if err != nil {
		return nil, err
	}

	values := parsedURL.Query()

	for key, value := range rb.queryParam {
		values[key] = value
		// values.Set(key, value[0])
	}
	parsedURL.RawQuery = values.Encode()
	return parsedURL, nil
}

func (rb *RequestBuilder) buildFilesAndForms() error {
	files := rb.files
	formData := rb.formData
	filesHeaders := rb.fileHeaders
	//handle file multipart
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	for k, v := range formData {
		for _, vv := range v {
			w.WriteField(k, vv)
		}
	}

	for field, path := range files {
		part, err := w.CreateFormFile(field, path)
		if err != nil {
			fmt.Printf("Upload %s failed!", path)
			panic(err)
		}
		file, err := os.Open(path)
		if err != nil {
			err = errors.WithMessagef(err, "Open %s", path)
			return err
		}
		_, err = io.Copy(part, file)
		if err != nil {
			return err
		}
	}
	for field, fileheader := range filesHeaders {
		part, err := w.CreateFormFile(field, fileheader.Filename)
		if err != nil {
			fmt.Printf("Upload %s failed!", field)
			panic(err)
		}
		_, err = io.Copy(part, bytes.NewReader([]byte(fileheader.content)))
		if err != nil {
			return err
		}
	}

	w.Close()
	// set file header example:
	// "Content-Type": "multipart/form-data; boundary=------------------------7d87eceb5520850c",
	rb.rawreq.Body = io.NopCloser(bytes.NewReader(b.Bytes()))
	rb.rawreq.ContentLength = int64(b.Len())
	rb.rawreq.Header.Set("Content-Type", w.FormDataContentType())
	return nil
}
