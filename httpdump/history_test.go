package httpdump_test

import (
	"testing"

	"github.com/alxarch/go-http/httpdump"
)

func Test_Ring(t *testing.T) {
	h := httpdump.NewHistory(10)
	h.Do(func(req httpdump.Request, res httpdump.Response) {

	})

}
