package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

var (
	pidChan chan int
	alive   bool
)
var ServUrl = "cheuzard.ddns.net"
var ServPort = "8080"

func main() {
	url := fmt.Sprintf("ws://%v:%v", ServUrl, ServPort)
	alive = true
	pidChan = make(chan int, 1)
	group := sync.WaitGroup{}
	fmt.Printf("trying to connect... ")
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		time.Sleep(time.Second)
	} else if conn != nil {
		go clientStarter(conn, ctx, &group)
		time.Sleep(time.Second)
		pid := <-pidChan
		p, err := os.FindProcess(pid)
		if err != nil {
			return
		}
		select {
		case <-ctx.Done():
			alive = false
			err := conn.Close()
			if err != nil {
				return
			}
			fmt.Printf("\n\nclosing....")
			err = p.Signal(syscall.SIGTERM)
			if err != nil {
				return
			}

		}
	}
}

func clientStarter(conn *websocket.Conn, ctx context.Context, group *sync.WaitGroup) {
	pidChan <- os.Getpid()
	group.Add(1)
	defer group.Done()
	reader := bufio.NewReader(os.Stdin)

	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
		}
	}(conn)
	go Receive(conn)
	fmt.Println()
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Context canceled, stopping clientStarter.")
			return
		default:
			text, _ := reader.ReadString('\n')
			fmt.Printf("\t|----messageSent---->")
			if alive {
				text = strings.TrimSpace(text)
				if len(text) == 0 {
					fmt.Printf("\b")
					continue
				}
				err := conn.WriteMessage(websocket.TextMessage, []byte(text))
				if err != nil {
					fmt.Println("Error connecting")
					break
				}
			}
		}
	}
}

func Receive(conn *websocket.Conn) {
	writer := bufio.NewWriter(os.Stdout)
	for {
		_, text, err := conn.ReadMessage()
		if err != nil {
			return
		}
		_, err = writer.WriteString("\n" + string(text) + "\t|------------------->")
		if err != nil {
			continue
		}
		err = writer.Flush()
		if err != nil {
			return
		}
	}
}
