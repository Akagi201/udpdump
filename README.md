# udpdump

Simple Golang UDP Server which dumps the incoming client sending message.

## Build
* docker: `docker build -t udpdump .`
* `go build main.go -o udpdump`

## Run
* `--host`: default host is `127.0.0.1`
* `--port`: default port is 2202
* `--file`: If set, the received data will be dumped to the file.
* `--buffer`: default is 10240
