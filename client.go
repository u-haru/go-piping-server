package go_piping_server

import (
	"context"
	"crypto/tls"
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
	log.Println(loc.Host)
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
			if loc.Scheme == "https" {
				tsv := tls.Client(sv, &tls.Config{InsecureSkipVerify: true})
				if err := tsv.Handshake(); err != nil {
					log.Println(err)
				}
				sv = tsv
			}
			req := r.Clone(context.Background())
			req.Write(sv)
			go func() {
				_, err := io.Copy(sv, cl)
				if err != nil {
					log.Println(err)
				}
				cl.Close()
			}()
			go func() {
				_, err := io.Copy(cl, sv)
				if err != nil {
					log.Println(err)
				}
				sv.Close()
			}()
		}()
	}
}
