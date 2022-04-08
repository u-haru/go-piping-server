package main

import (
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	go_piping_server "github.com/u-haru/go-piping-server"
)

var (
	Host string
)

func main() {
	sv := &go_piping_server.Server{}
	flag.StringVar(&Host, "h", "http://0.0.0.0:8001/piping", "Listening Address:Port")
	flag.StringVar(&sv.Target, "t", "127.0.0.1:8888", "Target host")
	flag.Parse()

	loc, err := url.ParseRequestURI(Host)
	if err != nil {
		log.Println(err)
		return
	}

	sv.Host = loc.Host
	sv.Path = loc.Path

	log.Println("Server running on " + sv.Host + sv.Path + " to " + sv.Target)
	go func() {
		if err := sv.ListenAndServe(); err != nil {
			log.Println(err)
			os.Exit(-1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	log.Printf("Signal %s received, shutting down...\n", (<-quit).String())
}
