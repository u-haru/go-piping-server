package go_piping_server

import (
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
)

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
	s.HandleFunc(s.Path, s.handle)
	s.HandleFunc("/", errorp)

	return http.Serve(li, s)
}

const message = `<!DOCTYPE html>
<html>
<head>
<title>Error</title>
<style>
body { width: 96%; margin: 0 auto; }
</style>
</head>
<body>
<center><h1>500 Internal Server Error</h1>
<hr>pipe</center>
</body>
</html>`

func errorp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Content-Length", strconv.Itoa(len(message)))
	w.WriteHeader(500)
	w.Write([]byte(message))
}

func (s *Server) handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "PUT", "PUSH":
		{
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
		}
	default:
		{
			errorp(w, r)
		}
	}
}
