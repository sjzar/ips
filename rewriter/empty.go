package rewriter

import (
	"net"

	"github.com/sjzar/ips/ipx"
)

// EmptyRewriter 空改写
type EmptyRewriter struct {
}

// Rewrite 改写
func (r *EmptyRewriter) Rewrite(ip net.IP, ipRange *ipx.Range, data map[string]string) (net.IP, *ipx.Range, map[string]string, error) {
	return ip, ipRange, data, nil
}
