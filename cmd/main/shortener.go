package main

import (
	"context"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net"
	"net/http"
	"time"
	"url-shortener-alt/internal/config"
	"url-shortener-alt/internal/url"
	"url-shortener-alt/internal/url/db"
	"url-shortener-alt/pkg/client/postgresql"
)

func main() {
	log.Println("create router")
	router := httprouter.New()

	cfg := config.GetConfig()

	client, err := postgresql.NewClient(context.TODO(), 3, cfg.Storage)
	if err != nil {
		log.Fatalf("%v", err)
	}
	repository := db.NewRepository(client)
	handler := url.NewHandler(repository)
	handler.Register(router)

	start(router, cfg)
}

func start(router *httprouter.Router, cfg *config.Config) {
	log.Println("start application")

	var listener net.Listener
	var listenError error

	log.Println("listen tcp")
	listener, listenError = net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen.BindIp, cfg.Listen.Port))
	log.Println(fmt.Sprintf("server is listening %s:%s", cfg.Listen.BindIp, cfg.Listen.Port))

	if listenError != nil {
		log.Fatalln(listenError)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatalln(server.Serve(listener))
}
