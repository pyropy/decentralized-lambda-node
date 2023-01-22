package main

import (
	"github.com/pyropy/decentralised-lambda/runtime"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hi"))
}

func main() {
	rt := runtime.NewLambdaRuntime()
	rt.RegisterHandlerFunc(handler)
	rt.Run()
}
