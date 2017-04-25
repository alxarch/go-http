package httpdump

import (
	"encoding/json"
	"mime"
	"net/http"
	"regexp"
)

type Request []byte

type request struct {
	Path    string
	Method  string
	Headers http.Header
	Body    json.RawMessage
}

var requestLineRX = regexp.MustCompile("^([A-Z]+) (/.*) HTTP/1\\.1$")

func (r Request) MarshalJSON() ([]byte, error) {
	rline, headers, body, err := Parse(r)
	if err != nil {
		return nil, err
	}
	req := &request{}
	if match := requestLineRX.FindStringSubmatch(rline); len(match) > 0 {
		req.Method = match[1]
		req.Path = match[2]
	}
	req.Headers = headers
	if body != nil {
		var tmp interface{}
		tmp = body
		if mediaType := req.Headers.Get("Content-Type"); mediaType != "" {
			if m, _, err := mime.ParseMediaType(mediaType); err == nil {
				switch m {
				case "application/json":
					tmp = json.RawMessage(body)
				case "text/plain", "text/html", "text/xml":
					tmp = string(body)
				}
			}
		}
		if msg, err := json.Marshal(tmp); err != nil {
			return nil, err
		} else {
			req.Body = json.RawMessage(msg)
		}

	}
	return json.Marshal(req)
}
