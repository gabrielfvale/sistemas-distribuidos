package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	log.Printf("Starting Home Assistant")

	listen, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatalf("Socket listen failed, %s", err)
		os.Exit(1)
	}

	log.Printf("Home Assistant listening")

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatalln(err)
			continue
		}
		go connHandler(conn)
	}
}

func connHandler(conn net.Conn) {
	defer conn.Close()

	var (
		buf = make([]byte, 1024)
		r   = bufio.NewReader(conn)
		w   = bufio.NewWriter(conn)
	)

ILOOP:
	for {
		n, err := r.Read(buf)
		data := string(buf[:n])

		switch err {
		case io.EOF:
			break ILOOP
		case nil:
			log.Println("Message received:", data)
			if isTransportOver(data) {
				break ILOOP
			}
		default:
			log.Fatalf("Receive data failed:%s", err)
			return
		}
	}
	w.Write([]byte("Pong"))
	w.Flush()
	log.Printf("Send: %s", "Pong")
}

func isTransportOver(data string) (over bool) {
	over = strings.HasSuffix(data, "\r\n\r\n")
	return
}
