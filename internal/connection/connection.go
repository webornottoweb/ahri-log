package connection

import (
	"bufio"
	"fmt"
	"io"

	"github.com/webornottoweb/ahri-log/configs"
	pubkey "github.com/webornottoweb/ahri-log/internal/pkg"
	"golang.org/x/crypto/ssh"
)

// Connection represents one server to be connected to
type Connection struct {
	Server configs.EndpointServer
	conn   *ssh.Client
}

// New returns created Connection instance
func New(server configs.EndpointServer) *Connection {
	return &Connection{Server: server}
}

// Init connection
func (c *Connection) Init() {
	config := &ssh.ClientConfig{
		User: string(configs.Auth.User),
		Auth: []ssh.AuthMethod{
			pubkey.GetAuth(string(configs.Auth.Key.Path)),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, err := ssh.Dial("tcp", string(c.Server.Host)+":"+fmt.Sprint(c.Server.Port), config)
	if err != nil {
		panic(err)
	}

	c.conn = conn
}

// Close connection
func (c *Connection) Close() {
	if c.conn == nil {
		return
	}

	c.conn.Close()
}

func (c *Connection) RunCommand(cmd string, stdOut chan Message, stdErr chan Message) {
	sess, err := c.conn.NewSession()
	if err != nil {
		panic(err)
	}
	defer sess.Close()

	sessStdOut, err := sess.StdoutPipe()
	if err != nil {
		panic(err)
	}

	go c.bindStream(stdOut, &sessStdOut)

	sessStdErr, err := sess.StderrPipe()
	if err != nil {
		panic(err)
	}

	go c.bindStream(stdErr, &sessStdErr)

	sess.Run(cmd)
}

func (c *Connection) bindStream(out chan Message, input *io.Reader) {
	reader := bufio.NewReader(*input)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

		out <- Message{
			Host:    string(c.Server.Host),
			Content: line,
		}
	}
}
