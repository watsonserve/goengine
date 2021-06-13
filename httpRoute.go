package goengine

import (
    "github.com/watsonserve/goutils"
    "net/http"
    "regexp"
)

type catcher_t struct {
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
    catcher   []*catcher_t
    subRouter []*sub_router_t
}

func InitHttpRoute() *HttpRoute {
    return &HttpRoute{
        index: make(map[string]ActionFunc),
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
        path:   path,
        length: len(path),
        handle: router.ServeHTTP,
    })
}

// @return go on
func (this *HttpRoute) ServeHTTP(res http.ResponseWriter, session *Session, req *http.Request) bool {
    if this.Range(res, session, req) {
        // 已被拦截，停止流程
        return false
    }
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
        handle(res, session, req)
        return false
    }

    var subRouteHandle FilterFunc = nil
    path_len := len(req.URL.Path)
    // 子路由
    for i := range this.subRouter {
        subRouter := this.subRouter[i]
        if subRouter.length <= path_len && subRouter.path == req.URL.Path[0: subRouter.length] {
            subRouteHandle = subRouter.handle
            break
        }
    }
    // 没有匹配的子路由
    if nil == subRouteHandle {
        goutils.Errorf("- 404 Not Found - %s\n", req.URL.Path)
        return true
    }
    return subRouteHandle(res, session, req)
}
