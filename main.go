package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/jessevdk/go-flags"
)

var opts struct {
	Host string `long:"host" default:"127.0.0.1" description:"IP to bind to"`
	Port uint16 `long:"port" default:"2202" description:"UDP port to bind to"`
}

func newUDPListener(host string, port uint16) (*net.UDPConn, error) {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%v:%d", host, port))

	if err != nil {
		return nil, err
	}

	l, err := net.ListenUDP("udp", addr)

	if err != nil {
		return nil, err
	}

	return l, nil
}

func handleClient(conn *net.UDPConn) {
	b := make([]byte, 1024)
	_, addr, err := conn.ReadFromUDP(b)
	if err != nil {
		log.Printf("Read from UDP failed, err: %v", err)
		return
	}
	log.Printf("Read from client(%v:%v): %v", addr.IP, addr.Port, string(b))
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		if !strings.Contains(err.Error(), "Usage") {
			fmt.Fprintf(os.Stderr, "error: %v\n", err.Error())
			os.Exit(1)
		} else {
			fmt.Fprintf(os.Stderr, "%v\n", err.Error())
			os.Exit(0)
		}
	}

	l, err := newUDPListener(opts.Host, opts.Port)
	if err != nil {
		panic(err)
	}

	log.Printf(">> Starting udpdump, listening at %v:%v...", opts.Host, opts.Port)

	for {
		handleClient(l)
	}
}
