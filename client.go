package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var (
	pidChan chan int
	alive   bool
)

func main() {
	fmt.Println("Please input the server address:")
	var ServUrl string
	if _, err := fmt.Scanln(&ServUrl); err != nil {
		log.Print("error getting Url")
	}

	fmt.Println("Please input the port (press enter for default: 8080):")
	var ServPort string
	if _, err := fmt.Scanln(&ServPort); err != nil {
		ServPort = "8080"
		fmt.Println("Default port selected")
	}

	ServUrl = strings.TrimSpace(ServUrl)
	ServPort = strings.TrimSpace(ServPort)
	url := fmt.Sprintf("ws://%v:%v", ServUrl, ServPort)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			fmt.Println("Trying to connect...")
			conn, _, err := websocket.DefaultDialer.Dial(url, nil)
			if err != nil {
				fmt.Println(err)
				time.Sleep(time.Second)
			} else {
				go clientStarter(conn, ctx, stop)
				pid := <-pidChan
				p, err := os.FindProcess(pid)
				if err != nil {
					return
				}

				select {
				case <-ctx.Done():
					alive = false
					if err := conn.Close(); err != nil {
						return
					}
					fmt.Println("\nClosing...")
					if err := p.Signal(syscall.SIGTERM); err != nil {
						return
					}
				}
			}
		}
	}
}

func clientStarter(conn *websocket.Conn, ctx context.Context, stop context.CancelFunc) {
	pidChan <- os.Getpid()

	reader := bufio.NewReader(os.Stdin)

	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
		}
	}(conn)
	go Receive(conn, stop)
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

func Receive(conn *websocket.Conn, stop context.CancelFunc) {
	writer := bufio.NewWriter(os.Stdout)
	for {
		_, text, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("connection terminated shutting down")
			stop()
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
