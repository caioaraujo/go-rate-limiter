package main

import (
	"fmt"
	"github.com/caioaraujo/go-rate-limiter/configs"
	"github.com/caioaraujo/go-rate-limiter/internal/infra/cache"
	"github.com/caioaraujo/go-rate-limiter/internal/infra/web/webserver"
	"net/http"
)

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}
	err = cache.Set("MAX_REQ_PERM", configs.MaxReqPermitidas)
	if err != nil {
		panic(err)
	}
	err = cache.Set("MET_BLOQUEIO", configs.MetodoBloqueio)
	if err != nil {
		panic(err)
	}
	err = cache.Set("TEMPO_BLOQUEIO", configs.TempoBloqueioSec)
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
