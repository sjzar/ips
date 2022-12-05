package rewriter

import (
	"net"

	"github.com/sjzar/ips/ipx"
)

type Rewriter interface {

	// Rewrite 改写
	Rewrite(net.IP, *ipx.Range, map[string]string) (net.IP, *ipx.Range, map[string]string, error)
}

var DefaultRewriter Rewriter = &EmptyRewriter{}
