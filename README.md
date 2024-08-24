# GoLang Chat Server

This project implements a simple chat server and client in Go, where multiple clients can connect to a server concurrently and exchange messages in real-time. Messages sent by any client are broadcasted to all other connected clients.

## Features

- **Concurrent Client Handling:** The server can handle multiple clients simultaneously.
- **Message Broadcasting:** Messages from one client are broadcasted to all other connected clients.
- **Simple and Lightweight:** The implementation is straightforward and lightweight, making it easy to understand and extend.

## Files

- `server.go`: Contains the server-side code.
- `client.go`: Contains the client-side code.

## Prerequisites

- **Go:** Ensure that Go is installed on your machine. You can download it from [here](https://golang.org/dl/).

## Usage
### server side
1. **Port Configuration**:
   when starting the server you will be prompted for the port you want to use
   press enter to use default, the server listens on port `8080`.

2. **Start the Server:**

   ```bash
   go run server.go
   ```

### client side

1. **port and Url configuration**

   when starting the client you will be prompted for address then port of the server you are trying to connect to

2. **Start the Clients:**

   In separate terminal windows, run:

   ```bash
   go run client.go
   ```

   if the address and port was set correctly the client will then connect to the server as intended
> Note:To simulate a group chat, you can open multiple terminals and run the client code (`client.go`) in each one. All connected clients will be able to send and receive messages in real-time.
   
3. **input your name:**
    
   after the connection is established you will be prompted to input your name,
   *for readability avoid large names 4 to 10 letter names are optimal* 


4. **Send Messages:**

   type messages in any client's terminal. The message will be broadcasted to all connected clients.

> Note: due to the limitation of working on terminal, incoming messages may interrupt and split the text you're typing. However, when you press enter to send the message, it will be sent as a complete and intact message, without losing any content. *this will be fixed later on*

## License

it's juste a little client server on terminal why would it have a Licence ??

## Contributing

Feel free to submit issues or pull requests if you have any improvements or bug fixes.
