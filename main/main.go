package main

import (
	"fmt"
	_ "github.com/joho/godotenv/autoload"
	"net"
	"os"
	"time"
)

func main() {
	fmt.Println(os.Getenv("PORT"))
	fmt.Println(os.Getenv("HOST"))
	fmt.Println(fmt.Sprint(os.Getenv("PORT"), ":", os.Getenv("HOST")))
	laddr, err := net.ResolveTCPAddr("tcp", fmt.Sprint(os.Getenv("HOST"), ":", os.Getenv("PORT")))
	if err != nil {
		panic(err)
	}

	l, err := net.ListenTCP("tcp", laddr)
	fmt.Println(l, "log laddr")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println("listening on", os.Getenv("HOST"))
	defer l.Close()

	messages := make(chan string, 5)
	for {
		fmt.Println("started a for loop")
		conn, err := l.Accept()
		if err != nil {
			panic(err)
		}

		go write(conn, messages)
		fmt.Println("connected", conn.RemoteAddr().String(), conn.RemoteAddr().Network())
	}
}

func write(conn net.Conn, messages chan string) {
	conn.Write([]byte(time.Now().String()))

	//messages <- "someone joined\n"
	//reader := bufio.NewReader(conn)
	//
	//outgoing := make(chan string)
	//go func() {
	//	inputReader := bufio.NewReader(os.Stdin)
	//	for {
	//		o, err := inputReader.ReadString('\n')
	//		if err != nil {
	//			fmt.Printf("outgoing error: %v", err)
	//			conn.Close()
	//			return
	//		}
	//		outgoing <- o
	//	}
	//}()
	//
	//incoming := make(chan string)
	//go func() {
	//	for {
	//		i, err := reader.ReadString('\n')
	//		if err != nil {
	//			messages <- "someone disconnected\n"
	//			conn.Close()
	//			return
	//		}
	//		incoming <- i
	//	}
	//}()
	//
	//for {
	//	select {
	//	case msg := <-messages:
	//		conn.Write([]byte(msg))
	//	case in := <-incoming:
	//		messages <- in
	//	}
	//}
}
