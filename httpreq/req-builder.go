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
	if len(rb.files) > 0 || len(rb.fileHeaders) > 0 {
		dataType = ContentTypeFormData
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
	case ContentTypeFormEncode:
		if len(rb.datas) > 0 {
			formEncodeValues := rb.buildFormEncode(rb.datas)
			rb.setBodyFormEncode(formEncodeValues)
		}
	case ContentTypeFormData:
		// multipart/form-data
		rb.buildFilesAndForms()
	}

	if rb.rawreq.Body == nil && rb.rawreq.Method != "GET" {
		rb.rawreq.Body = http.NoBody
	}

	rb.rawreq.URL = URL

	return rb.rawreq, nil
}

// build post Form encode
func (rb *RequestBuilder) buildFormEncode(datas map[string]string) (Forms url.Values) {
	Forms = url.Values{}
	for key, value := range datas {
		Forms.Add(key, value)
	}
	return Forms
}

// set form urlencode
func (rb *RequestBuilder) setBodyFormEncode(Forms url.Values) {
	data := Forms.Encode()
	rb.rawreq.Body = io.NopCloser(strings.NewReader(data))
	rb.rawreq.ContentLength = int64(len(data))
}

func (rb *RequestBuilder) buildURLParams(userURL string) (*url.URL, error) {
	params := rb.params
	paramsArray := rb.paramsList
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

	for key, value := range params {
		values.Set(key, value)
	}
	for key, vals := range paramsArray {
		for _, v := range vals {
			values.Add(key, v)
		}
	}
	parsedURL.RawQuery = values.Encode()
	return parsedURL, nil
}

func (rb *RequestBuilder) buildFilesAndForms() error {
	files := rb.files
	datas := rb.datas
	filesHeaders := rb.fileHeaders
	//handle file multipart
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	for k, v := range datas {
		w.WriteField(k, v)
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
