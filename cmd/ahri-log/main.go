package main

import (
	"fmt"
	"strings"

	"github.com/webornottoweb/ahri-log/configs"
	"github.com/webornottoweb/ahri-log/internal/connection"
)

func main() {
	var conns []*connection.Connection

	for i := 0; i < len(configs.Endpoints.Servers); i++ {
		conn := connection.New(configs.Endpoints.Servers[i])
		conn.Init()
		defer conn.Close()

		conns = append(conns, conn)
	}

	forever := make(chan bool)

	stdOut, stdErr := make(chan connection.Message, 255), make(chan connection.Message, 255)

	for i := 0; i < len(conns); i++ {
		go conns[i].RunCommand("ls -la", stdOut, stdErr)
	}

	fmt.Println("[LISTENING]")
	go func() {
		for str := range stdOut {
			fmt.Println(strings.Trim(str.Content, "\n"))
		}
	}()

	go func() {
		for str := range stdErr {
			fmt.Println(strings.Trim(str.Content, "\n"))
		}
	}()

	<-forever
}
