package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func server() {
	fmt.Println("Listening...")

	l, _ := net.Listen("tcp", ":2000")
	c, _ := l.Accept()

	for {
		m, _ := bufio.NewReader(c).ReadString('\n')
		c.Write([]byte(m))
	}
}

func client() {
	conn, err := net.Dial("tcp", "127.0.0.1:2000")
	if err != nil {
		panic(err)
	}

	for {
		fmt.Print("Text to send: ")

		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')

		fmt.Fprintf(conn, text + "\n")

		message, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Print("Message from server: " + message)
	}
}

func main() {
	go server()
	go client()
	for {}
}
