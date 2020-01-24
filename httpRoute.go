package goengine

import (
	"net/http"
	"regexp"
)

type catcher_t struct {
	route  *regexp.Regexp
	handle ActionFunc
}

type sub_router_t struct {
	path  string
	length int
	handle FilterFunc
}

type HttpRoute struct {
	filters_t
	index          map[string]ActionFunc
	catcher        []*catcher_t
	subRouter      []*sub_router_t
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
	this.subRouter = append(this.subRouter, &sub_router_t{
		path: path,
		length: len(path),
		handle: router.ServeHTTP,
	})
}

func (this *HttpRoute) ServeHTTP(res http.ResponseWriter, session *Session, req *http.Request) bool {
	handle := this.index[req.URL.Path]
	// 正则路由
	if nil == handle {
		for i := range this.catcher {
			catcher := this.catcher[i]
			if catcher.route.MatchString(req.URL.Path) {
				handle = catcher.handle
				break
			}
		}
	}
	// 发现action
	if nil != handle {
		if !this.Range(res, session, req) {
			handle(res, session, req)
		}
	
		return false
	}

	var subRouteHandle FilterFunc = nil
	// 子路由
	for i := range this.subRouter {
		subRouter := this.subRouter[i]
		if subRouter.path == req.URL.Path[0: subRouter.length] {
			subRouteHandle = subRouter.handle
			break
		}
	}
	// 没有匹配的子路由
	if nil == subRouteHandle {
		return true
	}
	if !this.Range(res, session, req) {
		return subRouteHandle(res, session, req)
	}
	return false
}
