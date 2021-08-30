package main

import (
	"chat/client"
	"chat/server"
	"flag"
)

const netw = "localhost:5000"

func main() {
	isServer := flag.Bool("s", false, "Start server instead of client")
	netid := flag.String("net", "localhost:5000", "ip:port")
	flag.Parse()
	switch *isServer {
	case true:
		server.Start(*netid)
	case false:
		client.Start(*netid)
	}
}
