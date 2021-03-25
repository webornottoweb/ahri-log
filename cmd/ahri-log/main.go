package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/webornottoweb/ahri-log/config"
	"golang.org/x/crypto/ssh"
)

// Message represents single stream message
type Message struct {
	Host    string
	Content string
}

func main() {
	conn := initConn()
	defer conn.Close()

	forever := make(chan bool)

	stdOut, stdErr := make(chan Message, 255), make(chan Message, 255)

	go runCommand("command", conn, stdOut, stdErr)
	fmt.Println("[LISTENING]")
	go func() {
		for str := range stdOut {
			fmt.Println(str.Content)
		}
	}()

	go func() {
		for str := range stdErr {
			fmt.Println(str.Content)
		}
	}()

	<-forever
}

func runCommand(cmd string, conn *ssh.Client, stdOut chan Message, stdErr chan Message) {
	sess, err := conn.NewSession()
	if err != nil {
		panic(err)
	}
	defer sess.Close()

	sessStdOut, err := sess.StdoutPipe()
	if err != nil {
		panic(err)
	}

	go bindStream(stdOut, &sessStdOut)

	sessStdErr, err := sess.StderrPipe()
	if err != nil {
		panic(err)
	}

	go bindStream(stdErr, &sessStdErr)

	err = sess.Run(cmd)
	if err != nil {
		panic(err)
	}
}

func bindStream(out chan Message, input *io.Reader) {
	reader := bufio.NewReader(*input)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}

		out <- Message{
			Host:    "test",
			Content: line,
		}
	}
}

func initConn() *ssh.Client {
	config := &ssh.ClientConfig{
		User: string(config.Endpoints.User),
		Auth: []ssh.AuthMethod{
			publickey(string(config.Endpoints.Key.Path)),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, err := ssh.Dial("tcp", "host", config)

	if err != nil {
		panic(err)
	}

	return conn
}

func publickey(path string) ssh.AuthMethod {
	key, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	signer, err := ssh.ParsePrivateKeyWithPassphrase(key, []byte(config.Endpoints.Key.Password))
	if err != nil {
		panic(err)
	}

	return ssh.PublicKeys(signer)
}
