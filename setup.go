package coredns_page

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/miekg/dns"
)

func init() {
	plugin.Register("coredns_page", setup)
}

func setup(c *caddy.Controller) error {
	handler, err := parse(c)
	if err != nil {
		return err
	}
	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		return &Handler{Next: next, HandlePage: handler}
	})
	return nil
}

func parse(c *caddy.Controller) (func(msg *dns.Msg) *dns.Msg, error) {
	for c.Next() {
		args := c.RemainingArgs()
		switch len(args) {
		case 0:
			return pageHandler(10), nil
		case 1:
			pageSize, err := strconv.Atoi(args[0])
			if err != nil {
				return nil, fmt.Errorf("invalid page size args: %v", args[0])
			}
			return pageHandler(pageSize), nil
		default:
			return nil, fmt.Errorf("too many args: %v", args)
		}
	}
	return nil, c.ArgErr()
}

func pageHandler(pageSize int) func(msg *dns.Msg) *dns.Msg {
	var currentPage int
	var mux sync.Mutex
	return func(msg *dns.Msg) *dns.Msg {
		mux.Lock()
		defer mux.Unlock()

		start := currentPage * pageSize
		if start >= len(msg.Answer) {
			start = 0
			currentPage = 0
		}
		end := min(start+pageSize, len(msg.Answer))

		msg.Answer = msg.Answer[start:end]
		currentPage++
		return msg
	}
}
