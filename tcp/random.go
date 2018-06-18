package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"math/rand"
	"time"
)

var messages chan string

func server() {
	fmt.Println("Listening...")

	l, _ := net.Listen("tcp", ":2000")

	messages <- "ping"

	c, _ := l.Accept()
	rand.Seed(time.Now().UTC().UnixNano())

	for i:=0; i<20; i++ {
		r := rand.Intn(10)
		buf := make([]byte, r)
		buf[0] = r
		n, _ := c.Write(buf)
		fmt.Printf("From server to client, bytes = %d\n", n)

		//m, _ := bufio.NewReader(c).ReadString('\n')
		//fmt.Println("Recieved from client: ", m)

		//c.Write([]byte(m))
	}
}

func client() {
	<-messages

	conn, err := net.Dial("tcp", "127.0.0.1:2000")
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(os.Stdin)
	buf := make([]byte, 20)

	for {
		n, _ := conn.Read(buf)
		fmt.Printf("Message from server: %d\n", n)

		fmt.Print("Press enter to continue...")
		reader.ReadString('\n')
	}
}

func main() {
	messages = make(chan string)

	go server()
	go client()

	select {}
}
