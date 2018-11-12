package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func main() {

	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide a server:port string!")
		return
	}

	CONNECT := arguments[1]

	tcpAddr, err := net.ResolveTCPAddr("tcp4", CONNECT)
	if err != nil {
		fmt.Println("ResolveTCPAddr:", err.Error())
		return
	}
	conn, err := net.DialTCP("tcp4", nil, tcpAddr)
	if err != nil {
		fmt.Println("DialTCP:", err.Error())
		return
	}
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")
		text, _ := reader.ReadString('\n')
		_, err := fmt.Fprintf(conn, text)
		if nil != err {
			if err == io.EOF {
				fmt.Print("Server closed")
			} else {
				fmt.Println(err)
			}
			conn.Close()
			break
		}

		message, err := bufio.NewReader(conn).ReadString('\n')
		if nil != err {
			if err == io.EOF {
				fmt.Print("Server closed")
			} else {
				fmt.Println(err)
			}
			conn.Close()
			break
		}
		fmt.Print("->: " + message)
		if strings.TrimSpace(string(text)) == "STOP" {
			fmt.Println("TCP client exiting...")
			conn.Close()
			break
		}
	}
}
