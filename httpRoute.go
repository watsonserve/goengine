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
	index     map[string]map[string]ActionFunc
	catchers  []*catcher_t
	subRouter []*sub_router_t
}

func InitHttpRoute() *HttpRoute {
	return &HttpRoute{
		index: make(map[string]map[string]ActionFunc),
	}
}

func (hr *HttpRoute) SetByMethod(path, method string, handle ActionFunc) {
	r := hr.index[path]
	if nil != r && len(r) > 0 {
		r[""] = handle
		return
	}
	hr.index[path] = map[string]ActionFunc{"": handle}
}

func (hr *HttpRoute) Set(path string, handle ActionFunc) {
	hr.SetByMethod(path, "", handle)
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
	method := req.Method
	r := hr.index[urlPath]

	handle := r[method]
	if nil == handle {
		handle = r[""]
	}
	if 0 < len(r) && nil == handle {
		res.WriteHeader(http.StatusMethodNotAllowed)
		res.Write([]byte(method + " " + req.URL.Path + " Method Not Allowed"))
		goutils.Errorf("- 405 Method Not Allowed - %s\n", req.URL.Path)
		return false
	}

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
		return true
	}
	return subRouteHandle(res, req)
}
