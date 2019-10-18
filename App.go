package goengine

import (
	"log"
	"net"
	"net/http"
	"os"
	"syscall"
)

type GoEngine struct {
	filters_t
	sessionManager *SessionManager
}

func New(sessionManager *SessionManager) *GoEngine {
	return &GoEngine{
		sessionManager: sessionManager,
	}
}

func (this *GoEngine) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if nil != err {
		res.WriteHeader(500)
		res.Write([]byte(err.Error()))
	}

	session := this.sessionManager.Get(&res, req)

	header := res.Header()
	header.Set("Cache-Control", "no-cache")

	for i := range this.filter {
		if !this.filter[i](res, session, req) {
			return
		}
	}

	res.WriteHeader(404)
	res.Write([]byte(req.URL.Path + " not found"))
}

func (this *GoEngine) ListenUnix(addr string) {
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