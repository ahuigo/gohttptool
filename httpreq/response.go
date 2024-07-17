package httpreq

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

type httpRes struct {
	R              *http.Response
	Attempt        int
	body           []byte
	doNotCloseBody bool
	URL            *url.URL
	client         *http.Client
	isdebugBody    bool
	dumpCurl       string
	dumpResponse   string
}

func BuildResponse(response *http.Response) *httpRes {
	r := &httpRes{
		R: response,
	}
	// resp.R.Body = ioutil.NopCloser(bytes.NewBuffer(resp.Body())) // important!!
	// r._DumpResponse(true)
	// r.Body()
	return r
}

func (resp *httpRes) SetDoNotCloseBody() *httpRes {
	resp.doNotCloseBody = true
	return resp
}

func (resp *httpRes) SetClientReq(url *url.URL, client *http.Client) *httpRes {
	resp.client = client
	resp.URL = url
	return resp
}

func (resp *httpRes) ResponseDebug() {
	fmt.Println("===========ResponseDebug ============")
	err := resp._DumpResponse(resp.isdebugBody)
	if err != nil {
		return
	}
	fmt.Println(resp.dumpResponse)
	fmt.Println("========== ResponseDebug(end) ============")
}

func (resp *httpRes) _DumpResponse(isdebugBody bool) error {
	message, err := httputil.DumpResponse(resp.R, isdebugBody)
	resp.dumpResponse = string(message)
	return err
}

func (resp *httpRes) GetDumpCurl() string {
	return resp.dumpCurl
}
func (resp *httpRes) GetDumpResponse() string {
	if resp.dumpResponse == "" {
		if resp.R.Body == nil {
			resp.R.Body = io.NopCloser(bytes.NewBuffer(resp.Body())) // important!!
		}
		resp._DumpResponse(true)
	}
	return resp.dumpResponse
}

func (resp *httpRes) Body() []byte {
	var err error
	if resp.body != nil {
		return resp.body
	}
	resp.body = []byte{}
	if !resp.doNotCloseBody {
		defer resp.R.Body.Close()
	}

	var Body = resp.R.Body
	if resp.R.Header.Get("Content-Encoding") == "gzip" {
		reader, err := gzip.NewReader(Body)
		if err != nil {
			return nil
		}
		Body = reader
	}

	resp.body, err = io.ReadAll(Body)
	if err != nil {
		return nil
	}

	return resp.body
}

func (resp *httpRes) Text() string {
	return string(resp.Body())
}

func (resp *httpRes) Size() int {
	return len(resp.Body())
}

func (resp *httpRes) RaiseForStatus() (code int, err error) {
	code = resp.R.StatusCode
	if resp.R.StatusCode >= 400 && resp.R.StatusCode != 401 {
		err = errors.New(resp.Text())
	}
	return
}

func (resp *httpRes) StatusCode() (code int) {
	return resp.R.StatusCode
}

func (resp *httpRes) Header() http.Header {
	return resp.R.Header
}

func (resp *httpRes) SaveFile(filename string) error {
	if resp.body == nil {
		resp.Body()
	}
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(resp.body)
	f.Sync()

	return err
}

func (resp *httpRes) Json(v interface{}) error {
	if resp.body == nil {
		resp.Body()
	}
	return json.Unmarshal(resp.body, v)
}

func (resp *httpRes) Cookies() (cookies []*http.Cookie) {
	client := resp.client

	if resp.URL == nil || client == nil {
		return resp.R.Cookies()
	}
	// cookies's type is `[]*http.Cookies`
	cookies = client.Jar.Cookies(resp.URL)
	return cookies
}

func (resp *httpRes) GetCookie(key string) (val string) {
	cookies := map[string]string{}
	for _, c := range resp.Cookies() {
		cookies[c.Name] = c.Value
	}
	val = cookies[key]
	return val
}

func (resp *httpRes) HasCookie(key string) (exists bool) {
	cookies := map[string]string{}
	for _, c := range resp.Cookies() {
		cookies[c.Name] = c.Value
	}
	_, exists = cookies[key]
	return exists
}
