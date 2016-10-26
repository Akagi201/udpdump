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
	Host   string `long:"host" default:"0.0.0.0" description:"IP to bind to"`
	Port   uint16 `long:"port" default:"2202" description:"UDP port to bind to"`
	File   string `long:"file" default:"" description:"dump received data to a dump file"`
	Buffer int    `long:"buffer" default:"10240" description:"max buffer size for the socket io"`
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
	b := make([]byte, opts.Buffer)
	n, addr, err := conn.ReadFromUDP(b)
	if err != nil {
		log.Printf("Read from UDP failed, err: %v", err)
		return
	}
	log.Printf("Read from client(%v:%v), len: %v, [%v]", addr.IP, addr.Port, n, string(b[:n]))

	if len(opts.File) != 0 {
		f, err := os.OpenFile(opts.File, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			log.Printf("Open file failed, err: %v", err)
			return
		}
		defer f.Close()
		if _, err = f.Write(b[:n]); err != nil {
			log.Printf("Write file failed, err: %v", err)
			return
		}
	}

	conn.WriteToUDP(b[:n], addr)
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		if !strings.Contains(err.Error(), "Usage") {
			log.Printf("error: %v\n", err.Error())
			os.Exit(1)
		} else {
			log.Printf("%v\n", err.Error())
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
