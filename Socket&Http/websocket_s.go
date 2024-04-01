package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)
// WebSocket server
type WSServer struct {
	address                  string // Address to open connection: localhost:3244
	onNewClientCallback      func(c *WSClient)
	onClientConnectionClosed func(c *WSClient, err error)
	onNewMessage             func(c *WSClient, message string)
	upgrader                 *websocket.Upgrader
	clients                  map[*WSClient]bool
	clientMutex              sync.Mutex
}

// WebSocket client
type WSClient struct {
	ws     *websocket.Conn
	server *WSServer
}


// Read client data from channel
func (c *WSClient) listen() {
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			c.server.onClientConnectionClosed(c, err)
			c.ws.Close()
			return
		}
		c.server.onNewMessage(c, string(message))
	}
}

// Send text message to client
func (c *WSClient) Send(message string) error {
	return c.ws.WriteMessage(websocket.TextMessage, []byte(message))
}

// Send bytes to client
func (c *WSClient) SendBytes(b []byte) error {
	return c.ws.WriteMessage(websocket.BinaryMessage, b)
}

func (c *WSClient) Close() error {
	return c.ws.Close()
}

// Called right after server starts listening new client
func (s *WSServer) SetNewClientCB(callback func(c *WSClient)) {
	s.onNewClientCallback = callback
}

// Called right after connection closed
func (s *WSServer) SetClientConnectionClosedCB(callback func(c *WSClient, err error)) {
	s.onClientConnectionClosed = callback
}

// Called when Client receives new message
func (s *WSServer) SetNewMessageCB(callback func(c *WSClient, message string)) {
	s.onNewMessage = callback
}

// Start WebSocket server
func (s *WSServer) Listen() {
	s.clients = make(map[*WSClient]bool)
	s.upgrader = &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws, err := s.upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("upgrade:", err)
			return
		}
		defer ws.Close()

		client := &WSClient{
			ws:     ws,
			server: s,
		}
		s.clientMutex.Lock()
		s.clients[client] = true
		s.clientMutex.Unlock()

		s.onNewClientCallback(client)

		client.listen()

		s.clientMutex.Lock()
		delete(s.clients, client)
		s.clientMutex.Unlock()
	})

	log.Println("WebSocket server started on", s.address)
	err := http.ListenAndServe(s.address, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// Creates new WebSocket server instance
func NewWSServer(address string) *WSServer {
	log.Println("Creating WebSocket server with address", address)
	server := &WSServer{
		address: address,
	}

	server.SetNewClientCB(func(c *WSClient) {})
	server.SetNewMessageCB(func(c *WSClient, message string) {})
	server.SetClientConnectionClosedCB(func(c *WSClient, err error) {})

	return server
}

func onNewWSClient(c *WSClient) {
	fmt.Println(c.ws.RemoteAddr().String(), "connected")
}

func onWSClientNewMessage(c *WSClient, message string) {
	fmt.Println(c.ws.RemoteAddr().String(), "Msg:", message)
	c.Send(strings.ToUpper(message))
}

func onWSClientClose(c *WSClient, err error) {
	fmt.Println(c.ws.RemoteAddr().String(), "closed")
}

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide a port number!")
		return
	}

	SERVER_ADDRESS := "0.0.0.0:" + arguments[1]

	Server := NewWSServer(SERVER_ADDRESS)
	Server.SetNewClientCB(onNewWSClient)
	Server.SetNewMessageCB(onWSClientNewMessage)
	Server.SetClientConnectionClosedCB(onWSClientClose)
	Server.Listen()
}
