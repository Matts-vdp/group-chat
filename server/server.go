package server

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

// handles a incomming connection
func handleConnection(conn net.Conn, inch, outch chan string, id int, closeCh chan int) {
	out := bufio.NewWriter(conn)
	in := bufio.NewReader(conn)
	go receiver(in, inch, id, closeCh)
	go sender(out, outch)
}

// processes any incoming strings from in and sends them on channel ch
func receiver(in *bufio.Reader, ch chan string, id int, closeCh chan int) {
	for {
		str, err := in.ReadString(';')
		if err != nil {
			closeCh <- id
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
func broadcast(newOut chan chan string) (chan string, chan int) {
	ch := make(chan string)
	closeCh := make(chan int)
	outputs := make(map[int]chan string)
	go func() {
		i := 0
		for {
			select {
			case new := <-newOut:
				outputs[i] = new
				i++
			case id := <-closeCh:
				fmt.Println("client " + fmt.Sprintf("%d", id) + " disconnected")
				delete(outputs, id)
			case str := <-ch:
				for _, out := range outputs {
					out <- str
				}
			}
		}
	}()
	return ch, closeCh
}

func Start(netid string) {
	ln, err := net.Listen("tcp", netid)
	if err != nil {
		fmt.Println("Connection setup failed")
		return
	}
	addBroadcastCh := make(chan chan string)
	broadcastCh, closeCh := broadcast(addBroadcastCh)
	id := 0
	fmt.Println("started")
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("client failed to connect")
			continue
		}
		fmt.Println("client " + fmt.Sprintf("%d", id) + " connected")
		outch := make(chan string)
		addBroadcastCh <- outch
		handleConnection(conn, broadcastCh, outch, id, closeCh)
		id++
	}
}
