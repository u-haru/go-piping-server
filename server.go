package go_piping_server

import (
	"io"
	"log"
	"net"
	"net/http"
)

// Websocketに流れてきたパケットをsocks5として解釈するサーバー
type Server struct {
	Path   string
	Host   string
	Target string

	http.ServeMux
}

func (s *Server) ListenAndServe() (err error) {
	if s.Host == "" {
		s.Host = ":80"
	}
	ln, err := net.Listen("tcp", s.Host)
	if err != nil {
		return
	}
	return s.Serve(ln)
}

func (s *Server) Serve(li net.Listener) (err error) {
	s.HandleFunc(s.Path, func(w http.ResponseWriter, r *http.Request) {
		hj, ok := w.(http.Hijacker)
		if !ok {
			log.Println("Error")
			return
		}
		sv, err := net.Dial("tcp", s.Target)
		if err != nil {
			log.Println(err)
			return
		}
		cl, _, err := hj.Hijack()
		if err != nil {
			log.Println(err)
			return
		}
		go func() {
			io.Copy(cl, sv)
			sv.Close()
		}()
		go func() {
			io.Copy(sv, cl)
			cl.Close()
		}()
	})

	return http.Serve(li, s)
}
