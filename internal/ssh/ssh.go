package ssh

import (
		"fmt"
		"log"
		"bytes"
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

		session, err := client.NewSession()
		if err != nil {
				log.Fatal("Failed to create new session", err)
		}
		defer session.Close()

		var b bytes.Buffer
		session.Stdout = &b
		command := "/usr/bin/cat /proc/net/tcp"

		fmt.Println("middle")

		if err := session.Run(command); err != nil {
				log.Fatal("Faild to run command", err)
		}
		fmt.Println("output: %s", b.String())
}
