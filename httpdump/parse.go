package httpdump

import (
	"bufio"
	"bytes"
	"net/http"
	"regexp"
)

const ScanBufferSize = 8 * 1024

var headersRx = regexp.MustCompile(" *(.+): *([^\\r]+)")

func Parse(dump []byte) (string, http.Header, []byte, error) {

	n := 0
	buffer := bytes.NewBuffer(dump)
	scanner := bufio.NewScanner(buffer)
	scanner.Split(bufio.SplitFunc(func(data []byte, atEOF bool) (int, []byte, error) {
		adv, token, err := bufio.ScanLines(data, atEOF)
		n += adv
		return adv, token, err
	}))
	var rline string
	rlineOK := false

	headers := make(http.Header)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			break
		}
		if !rlineOK {
			rline = line
			rlineOK = true
			continue
		}
		match := headersRx.FindStringSubmatch(line)
		if len(match) > 0 {
			name := match[1]
			value := match[2]
			headers[name] = append(headers[name], value)
		}
	}
	return rline, headers, dump[n:], nil
}
