package middleware

import (
	"github.com/kataras/iris/context"
)

type RemoteCaller interface {
	CheckLimit(ip string) bool
}

type IpLimit struct {
	r RemoteCaller
}

func NewIpLimit(s RemoteCaller) context.Handler {
	l := &IpLimit{
		r: s,
	}
	return l.ServeHTTP
}

func (g *IpLimit) ServeHTTP(ctx context.Context) {
	ip := ctx.RemoteAddr()

	ok := g.r.CheckLimit(ip)
	if !ok {
		ctx.JSON("ip limit")
		ctx.StopExecution()
		return
	}
	ctx.Next()
}
