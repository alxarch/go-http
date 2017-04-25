package httpdump_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"testing"

	"github.com/alxarch/go-http/httpdump"
)

func Test_MarshalJSON(t *testing.T) {
	b := bytes.NewBuffer([]byte(`{"foo": "bar"}`))
	req := httptest.NewRequest(http.MethodPost, "/", b)
	req.Header.Set("Content-Type", "application/json")
	data, _ := httputil.DumpRequest(req, true)
	d := httpdump.Request(data)
	out, err := d.MarshalJSON()
	if err != nil {
		t.Error(err)
	}
	if string(out) != `{"Path":"/","Method":"POST","Headers":{"Content-Type":["application/json"],"Host":["example.com"]},"Body":{"foo":"bar"}}` {
		t.Error("Not equal")
	}
}
