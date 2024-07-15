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

func (session *request) ToRequest() (*http.Request, error) {
	var dataType = ContentType(session.rawreq.Header.Get("Content-Type"))
	var origurl = session.url
	if len(session.files) > 0 || len(session.fileHeaders) > 0 {
		dataType = ContentTypeFormData
	}

	URL, err := session.buildURLParams(origurl)
	if err != nil {
		return nil, err
	}
	if URL.Scheme == "" || URL.Host == "" {
		err = &url.Error{Op: "parse", URL: origurl, Err: fmt.Errorf("failed")}
		return nil, err
	}

	switch dataType {
	case ContentTypeFormEncode:
		if len(session.datas) > 0 {
			formEncodeValues := session.buildFormEncode(session.datas)
			session.setBodyFormEncode(formEncodeValues)
		}
	case ContentTypeFormData:
		// multipart/form-data
		session.buildFilesAndForms()
	}

	if session.rawreq.Body == nil && session.rawreq.Method != "GET" {
		session.rawreq.Body = http.NoBody
	}

	session.rawreq.URL = URL

	return session.rawreq, nil
}

// build post Form encode
func (session *request) buildFormEncode(datas map[string]string) (Forms url.Values) {
	Forms = url.Values{}
	for key, value := range datas {
		Forms.Add(key, value)
	}
	return Forms
}

// set form urlencode
func (session *request) setBodyFormEncode(Forms url.Values) {
	data := Forms.Encode()
	session.rawreq.Body = io.NopCloser(strings.NewReader(data))
	session.rawreq.ContentLength = int64(len(data))
}

func (r *request) buildURLParams(userURL string) (*url.URL, error) {
	params := r.params
	paramsArray := r.paramsList
	if strings.HasPrefix(userURL, "/") {
		userURL = "http://localhost" + userURL
	}else if userURL == ""{
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

func (r *request) buildFilesAndForms() error {
	files := r.files
	datas := r.datas
	filesHeaders := r.fileHeaders
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
	r.rawreq.Body = io.NopCloser(bytes.NewReader(b.Bytes()))
	r.rawreq.ContentLength = int64(b.Len())
	r.rawreq.Header.Set("Content-Type", w.FormDataContentType())
	return nil
}
