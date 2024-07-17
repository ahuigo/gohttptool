package httpreq

import (
	"bytes"
	"encoding/json"
	"sync"
)

var (
	bufPool = &sync.Pool{New: func() interface{} { return &bytes.Buffer{} }}
)

var noescapeJSONMarshalIndent = func(v interface{}) (*bytes.Buffer, error) {
	buf := acquireBuffer()
	encoder := json.NewEncoder(buf)
	encoder.SetEscapeHTML(false)
	// encoder.SetIndent("", "   ")

	if err := encoder.Encode(v); err != nil {
		releaseBuffer(buf)
		return nil, err
	}
	return buf, nil
}

func acquireBuffer() *bytes.Buffer {
	return bufPool.Get().(*bytes.Buffer)
}

func releaseBuffer(buf *bytes.Buffer) {
	if buf != nil {
		buf.Reset()
		bufPool.Put(buf)
	}
}
