# GoLang Chat Server

This project implements a simple chat server and client in Go, where multiple clients can connect to a server concurrently and exchange messages in real-time. Messages sent by any client are broadcasted to all other connected clients.

## Features

- **Concurrent Client Handling:** The server can handle multiple clients simultaneously.
- **Message Broadcasting:** Messages from one client are broadcasted to all other connected clients.
- **Simple and Lightweight:** The implementation is straightforward and lightweight, making it easy to understand and extend.

## Files

- `server.go`: Contains the server-side code.
- `client.go`: Contains the client-side code.

## Getting Started

### Prerequisites

- **Go:** Ensure that Go is installed on your machine. You can download it from [here](https://golang.org/dl/).

### Running the Server

 **Compile and Run the Server:**

   Open a terminal and navigate to the directory containing `server.go`.

   ```bash
   go run server.go
   ```

   The server will start and listen for incoming client connections on the specified port (default is `8080`).

### Running the Client

1. **Compile and Run the Client:**

   Open another terminal and navigate to the directory containing `client.go`.

   ```bash
   go run client.go
   ```

   The client will prompt you to enter your name and then connect to the server.

### Connecting Multiple Clients

To simulate a group chat, you can open multiple terminals and run the client code (`client.go`) in each one. All connected clients will be able to send and receive messages in real-time.

## Usage
### server side
1. **Port Configuration**:

By default, the server listens on port `8000`. You can modify the port by editing the server code in `server.go` 
```go
func main() {
	port := ":8080"
	http.HandleFunc("/", connection)
	println("handle func set")
```

2. **Start the Server:**

   ```bash
   go run server.go
   ```

### client side

1. **Start the Clients:**

   In separate terminal windows, run:

   ```bash
   go run client.go
   ```

   if the address and port was set correctly the client will then connect to the server as intended
2. **input your name:**
    
after the connection is established you will be prompted to input your name, 
*for readability avoid large names 4 to 10 letter names are optimal* 

3. **Send Messages:**

   type messages in any client's terminal. The message will be broadcasted to all connected clients.

> Note: due to the limitation of working on terminal when a message is received the text being written will be split by it but when pressing enter the message being sent will still keep it's integrity without loss, *this will be fixed later on*

## License

it's juste a little client server on terminal why would it have a Licence ??

## Contributing

Feel free to submit issues or pull requests if you have any improvements or bug fixes.
