package main

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strings"
	"sync"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type client struct {
	name     string
	conn     *websocket.Conn
	incoming chan string
}

var Messages = make(chan string, 70)
var uList = sync.Map{}

func main() {
	port := ":8080" // Load from config or env variable
	http.HandleFunc("/", connection)
	println("handle func set")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go broadcaster()
	server := &http.Server{Addr: port}

	go func() {
		<-ctx.Done()
		err := server.Shutdown(context.Background())
		if err != nil {
			fmt.Printf("error at shutdown  %v", err)
		}
	}()

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func broadcaster() {
	for {
		select {
		case msg := <-Messages:
			uList.Range(func(key, value interface{}) bool {
				value.(*client).incoming <- msg
				return true
			})

		}
	}
}

func connection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("\n\nconnection failed to start: %v\n", err)
	}
	println("connection upgraded")
	client := &client{
		name:     NameGetter(conn),
		conn:     conn,
		incoming: make(chan string, 10),
	}
	println("client created")
	go HandleUser(client)
}

func HandleUser(client *client) {
	println("user handler started")
	uList.Store(client.name, client)
	ctx, cancel := context.WithCancel(context.Background())
	defer uList.Delete(client.name)
	defer cancel()
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Printf("\n\nconnection failed to close: %v", err)
		}
	}(client.conn)
	go Receiver(client, ctx)

	for {
		select {
		case msg := <-client.incoming:
			err := client.conn.WriteMessage(websocket.TextMessage, []byte(msg))
			fmt.Printf("message sent to users")
			if err != nil {
				fmt.Printf("\n\nerror when writing message to websocket: %v", err)
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

func Receiver(client *client, ctx context.Context) {
	println("receiver started")
	for {
		_, messageByte, err := client.conn.ReadMessage()
		fmt.Printf("a message was received: %v", string(messageByte))
		if err != nil {
			return
		}
		message := string(messageByte)
		if len(message) > 0 && message == "/exit" {
			uList.Delete(client.name)
			err := client.conn.Close()
			if err != nil {
				return
			}
		} else {
			select {
			case Messages <- fmt.Sprintf("[%s] %s", client.name, message):
			case <-ctx.Done():
				return
			}
		}
	}
}

func NameGetter(conn *websocket.Conn) string {
	request := "please only send your name in the next message:\n"
	err := conn.WriteMessage(websocket.TextMessage, []byte(request))
	if err != nil {
		return ""
	}
	_, answer, err := conn.ReadMessage()
	if err != nil {
	}
	name := strings.TrimSpace(string(answer))
	return name
}
