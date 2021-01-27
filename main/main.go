package main

import (
	"bufio"
	"fmt"
	_ "github.com/joho/godotenv/autoload"
	"net"
	"os"
)

func main() {
	laddr, err := net.ResolveTCPAddr("tcp", os.Getenv("HOST"))
	if err != nil {
		panic(err)
	}

	l, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		panic(err)
	}
	defer l.Close()

	messages := make(chan string, 5)
	for {
		conn, err := l.Accept()
		if err != nil {
			panic(err)
		}

		go write(conn, messages)
		fmt.Println("connected", conn.RemoteAddr().String(), conn.RemoteAddr().Network())
	}
}

func write(conn net.Conn, messages chan string) {
	messages <- "someone joined\n"
	reader := bufio.NewReader(conn)

	outgoing := make(chan string)
	go func() {
		inputReader := bufio.NewReader(os.Stdin)
		for {
			o, err := inputReader.ReadString('\n')
			if err != nil {
				fmt.Printf("outgoing error: %v", err)
				return
			}
			outgoing <- o
		}
	}()

	incoming := make(chan string)
	go func() {
		for {
			i, err := reader.ReadString('\n')
			if err != nil {
				messages <- "someone disconnected\n"
				conn.Close()
				return
			}
			incoming <- i
		}
	}()

	for {
		select {
		case msg := <-messages:
			conn.Write([]byte(msg))
		case in := <-incoming:
			messages <- in
		}
	}
}
