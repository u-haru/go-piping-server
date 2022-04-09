package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	piping "github.com/u-haru/go-piping-server"
)

var Host string

func main() {
	cl := &piping.Client{}
	flag.StringVar(&Host, "h", "0.0.0.0:8000", "Listening Address:Port")
	flag.StringVar(&cl.Target, "t", "http://127.0.0.1:8001/piping", "Target Path")
	flag.Parse()

	go func() {
		log.Println("Client running on " + Host + " to " + cl.Target)
		if err := cl.ListenAndServe(Host); err != nil {
			log.Println(err)
			os.Exit(-1)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Printf("Signal %s received, shutting down...\n", (<-c).String())
}
