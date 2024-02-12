package main

import (
	"fmt"
	"github.com/caioaraujo/go-rate-limiter/internal/infra/web/webserver"
	"net/http"
)

func main() {
	//configs, err := configs.LoadConfig(".")
	//if err != nil {
	//	panic(err)
	//}
	//
	//fmt.Println("some config:", configs.MaxReqPermitidas) // TODO: REMOVER

	wserver := webserver.NewWebServer()
	wserver.AddHandler("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})
	fmt.Println("Starting web server on port 8080")
	go wserver.Start()
	select {}
}