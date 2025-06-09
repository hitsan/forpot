package ssh

import (
		"fmt"
		"log"
		"golang.org/x/crypto/ssh"
)

func CreateSshConfig(user string, password string) ssh.ClientConfig {
		config := ssh.ClientConfig{
				User: user,
				Auth: []ssh.AuthMethod{
						ssh.Password(password),
				},
				HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}
		return config
}

func Connect(config ssh.ClientConfig, host string, port string) {
		addr := fmt.Sprintf("%s:%s", host, port)
		client, err := ssh.Dial("tcp", addr, &config)
		if err != nil {
				log.Fatal("Failed to dial: ", err)
		}
		defer client.Close()
}
