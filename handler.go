package coredns_page

import (
	"context"

	"github.com/coredns/coredns/plugin"
	"github.com/miekg/dns"
)

type Handler struct {
	Next plugin.Handler

	HandlePage func(res *dns.Msg) *dns.Msg
}

func (h *Handler) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	return plugin.NextOrFailure(h.Name(), h.Next, ctx, &PageResponseWriter{ResponseWriter: w, HandlePage: h.HandlePage}, r)
}

func (h *Handler) Name() string {
	return "coredns_page"
}

type PageResponseWriter struct {
	dns.ResponseWriter
	HandlePage func(res *dns.Msg) *dns.Msg
}

func (w *PageResponseWriter) WriteMsg(res *dns.Msg) error {
	if res.Rcode != dns.RcodeSuccess {
		return w.ResponseWriter.WriteMsg(res)
	}

	if res.Question[0].Qtype == dns.TypeAXFR || res.Question[0].Qtype == dns.TypeIXFR {
		return w.ResponseWriter.WriteMsg(res)
	}

	return w.ResponseWriter.WriteMsg(w.HandlePage(res))
}
