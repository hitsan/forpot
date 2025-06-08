package ssh

import (
		"log"
		"golang.org/x/crypto/ssh"
)

func CreateSshConfig(user string) ssh.ClientConfig {
		config := ssh.ClientConfig{
				User: user,
				Auth: []ssh.AuthMethod{
						ssh.Password(""),
				},
				HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}
		return config
}

func connect(config ssh.ClientConfig, host string, port string) {
		client, err := ssh.Dial("tcp", host, &config)
		if err != nil {
				log.Fatal("Failed to dial: ", err)
		}
		defer client.Close()

//		destination := host + ":" + port
//		listener, err := net.Listen("tcp", destination)
//		if err != nil {
//				log.Fatalf("Failed to establish listener")
//		}
//		defer listener.Close()
//
//		for {
//				localConn, err := listener.Accept()
//				if err != nil {
//						log.Printf("failed to accept connection")
//				}
//
//				source = "localhost:" + port
//				remoteConn, err := client.Dial("tcp", "127.0.0.1:80")
//				if err != nil {
//						log.Printf("failed to accept connection")
//				}
//		}
}
