package main

import (
	"context"
	"log"
	"net"
	"time"
)

func main() {
	log.Println("hello world")
	r := net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: 10 * time.Second,
			}
			return d.DialContext(ctx, "udp", "223.5.5.5:53")
		},
	}
	ips, _ := r.LookupHost(context.Background(), "www.gzhu.app")
	log.Println(ips)
}
