package go_piping_server

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
)

const errpage = `<!DOCTYPE html><html><head><title>%d %s</title></head><body><center><h1>%d %s</h1></center><hr><center>%s</center></body></html>`

func MakeErrorPage(status int, message, version string) func(http.ResponseWriter, *http.Request) {
	page := fmt.Sprintf(errpage, status, message, status, message, version)
	n := strconv.Itoa(len(page))
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Content-Length", n)
		w.WriteHeader(status)
		w.Write([]byte(page))
	}
}

var BadGateway = MakeErrorPage(502, "Bad Gateway", "nginx/1.20.2")
var BadRequest = MakeErrorPage(400, "Bad Request", "nginx/1.20.2")

type Handler string

func (s Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "PUT", "PUSH":
		{
			hj, ok := w.(http.Hijacker)
			if !ok {
				BadGateway(w, r)
				return
			}
			sv, err := net.Dial("tcp", string(s))
			if err != nil {
				log.Println(err)
				BadGateway(w, r)
				return
			}
			cl, _, err := hj.Hijack()
			if err != nil {
				log.Println(err)
				BadGateway(w, r)
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
			BadRequest(w, r)
		}
	}
}

func HandleFunc(target string) func(http.ResponseWriter, *http.Request) {
	return Handler(target).ServeHTTP
}
