package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

// used to read incoming messages and send them on channel ch
func reader(in *bufio.Reader) chan string {
	ch := make(chan string)
	go func() {
		for {
			str, err := in.ReadString(';')
			if err != nil {
				if err.Error() != "EOF" {
					fmt.Println(err)
					break
				}
			}
			str = strings.Trim(str, "\n;")
			ch <- str
		}
	}()
	return ch

}

// used to send messages
// all strings on channel ch are send to the out connection
func sender(out *bufio.Writer) chan string {
	ch := make(chan string)
	go func() {
		for {
			out.WriteString(<-ch + ";")
			out.Flush()
		}
	}()
	return ch
}

// used to read keyboard input
func keyReader() chan string {
	ch := make(chan string)
	go func() {
		keyb := bufio.NewReader(os.Stdin)
		for {
			str, err := keyb.ReadString('\n')
			if err != nil {
				fmt.Println("keyerr")
				break
			}
			str = strings.Trim(str, "\n;")
			ch <- str
		}
	}()
	return ch
}

// handles the main program loop
func mainloop(inCh, outCh chan string) {
	fmt.Println("application started")
	keyCh := keyReader()
	for {
		select {
		case str := <-inCh:
			fmt.Println(str)
		case str := <-keyCh:
			outCh <- str
		}
	}
}

const ip = "localhost"
const port = "5000"

func main() {
	conn, err := net.Dial("tcp", ip+":"+port)
	if err != nil {
		fmt.Println("connection failed")
		return
	}
	in := bufio.NewReader(conn)
	out := bufio.NewWriter(conn)
	inCh := reader(in)
	outCh := sender(out)
	mainloop(inCh, outCh)
	conn.Close()
}
