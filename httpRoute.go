package goengine

import (
	"log"
	"net"
	"net/http"
	"os"
	"syscall"
)

type ActionFunc func(http.ResponseWriter, *Session, *http.Request)
type FilterFunc func(http.ResponseWriter, *Session, *http.Request) bool

type HttpRoute struct {
	sessionManager *SessionManager
	index          map[string]ActionFunc
	filter         []FilterFunc
}

func InitHttpRoute(sessionManager *SessionManager) *HttpRoute {
	return &HttpRoute{
		sessionManager: sessionManager,
		index:          make(map[string]ActionFunc),
	}
}

func (this *HttpRoute) Use(handle FilterFunc) {
	this.filter = append(this.filter, handle)
}

func (this *HttpRoute) Set(path string, handle ActionFunc) {
	this.index[path] = handle
}

func (this *HttpRoute) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if nil != err {
		res.WriteHeader(500)
		res.Write([]byte(err.Error()))
	}

	handle := this.index[req.URL.Path]
	if nil == handle {
		res.WriteHeader(404)
		res.Write([]byte(req.URL.Path + " not found"))
		return
	}

	session := this.sessionManager.Get(&res, req)

	header := res.Header()
	header.Set("Cache-Control", "no-cache")

	for i := range this.filter {
		if !this.filter[i](res, session, req) {
			return
		}
	}

	handle(res, session, req)
}

func (this *HttpRoute) ListenUnix(addr string) {
	_ = os.Remove(addr)
	syscall.Umask(0111)
	ln, err := net.Listen("unix", addr)
	if nil != err {
		log.Fatal("failed to start server", err)
		return
	}

	if err = http.Serve(ln, this); nil != err {
		log.Fatal("failed to start server", err)
	}

}
