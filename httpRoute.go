package goengine

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/watsonserve/goutils"
)

const R_REGEXP = 1
const R_PREFIX = 2

type catcher_t struct {
	rType  int
	fix    string
	route  *regexp.Regexp
	handle ActionFunc
}

type sub_router_t struct {
	path   string
	length int
	handle FilterFunc
}

type HttpRoute struct {
	filters_t
	index     map[string]ActionFunc
	catchers  []*catcher_t
	subRouter []*sub_router_t
}

func InitHttpRoute() *HttpRoute {
	return &HttpRoute{
		index: make(map[string]ActionFunc),
	}
}

func (hr *HttpRoute) Set(path string, handle ActionFunc) {
	hr.index[path] = handle
}

func (hr *HttpRoute) SetWith(path string, handle ActionFunc) {
	hr.SetRegexp(regexp.MustCompile(path), handle)
}

func (hr *HttpRoute) SetRegexp(route *regexp.Regexp, handle ActionFunc) {
	hr.catchers = append(hr.catchers, &catcher_t{
		rType:  R_REGEXP,
		route:  route,
		handle: handle,
	})
}

func (hr *HttpRoute) StartWith(path string, handle ActionFunc) {
	hr.catchers = append(hr.catchers, &catcher_t{
		rType:  R_PREFIX,
		fix:    path,
		handle: handle,
	})
}

func (hr *HttpRoute) UseRouter(path string, router *HttpRoute) {
	hr.subRouter = append(hr.subRouter, &sub_router_t{
		path:   path,
		length: len(path),
		handle: router.ServeHTTP,
	})
}

// @return go on
func (hr *HttpRoute) ServeHTTP(res http.ResponseWriter, req *http.Request) bool {
	if hr.Range(res, req) {
		// 已被拦截，停止流程
		return false
	}
	urlPath := req.URL.Path
	handle := hr.index[urlPath]
	// 正则路由
	if nil == handle {
		for _, catcher := range hr.catchers {
			if R_PREFIX == catcher.rType && strings.HasPrefix(urlPath, catcher.fix) ||
				R_REGEXP == catcher.rType && catcher.route.MatchString(urlPath) {
				handle = catcher.handle
				break
			}
		}
	}
	// 发现action
	if nil != handle {
		handle(res, req)
		return false
	}

	var subRouteHandle FilterFunc = nil
	path_len := len(req.URL.Path)
	// 子路由
	for _, subRouter := range hr.subRouter {
		if subRouter.length <= path_len && subRouter.path == req.URL.Path[0:subRouter.length] {
			subRouteHandle = subRouter.handle
			break
		}
	}
	// 没有匹配的子路由
	if nil == subRouteHandle {
		goutils.Errorf("- 404 Not Found - %s\n", req.URL.Path)
		return true
	}
	return subRouteHandle(res, req)
}
