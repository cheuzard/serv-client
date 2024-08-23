package main

import (
	"bufio"
	"fmt"
	"github.com/gorilla/websocket"
	"os"
	"sync"
	"time"
)

var (
	writeMutex sync.Mutex
)

func main() {
	done := true
	for i := 1; i < 4 || done; i++ {
		done = false
		fmt.Printf("connection try number %v", i)
		conn, _, err := websocket.DefaultDialer.Dial("ws://192.168.1.40:8080", nil)
		if err != nil {
			time.Sleep(time.Second * 2)
			fmt.Println("Error connecting %v", err)
		} else if conn != nil {
			defer func(conn *websocket.Conn) {
				err := conn.Close()
				if err != nil {

				}
			}(conn)
			go Receive(conn)
			println()
			done = true
			for {
				reader := bufio.NewReader(os.Stdin)
				text, _ := reader.ReadString('\n')
				if len(text) < 0 {
					continue
				}
				err := conn.WriteMessage(websocket.TextMessage, []byte(text))
				if err != nil {
					fmt.Println("Error connecting:", err)
					break
				}
				fmt.Printf("sending....\n")
			}
		}

	}
}

func Receive(conn *websocket.Conn) {
	writer := bufio.NewWriter(os.Stdout)
	for {
		_, text, err := conn.ReadMessage()
		//println("a message was received")
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Printf("unexpected close error: %v\n", err)
			} else {
				fmt.Printf("error when receiving msg: %v\n", err)
			}
			break // Exit the loop on error
		}

		_, err = writer.WriteString("\n" + string(text))
		if err != nil {
			return
		}
		writeMutex.Lock()
		err = writer.Flush()
		writeMutex.Unlock()
		if err != nil {
			return
		}
	}
}
