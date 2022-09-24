package goengine

import (
	"log"
	"net"
	"net/http"
	"os"
)

// unix "golang.org/x/sys/unix"
type GoEngine struct {
	filters_t
	sessionManager SessionManager
}

func New(sessionManager SessionManager) *GoEngine {
	return &GoEngine{
		sessionManager: sessionManager,
	}
}

func (this *filters_t) UseRouter(router *HttpRoute) {
	this.Use(router.ServeHTTP)
}

func (this *GoEngine) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if nil != err {
		res.WriteHeader(500)
		res.Write([]byte(err.Error()))
	}

	var session *Session
	session = nil
	if nil != this.sessionManager {
		session = this.sessionManager.Get(req)
	}

	header := res.Header()
	header.Set("Cache-Control", "no-cache")
	ctx := req.Context()
	ctx.WithValue(ctx, "session", session)
	req = req.WithContext(ctx)

	if this.Range(res, req) {
		return
	}
	res.WriteHeader(404)
	res.Write([]byte(req.URL.Path + " not found"))
}

func (this *GoEngine) ListenUnix(addr string) {
	_ = os.Remove(addr)
	// unix.Umask(0666)
	ln, err := net.Listen("unix", addr)
	if nil != err {
		log.Fatal("failed to start server", err)
	}

	if err = http.Serve(ln, this); nil != err {
		log.Fatal("failed to start server", err)
	}
}

func (this *GoEngine) ListenTCP(port string) {
	err := http.ListenAndServe(port, this)
	if nil != err {
		log.Fatal("failed to start server", err)
	}
}
