package main

import (
	"flag"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	piping "github.com/u-haru/go-piping-tunnel"
)

var (
	Host   string
	Target string
)

func main() {
	flag.StringVar(&Host, "h", "http://0.0.0.0:8001/piping", "Listening Address:Port")
	flag.StringVar(&Target, "t", "127.0.0.1:8888", "Target host")
	flag.Parse()

	loc, err := url.ParseRequestURI(Host)
	if err != nil {
		log.Println(err)
		os.Exit(-1)
	}
	http.Handle(loc.Path, piping.Handler(Target))
	http.HandleFunc("/", piping.BadGateway)

	go func() {
		log.Println("Server running on " + Host + " to " + Target)
		if err := http.ListenAndServe(loc.Host, nil); err != nil {
			log.Println(err)
			os.Exit(-1)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Printf("Signal %s received, shutting down...\n", (<-c).String())
}
