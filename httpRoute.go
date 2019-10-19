package goengine

import (
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"syscall"
)

type catcher_t {
	route  Regexp
	handle ActionFunc
}

type HttpRoute struct {
	filters_t
	index          map[string]ActionFunc
	catcher        []*catcher_t
}

func InitHttpRoute(sessionManager *SessionManager) *HttpRoute {
	return &HttpRoute{
		sessionManager: sessionManager,
		index:          make(map[string]ActionFunc),
	}
}

func (this *HttpRoute) Set(path string, handle ActionFunc) {
	this.index[path] = handle
}

func (this *HttpRoute) SetWith(path string, handle ActionFunc) {
	route := regexp.MustCompile(path)

	append(this.catcher, &catcher_t{
		route: route,
		handle: handle,
	})
}

func (this *HttpRoute) ServeHTTP(res http.ResponseWriter, session *Session, req *http.Request) bool {
	handle := this.index[req.URL.Path]
	if nil == handle {
		for i := range this.catcher {
			catcher := this.catcher[i]
			if catcher.route.MatchString(req.URL.Path) {
				handle = catcher.handle
				break
			}
		}
	}
	if nil == handle {
		res.WriteHeader(404)
		res.Write([]byte(req.URL.Path + " not found"))
		return
	}

	if this.Range(res, session, req) {
		return
	}

	handle(res, session, req)
}
