package main

import (
	"fmt"
	"net/http"

	"github.com/caioaraujo/go-rate-limiter/configs"
	"github.com/caioaraujo/go-rate-limiter/internal/infra/cache"
	"github.com/caioaraujo/go-rate-limiter/internal/infra/web/webserver"
)

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}
	client := cache.Connect()
	err = cache.Set(client, "MAX_REQ_PERM", configs.MaxReqPermitidas, 0)
	if err != nil {
		panic(err)
	}
	err = cache.Set(client, "MET_BLOQUEIO", configs.MetodoBloqueio, 0)
	if err != nil {
		panic(err)
	}
	err = cache.Set(client, "TEMPO_BLOQUEIO", configs.TempoBloqueioSec, 0)
	if err != nil {
		panic(err)
	}

	wserver := webserver.NewWebServer()
	wserver.AddHandler("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})
	fmt.Println("Starting web server on port 8080")
	go wserver.Start()
	select {}
}
