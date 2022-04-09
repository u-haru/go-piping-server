package go_piping_server

import (
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
)

type Server struct {
	Target string

	http.ServeMux
}

const Error502page = `<!DOCTYPE html>
<html><head><title>502 Bad Gateway</title></head>
<body>
<center><h1>502 Bad Gateway</h1></center>
<hr><center>nginx/1.20.2</center>
</body></html>`

func BadGateway(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Content-Length", strconv.Itoa(len(Error502page)))
	w.WriteHeader(500)
	w.Write([]byte(Error502page))
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
			BadGateway(w, r)
		}
	}
}

func Handler(target string) *Server {
	return &Server{
		Target: target,
	}
}

func HandleFunc(target string) func(http.ResponseWriter, *http.Request) {
	return Handler(target).ServeHTTP
}
