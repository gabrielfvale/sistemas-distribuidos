package pkg

import (
	"bufio"
	"encoding/json"
	"log"
	"net"
)

func ConnectToEnviroment() net.Conn {
	conn, err := net.Dial("tcp", "127.0.0.1:8010")
	if err != nil {
		log.Printf("CONN ERR %v", err)
	}
	return conn
}

func ReadEnviromentData(conn net.Conn) Environment {
	sent := EnvironmentMessage{Type: "read"}
	sentData, _ := json.Marshal(sent)
	sentData = append(sentData, '\n')
	conn.Write(sentData)

	message, _ := bufio.NewReader(conn).ReadBytes('\n')
	received := Environment{}
	err := json.Unmarshal(message, &received)
	if err != nil {
		log.Printf("Error")
	}
	return received
}

func WriteEnviromentData(conn net.Conn, field string, value int32) {
	sent := EnvironmentMessage{
		Type: "write",
		Data: EnviromentMessageData{
			Field: field,
			Value: value,
		},
	}

	sentData, err := json.Marshal(sent)
	if err != nil {
		log.Fatalf("Error marshalling data")
	}
	sentData = append(sentData, '\n')
	conn.Write(sentData)
}
