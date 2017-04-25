package httpdump

import (
	"encoding/json"
	"mime"
	"net/http"
	"regexp"
	"strconv"
)

type Response []byte

type response struct {
	StatusCode    int
	StatusMessage string
	Headers       http.Header
	Body          json.RawMessage
}

var responseLineRX = regexp.MustCompile("^HTTP/\\d+\\.\\d+ (\\d+) (.*)$")

func (r Response) MarshalJSON() ([]byte, error) {
	rline, headers, body, err := Parse(r)
	if err != nil {
		return nil, err
	}
	res := &response{}
	match := responseLineRX.FindStringSubmatch(rline)
	if len(match) > 0 {
		if code, err := strconv.Atoi(match[1]); err != nil {
			return nil, err
		} else {
			res.StatusCode = code
		}
		res.StatusMessage = match[2]
	}
	res.Headers = headers
	if body != nil {
		var tmp interface{}
		tmp = body
		if mediaType := res.Headers.Get("Content-Type"); mediaType != "" {
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
			res.Body = json.RawMessage(msg)
		}

	}
	return json.Marshal(res)
}
