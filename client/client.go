package client

import (
	"bufio"
	"chat/mess"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
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
			m, err := mess.FromJson([]byte(str))
			if err != nil {
				log.Fatal(err)
			}
			ch <- m.String()
		}
	}()
	return ch

}

// used to send messages
// all strings on channel ch are send to the out connection
func sender(out *bufio.Writer) chan mess.Message {
	ch := make(chan mess.Message)
	go func() {
		for {
			js, err := (<-ch).ToJson()
			if err != nil {
				log.Fatal(err)
			}
			out.WriteString(js + ";")
			out.Flush()
		}
	}()
	return ch
}

// displays incoming messages in the textfield
func displayIncoming(inCh chan string, textbox *tview.TextView) {
	for str := range inCh {
		t := textbox.GetText(true)
		t = t + str + "\n\n"
		textbox.SetText(t)
	}
}

func startUi(inCh chan string, outCh chan mess.Message) {
	username := "Anonymous"

	app := tview.NewApplication().EnableMouse(true)
	idbox := tview.NewInputField()
	idbox.SetFieldBackgroundColor(tcell.NewRGBColor(0, 0, 0)).SetBorder(true).SetTitle("User name")
	idbox.SetChangedFunc(func(txt string) {
		username = txt
	})
	box := tview.NewTextView().SetText("")
	box.SetScrollable(true)
	box.SetChangedFunc(func() {
		app.Draw()
	})
	box2 := tview.NewInputField()
	box2.SetFieldBackgroundColor(tcell.NewRGBColor(0, 0, 0)).SetBorder(true)
	box2.SetDoneFunc(func(tcell.Key) {
		outCh <- mess.Message{Sender: username, Data: box2.GetText()}
		box2.SetText("")
	})

	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	flex.AddItem(idbox, 3, 3, false)
	flex.AddItem(box, 0, 3, false)
	flex.AddItem(box2, 3, 1, true)

	go displayIncoming(inCh, box)
	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}

func Start(netid string) {
	conn, err := net.Dial("tcp", netid)
	if err != nil {
		fmt.Println("connection failed")
		return
	}
	in := bufio.NewReader(conn)
	out := bufio.NewWriter(conn)
	inCh := reader(in)
	outCh := sender(out)
	startUi(inCh, outCh)
	conn.Close()
}
