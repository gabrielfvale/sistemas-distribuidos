package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"reflect"

	"github.com/gabrielfvale/ti0151-sistemas/app/pkg"
)

var enviroment = pkg.Environment{
	Temperature: 27,
	Luminosity:  0,
	Smoke:       0,
}

func handleConnection(c net.Conn) {
	for {
		netData, err := bufio.NewReader(c).ReadBytes('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
			return
		}

		receivedMessage := pkg.EnvironmentMessage{}
		json.Unmarshal(netData, &receivedMessage)

		fmt.Printf("Message Received: %v \n", receivedMessage)

		if receivedMessage.Type == "write" {
			data := receivedMessage.Data
			reflect.ValueOf(&enviroment).Elem().FieldByName(data.Field).SetInt(int64(data.Value))
		}

		data, err := json.Marshal(enviroment)
		if err != nil {
			log.Printf("error marshal")
		}
		data = append(data, '\n')
		c.Write(data)
	}
	c.Close()
}

func main() {
	fmt.Println("Start enviroment server...")
	ln, _ := net.Listen("tcp", ":8010")
	defer ln.Close()

	for {
		c, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go handleConnection(c)
	}
}
