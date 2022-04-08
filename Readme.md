# go-piping-server

```piping server written by golang```

This is experimental.

```cmd/server``` hosts http server and make connection to specified address when client requests specific path.  
example: http://0.0.0.0:8001/piping => localhost:8888  
(maybe I should say listener, not a server)

```cmd/client``` make connection to specific piping server.  
example: http://0.0.0.0:8001/piping => localhost:8888

Usage of server:  
  -h string
        Listening Address:Port (default "http://0.0.0.0:8001/piping")
  -t string
        Target host (default "127.0.0.1:8888")

Usage of client:  
  -h string
        Listening Address:Port (default "0.0.0.0:8000")
  -t string
        Target Path (default "http://127.0.0.1:8001/piping")