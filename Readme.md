# go-piping-tunnel

```piping server written by golang```

This is experimental.

```cmd/server``` hosts http server and make connection to specified address when client requests specific path.  
Example: http://0.0.0.0:8001/piping => localhost:8888  

```cmd/client``` make connection to specific piping server.  
Example: localhost:8000 => http://0.0.0.0:8001/piping

If you want to use proxy,set ALL_PROXY envroiment.  
Example:```ALL_PROXY=http://localhost:8080 ./client```

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

## Usage
localhost:8000 => localhost:22
``` bash
// Server side
$ server -h http://localhost:8001/piping -t localhost:22
```
``` bash
// Client side
$ client -h localhost:8000 -t http://localhost:8001/piping
```
Then you can ```ssh localhost -p 8000```
