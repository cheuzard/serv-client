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
			log.Printf("error in spliting port %v", err)
		}
		log.Printf("%v", ip)
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
	fmt.Print("Please input the port (press enter to use default: 8080): \n")
	if _, err := fmt.Scanln(&port); err != nil {
		port = "8080"
		fmt.Println("Default port used")
	}

	http.HandleFunc("/", connection)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go startServer(ctx)

	<-ctx.Done()
	log.Println("Shutting down")
}
func startServer(ctx context.Context) {
	finalPort := fmt.Sprintf(":%v", port)
	go broadcaster()
	server := &http.Server{Addr: finalPort}

	go func() {
		<-ctx.Done()
		if err := server.Shutdown(context.Background()); err != nil {
			log.Printf("\nerror at shutdown  %v", err)
			return
		}
	}()

	log.Printf("listening on port:%v\n", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}

	log.Println("server stopped")
}

func broadcaster() {
	log.Println("broadcaster is up")
	for {
		select {
		case msg := <-Messages:
			log.Println("sending...")
			index := strings.Index(msg, "]")
			uList.Range(func(key, value interface{}) bool {
				if index != -1 && value.(*client).name == msg[1:index] {
					return true
				}
				value.(*client).incoming <- msg
				return true
			})
		}
	}
}

func connection(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	log.Printf("connection started %v\n", origin)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("connection failed to start: %v\n", err)
		return
	}
	log.Println("connection upgraded")
	client := &client{
		name:     NameGetter(conn),
		conn:     conn,
		incoming: make(chan string, 10),
	}
	log.Printf("client created: %v\n", client.name)
	Messages <- fmt.Sprintf("  >%v just joined", client.name)
	log.Printf("  >%v just joined", client.name)
	go HandleUser(client)
}

func HandleUser(client *client) {
	uList.Store(client.name, client)
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		uList.Delete(client.name)
		cancel()
		if err := client.conn.Close(); err != nil {
			return
		}
	}()
	go Receiver(client, ctx, cancel)

	for {
		select {
		case msg := <-client.incoming:
			if err := client.conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
				if websocket.IsUnexpectedCloseError(err) {
					uList.Delete(client.name)
				}
				return
			}
		case <-ctx.Done():
			log.Println("terminating user")
			return
		}
	}
}

func Receiver(client *client, ctx context.Context, cancel context.CancelFunc) {
	defer func() {
		Messages <- fmt.Sprintf("  >%v just left", client.name)
		log.Printf("  >%v just left", client.name)
		close(client.incoming)
		uList.Delete(client.name)
		cancel()
		if err := client.conn.Close(); err != nil {
			return
		}
	}()

	log.Println("receiver started")
	for {
		_, messageByte, err := client.conn.ReadMessage()
		if err != nil {
			return
		}
		message := string(messageByte)
		if len(message) > 0 && message == "/exit" {
			return
		}
		select {
		case Messages <- fmt.Sprintf("[%s]: %s", client.name, message):
		case <-ctx.Done():
			return
		}
	}
}

func NameGetter(conn *websocket.Conn) string {
	time.Sleep(100 * time.Millisecond)
	for {
		exists := false
		request := "please input your name (between 3 and 10 characters):\n"
		if err := conn.WriteMessage(websocket.TextMessage, []byte(request)); err != nil {
			return ""
		}
		_, answer, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("failed getting name")
			return ""
		}
		name := strings.TrimSpace(string(answer))
		if len(name) >= 3 && len(name) <= 10 {
			uList.Range(func(key, _ interface{}) bool {
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
			if err := conn.WriteMessage(websocket.TextMessage, []byte("\n\t!!!name too small or too long!!!\n")); err != nil {
				return ""
			}
		}
		if exists {
			if err := conn.WriteMessage(websocket.TextMessage, []byte("\n\t!!!name already exists!!!\n")); err != nil {
				return ""
			}
		}
	}
}
