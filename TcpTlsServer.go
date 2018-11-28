package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"crypto/tls"
)

// Client holds info about connection
type Client struct {
	conn   net.Conn
	Server *server
}

// TCP server
type server struct {
	address                  string // Address to open connection: localhost:3244
	onNewClientCallback      func(c *Client)
	onClientConnectionClosed func(c *Client, err error)
	onNewMessage             func(c *Client, message string)
}

// Read client data from channel
func (c *Client) listen() {
	reader := bufio.NewReader(c.conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			c.conn.Close()
			c.Server.onClientConnectionClosed(c, err)
			return
		}
		c.Server.onNewMessage(c, message)
	}
}

// Send text message to client
func (c *Client) Send(message string) error {
	_, err := c.conn.Write([]byte(message))
	return err
}

// Send bytes to client
func (c *Client) SendBytes(b []byte) error {
	_, err := c.conn.Write(b)
	return err
}

func (c *Client) Conn() net.Conn {
	return c.conn
}

func (c *Client) Close() error {
	return c.conn.Close()
}

// Called right after server starts listening new client
func (s *server) SetNewClientCB(callback func(c *Client)) {
	s.onNewClientCallback = callback
}

// Called right after connection closed
func (s *server) SetClientConnectionClosedCB(callback func(c *Client, err error)) {
	s.onClientConnectionClosed = callback
}

// Called when Client receives new message
func (s *server) SetNewMessageCB(callback func(c *Client, message string)) {
	s.onNewMessage = callback
}

// Start network server
func (s *server) Listen() {

	cer, err := tls.LoadX509KeyPair("server.crt", "server.key")
    if err != nil {
        log.Println(err)
        return
    }

    config := &tls.Config{Certificates: []tls.Certificate{cer}}
    listener, err := tls.Listen("tcp", s.address, config)
    if err != nil {
        log.Println(err)
        return
    }
    defer listener.Close()

	for {
		conn, _ := listener.Accept()
		client := &Client{
			conn:   conn,
			Server: s,
		}
		go client.listen()
		s.onNewClientCallback(client)
	}
}

// Creates new tcp server instance
func New(address string) *server {
	log.Println("Creating server with address", address)
	server := &server{
		address: address,
	}

	server.SetNewClientCB(func(c *Client) {})
	server.SetNewMessageCB(func(c *Client, message string) {})
	server.SetClientConnectionClosedCB(func(c *Client, err error) {})

	return server
}

/* tools function */
func onNewClient(c *Client){
	fmt.Println(c.conn.RemoteAddr().String(), "connected")
}

func onClientNewMessage(c *Client, message string) {
	fmt.Print(c.conn.RemoteAddr().String(), " Msg:", message)
	c.Send(strings.ToUpper(message))
}

func onClientClose(c *Client, err error) {
	fmt.Println(c.conn.RemoteAddr().String(), "closed")
}




func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide a port number!")
		return
	}

	SERVER_ADDRESS := "0.0.0.0" + ":" + arguments[1]

	// s, err := net.ResolveTCPAddr("tcp", SERVER_ADDRESS)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	Server := New(SERVER_ADDRESS)
	Server.SetNewClientCB(onNewClient)
	Server.SetNewMessageCB(onClientNewMessage)
	Server.SetClientConnectionClosedCB(onClientClose)
	Server.Listen()
}
