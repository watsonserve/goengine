package goengine_test

import (
    // "fmt"
    "net/http"
    "github.com/watsonserve/goengine"
)

func filter (resp http.ResponseWriter, sess *goengine.Session, req *http.Request) bool {
    // do something
    return true
}

func actionFoo (resp http.ResponseWriter, sess *goengine.Session, req *http.Request) {
    resp.Write([]byte(""))
}

func actionBar (resp http.ResponseWriter, sess *goengine.Session, req *http.Request) {
    resp.Write([]byte(""))
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
