package go_piping_server

import (
	"context"
	"errors"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"

	"golang.org/x/net/proxy"
)

type Client string

func (c Client) ListenAndServe(host string) (err error) {
	if host == "" {
		host = ":80"
	}
	ln, err := net.Listen("tcp", host)
	if err != nil {
		return
	}
	return c.Serve(ln)
}

func (c Client) Serve(li net.Listener) error {
	if c == Client("") {
		return errors.New("target isn't specified")
	}
	loc, err := url.ParseRequestURI(string(c))
	if err != nil {
		return err
	}
	r, err := http.NewRequest("PUT", string(c), nil)
	if err != nil {
		return err
	}
	proxy.RegisterDialerType("http", newHTTPProxy)
	httpDialer := proxy.FromEnvironment()
	for {
		cl, err := li.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go func() {
			sv, err := httpDialer.Dial("tcp", loc.Host)
			if err != nil {
				log.Println(err)
				cl.Close()
				return
			}
			req := r.Clone(context.Background())
			req.Write(sv)
			go func() {
				io.Copy(sv, cl)
				cl.Close()
			}()
			go func() {
				io.Copy(cl, sv)
				sv.Close()
			}()
		}()
	}
}
