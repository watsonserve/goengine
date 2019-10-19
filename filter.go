package goengine

import (
	"log"
	"net"
	"net/http"
	"os"
	"syscall"
)

type filters_t struct {
	filter         []FilterFunc
}

func (this *filters_t) Use(handle FilterFunc) {
	this.filter = append(this.filter, handle)
}

func (this *filters_t) UseRouter(router *HttpRoute) {
	this.filter = append(this.filter, router.ServeHTTP)
}

func (this *filters_t) Range(res http.ResponseWriter, session *Session, req *http.Request) bool {
	for i := range this.filter {
		if !this.filter[i](res, session, req) {
			return true
		}
	}
	return false
}
