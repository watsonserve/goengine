package goengine_test

import (
    "fmt"
    "net/http"
    "github.com/watsonserve/goengine"
)

func filter (resp http.ResponseWriter, sess *Session, req *http.Request) {
    // do something
    return true
}

func actionFoo (resp http.ResponseWriter, sess *Session, req *http.Request) {
    resp.Write("")
}

func actionBar (resp http.ResponseWriter, sess *Session, req *http.Request) {
    resp.Write("")
}

func ExampleGoengine() {
    router := goengine.InitHttpRoute()
    router.Use(filter)
    router.SetWith("^/foo/.+", actionFoo)
    router.Set("/bar", actionBar)

    engine := goengine.New(nil)
    engine.UseRouter(router)
    engine.ListenTCP(":7070")
}
