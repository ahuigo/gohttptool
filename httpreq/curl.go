package httpreq

import (
	"bytes"
	"io"
	"net/http"
	"net/http/cookiejar"

	"net/url"
	"strings"

	"github.com/ahuigo/gohttptool/shell"
)

func (r *RequestBuilder) FromCurl(curl string) {

}

func (r *RequestBuilder) ToCurl() (curl string, err error) {
	if httpreq, err := r.ToRequest(); err != nil {
		return "", err
	} else {
		curl := buildCurlRequest(httpreq, nil)
		return curl, nil
	}
}

func buildCurlRequest(req *http.Request, httpCookiejar http.CookieJar) (curlString string) {
	buf := acquireBuffer()
	defer releaseBuffer(buf)

	// 1. Generate curl raw headers
	buf.WriteString("curl -X " + req.Method + " ")
	// req.Host + req.URL.Path + "?" + req.URL.RawQuery + " " + req.Proto + " "
	headers := dumpCurlHeaders(req)
	for _, kv := range *headers {
		buf.WriteString(`-H ` + shell.Quote(kv[0]+": "+kv[1]) + ` `)
	}

	// 2. Generate curl cookies
	if cookieJar, ok := httpCookiejar.(*cookiejar.Jar); ok {
		cookies := cookieJar.Cookies(req.URL)
		if len(cookies) > 0 {
			buf.WriteString(` -H ` + shell.Quote(dumpCurlCookies(cookies)) + " ")
		}
	}

	// 3. Generate curl body
	if req.Body != nil {
		buf2, _ := io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewBuffer(buf2)) // important!!
		buf.WriteString(`-d ` + shell.Quote(string(buf2)))
	}

	urlString := shell.Quote(req.URL.String())
	if urlString == "''" {
		urlString = "'http://unexecuted-request'"
	}
	buf.WriteString(" " + urlString)
	return buf.String()
}

// dumpCurlCookies dumps cookies to curl format
func dumpCurlCookies(cookies []*http.Cookie) string {
	sb := strings.Builder{}
	sb.WriteString("Cookie: ")
	for _, cookie := range cookies {
		sb.WriteString(cookie.Name + "=" + url.QueryEscape(cookie.Value) + "&")
	}
	return strings.TrimRight(sb.String(), "&")

}

// dumpCurlHeaders dumps headers to curl format
func dumpCurlHeaders(req *http.Request) *[][2]string {
	headers := [][2]string{}
	for k, vs := range req.Header {
		for _, v := range vs {
			headers = append(headers, [2]string{k, v})
		}
	}
	n := len(headers)
	for i := 0; i < n; i++ {
		for j := n - 1; j > i; j-- {
			jj := j - 1
			h1, h2 := headers[j], headers[jj]
			if h1[0] < h2[0] {
				headers[jj], headers[j] = headers[j], headers[jj]
			}
		}
	}
	return &headers
}
