# 🛠️ Go http tool
[![tag](https://img.shields.io/github/tag/ahuigo/gohttptool.svg)](https://github.com/ahuigo/gohttptool/tags)
![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.21-%23007d9c)
[![GoDoc](https://godoc.org/github.com/ahuigo/gohttptool?status.svg)](https://pkg.go.dev/github.com/ahuigo/gohttptool)
![Build Status](https://github.com/ahuigo/gohttptool/actions/workflows/test.yml/badge.svg)
[![Go report](https://goreportcard.com/badge/github.com/ahuigo/gohttptool)](https://goreportcard.com/report/github.com/ahuigo/gohttptool)
[![Coverage](https://img.shields.io/codecov/c/github/ahuigo/gohttptool)](https://codecov.io/gh/ahuigo/gohttptool)
[![Contributors](https://img.shields.io/github/contributors/ahuigo/gohttptool)](https://github.com/ahuigo/gohttptool/graphs/contributors)
[![License](https://img.shields.io/github/license/ahuigo/gohttptool)](./LICENSE)

## Features
- [x] Build http request in golang
- [x] Generate curl command for http request

## Unittest with gonic

```

func CreateTestCtx(req *http.Request) (resp *httptest.ResponseRecorder, ctx *gin.Context) {
	resp = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(resp)
	ctx.Request = req
	return
}

func TestGonicApi(t *testing.T) {
	// 1. build request
	req, _ := httpreq.R().
		SetQueryParams(map[string]string{
			"job_id":   "1234",
		}).
		SetReq("GET", "http://any/api/v1/spark/job").
		GenRequest()
	curl := httpreq.GenCurlCommand(req, nil)
	println(curl)
	resp, ctx := CreateTestCtx(req)

	// 2. execute
	sparkServer := GetGonicSparkServer()
	sparkServer.GetJobInfo(ctx)
	if resp.Code != http.StatusOK {
		errors := ctx.Errors.Errors()
		fmt.Println("output", errors)
		t.Errorf("Expect code 200, but get %d body:%v", resp.Code, resp.Body)
	} else {
        data := map[string]string{}
		httpreq.BuildResponse(resp.Result()).Json(&data)
		if data["status"] == "" {
			t.Fatalf("Bad response: %v", data)
		}
	}
}
```
