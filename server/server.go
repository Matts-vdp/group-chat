package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

// handles a incomming connection
func handleConnection(conn net.Conn, inch, outch chan string) {
	out := bufio.NewWriter(conn)
	in := bufio.NewReader(conn)
	go receiver(in, inch)
	go sender(out, outch)
}

// processes any incoming strings from in and sends them on channel ch
func receiver(in *bufio.Reader, ch chan string) {
	for {
		str, err := in.ReadString(';')
		if err != nil {
			return
		}
		str = strings.Trim(str, "\n")
		ch <- str
	}
}

// reads any incoming strings on channel ch and sends them on the out channel
func sender(out *bufio.Writer, ch chan string) {
	for {
		_, err := out.WriteString(<-ch)
		if err != nil {
			return
		}
		err = out.Flush()
		if err != nil {
			return
		}
	}
}

// broadcast all incoming strings on the ch channel to all registered outputs
// the newOut channel can be used to add channels to the registered outputs
func broadcast(newOut chan chan string) chan string {
	ch := make(chan string)
	outputs := make([]chan string, 0)
	go func() {
		for {
			select {
			case new := <-newOut:
				outputs = append(outputs, new)
			case str := <-ch:
				for _, out := range outputs {
					out <- str
				}
			}
		}
	}()
	return ch
}

const ip = "localhost"
const port = "5000"

func main() {
	ln, err := net.Listen("tcp", ip+":"+port)
	if err != nil {
		fmt.Println("Connection setup failed")
		return
	}
	addBroadcastCh := make(chan chan string)
	broadcastCh := broadcast(addBroadcastCh)
	i := 0
	fmt.Println("started")
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("client failed to connect")
			continue
		}
		fmt.Println("client connected")
		ch := make(chan string)
		addBroadcastCh <- ch
		handleConnection(conn, broadcastCh, ch)
		i++
	}
}
