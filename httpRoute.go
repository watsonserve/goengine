package goengine

import (
	"net/http"
	"regexp"
	"fmt"
)

type catcher_t struct {
	route  *regexp.Regexp
	handle ActionFunc
}

type HttpRoute struct {
	filters_t
	index          map[string]ActionFunc
	catcher        []*catcher_t
}

func InitHttpRoute() *HttpRoute {
	return &HttpRoute{
		index:          make(map[string]ActionFunc),
	}
}

func (this *HttpRoute) Set(path string, handle ActionFunc) {
	this.index[path] = handle
}

func (this *HttpRoute) SetWith(path string, handle ActionFunc) {
	this.SetRegexp(regexp.MustCompile(path), handle)
}

func (this *HttpRoute) SetRegexp(route *regexp.Regexp, handle ActionFunc) {
	this.catcher = append(this.catcher, &catcher_t{
		route: route,
		handle: handle,
	})
}

func (this *HttpRoute) UseRouter(path string, router *HttpRoute) {
	route := regexp.MustCompile(path)
	this.Use(func(res http.ResponseWriter, session *Session, req *http.Request) bool {
		if route.MatchString(req.URL.Path) {
			return router.ServeHTTP(res, session, req)
		}
		return false
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
		return true
	}

	if !this.Range(res, session, req) {
		handle(res, session, req)
	}

	return false
}
