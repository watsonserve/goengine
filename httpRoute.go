package goengine

import (
	"net/http"
	"regexp"
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

func (this *HttpRoute) SetWith(path string, handle *HttpRoute) {
	route := regexp.MustCompile(path)

	this.catcher = append(this.catcher, &catcher_t{
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
		return false
	}

	if this.Range(res, session, req) {
		return false
	}

	handle(res, session, req)
	return true
}
