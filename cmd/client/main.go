package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	go_piping_server "github.com/u-haru/go-piping-server"
)

var (
	Host   string
	Target string
)

func main() {
	sv := &go_piping_server.Client{}
	flag.StringVar(&sv.Host, "h", "0.0.0.0:8000", "Listening Address:Port")
	flag.StringVar(&sv.Target, "t", "http://127.0.0.1:8001/piping", "Target Path")
	flag.Parse()

	log.Println("Client running on " + sv.Host + " to " + sv.Target)
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
