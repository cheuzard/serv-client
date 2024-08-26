package main

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

var port string
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			fmt.Printf("error in spliting port %v", err)
		}
		cheu, _ := net.LookupHost("cheuzard.ddns.net")
		return ip == cheu[0]
	},
}

type client struct {
	name     string
	conn     *websocket.Conn
	incoming chan string
}

var Messages = make(chan string, 70)
var uList = sync.Map{}

func main() {
	fmt.Printf("please input the port(press enter to use default :8080):\n")
	_, err := fmt.Scanf("%v", &port)
	if err != nil {
		port = "8080"
		fmt.Printf("default port used\n")
	}
	http.HandleFunc("/", connection)
	//println("handle func set")
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	go startServer(ctx)
	select {
	case <-ctx.Done():
		stop()
		fmt.Printf("shutting down")
	}
}
func startServer(ctx context.Context) {
	finalPort := fmt.Sprintf(":%v", port)
	go broadcaster()
	server := &http.Server{Addr: finalPort}

	go func() {
		<-ctx.Done()
		err := server.Shutdown(context.Background())
		if err != nil {
			fmt.Printf("\nerror at shutdown  %v", err)
			return
		}
	}()
	fmt.Printf("listening on port:%v\n", port)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("server dead")
}

func broadcaster() {
	fmt.Printf("broadcaster is up")
	for {
		select {
		case msg := <-Messages:
			fmt.Printf("sending...\n")
			index := strings.Index(msg, "]")
			uList.Range(func(key, value interface{}) bool {
				if index != -1 && value.(*client).name == msg[1:index] {
					return true
				} else {
					value.(*client).incoming <- msg
					return true
				}
			})

		}
	}
}

func connection(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	fmt.Printf("connection started%v\n", origin)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("\n\nconnection failed to start: %v\n", err)
		return
	}
	println("connection upgraded")
	client := &client{
		name:     NameGetter(conn),
		conn:     conn,
		incoming: make(chan string, 10),
	}
	fmt.Printf("client created : %v\n", client.name)
	Messages <- fmt.Sprintf("  >%v just joined", client.name)
	fmt.Printf("  >%v just joined", client.name)
	go HandleUser(client)
}

func HandleUser(client *client) {
	//println("user handler started")
	uList.Store(client.name, client)
	ctx, cancel := context.WithCancel(context.Background())
	defer uList.Delete(client.name)
	defer cancel()
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			return
		}
	}(client.conn)
	go Receiver(client, ctx, cancel)

	for {
		select {
		case msg := <-client.incoming:
			err := client.conn.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				if websocket.IsUnexpectedCloseError(err) {
					uList.Delete(client.name)
				}
				return
			}
		case <-ctx.Done():
			fmt.Printf("terminating user\n")
			return
		}
	}
}

func Receiver(client *client, ctx context.Context, cancel context.CancelFunc) {
	println("receiver started")
	for {
		_, messageByte, err := client.conn.ReadMessage()
		fmt.Printf("received :[%v] %v\n", client.name, string(messageByte))
		if err != nil {
			if websocket.IsUnexpectedCloseError(err) {
				Messages <- fmt.Sprintf("  >%v just left", client.name)
				fmt.Printf("  >%v just left", client.name)
				close(client.incoming)
				uList.Delete(client.name)
				cancel()
			}
			return
		}
		message := string(messageByte)
		if len(message) > 0 && message == "/exit" {
			Messages <- fmt.Sprintf("  >%v just left", client.name)
			fmt.Printf("  >%v just left", client.name)
			close(client.incoming)
			uList.Delete(client.name)
			cancel()
			err := client.conn.Close()
			if err != nil {
				return
			}
			return
		} else {
			select {
			case Messages <- fmt.Sprintf("[%s]: %s", client.name, message):
			case <-ctx.Done():
				return
			}
		}
	}
}

func NameGetter(conn *websocket.Conn) string {
	time.Sleep(time.Millisecond * 100)
	for {
		exists := false
		request := "please input your name (between 3 and 10 characters):\n"
		err := conn.WriteMessage(websocket.TextMessage, []byte(request))
		if err != nil {
			return ""
		}
		_, answer, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("failed getting name")
			return ""
		}
		name := strings.TrimSpace(string(answer))
		if len(name) > 3 && len(name) < 11 {
			uList.Range(func(key, value interface{}) bool {
				if key == name {
					exists = true
					return false
				}
				return true
			})
			if !exists {
				return name
			}
		} else {
			err = conn.WriteMessage(websocket.TextMessage, []byte("\n\t!!!name too small or too long!!!\n"))
			if err != nil {
				return ""
			}
		}
		if exists {
			err = conn.WriteMessage(websocket.TextMessage, []byte("\n\t!!!name already exists!!!\n"))
			if err != nil {
				return ""
			}
		}
	}
}
