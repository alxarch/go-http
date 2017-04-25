package httpdump_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"testing"

	"github.com/alxarch/go-http/httpdump"
)

func Test_ResponseMarshalJSON(t *testing.T) {
	b := bytes.NewBuffer([]byte(`{"foo": "bar"}`))
	res := &http.Response{}
	res.Status = "OK"
	res.StatusCode = 200
	res.Body = ioutil.NopCloser(b)
	res.Header = http.Header{}
	res.Header.Set("Content-Type", "application/json")
	data, _ := httputil.DumpResponse(res, true)
	d := httpdump.Response(data)
	out, err := d.MarshalJSON()
	if err != nil {
		t.Error(err)
	}
	if string(out) != `{"StatusCode":200,"StatusMessage":"OK","Headers":{"Content-Type":["application/json"]},"Body":{"foo":"bar"}}` {
		t.Error("Not equal")
	}
}
