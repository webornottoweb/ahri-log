package pubkey

import (
	"fmt"
	"io/ioutil"

	"github.com/webornottoweb/ahri-log/configs"
	"golang.org/x/crypto/ssh"
)

// GetAuth metod for provided ssh key path
func GetAuth(path string) ssh.AuthMethod {
	key, err := ioutil.ReadFile(path)
	if err != nil {
		panic(fmt.Sprintf("[%s] could not connect: [%s]", path, err))
	}

	signer, err := ssh.ParsePrivateKeyWithPassphrase(key, []byte(configs.Auth.Key.Password))
	if err != nil {
		panic(fmt.Sprintf("[%s] can't parse private key: [%s]", path, err))
	}

	return ssh.PublicKeys(signer)
}
