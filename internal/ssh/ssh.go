package ssh

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"

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

func fetchUid(client *ssh.Client) string {
	session, err := client.NewSession()
	if err != nil {
		return ""
	}
	defer session.Close()

	var output bytes.Buffer
	session.Stdout = &output
	if err := session.Run("cat /proc/self/status | grep Uid"); err != nil {
		return ""
	}
	items := strings.Fields(output.String())
	uid := items[1]
	return uid
}

func InitSshConnect(config ssh.ClientConfig, host string, port string) {
	addr := fmt.Sprintf("%s:%s", host, port)
	client, err := ssh.Dial("tcp", addr, &config)
	if err != nil {
		log.Fatal("Failed to dial: ", err)
	}
	defer client.Close()

	uid := fetchUid(client)

	for {
		session, err := client.NewSession()
		if err != nil {
			log.Fatal("Failed to create new session", err)
		}
		defer session.Close()

		var output bytes.Buffer
		session.Stdout = &output
		if err := session.Run("cat /proc/net/tcp"); err != nil {
			log.Println("Command errpr:", err)
		}
		ports := FindForwardablePorts(output.String(), uid)
		fmt.Println(ports)
		time.Sleep(5 * time.Second)
	}
}

// PortForward creates a local port forwarding tunnel
func PortForward(client *ssh.Client, localPort, remoteHost, remotePort string) error {
	localAddr := fmt.Sprintf("localhost:%s", localPort)
	remoteAddr := fmt.Sprintf("%s:%s", remoteHost, remotePort)

	listener, err := net.Listen("tcp", localAddr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %v", localAddr, err)
	}
	defer listener.Close()

	log.Printf("Port forwarding: %s -> %s", localAddr, remoteAddr)

	for {
		localConn, err := listener.Accept()
		if err != nil {
			return fmt.Errorf("failed to accept connection: %v", err)
		}

		go func() {
			defer localConn.Close()

			remoteConn, err := client.Dial("tcp", remoteAddr)
			if err != nil {
				log.Printf("failed to dial remote: %v", err)
				return
			}
			defer remoteConn.Close()

			// Copy data bidirectionally
			go io.Copy(localConn, remoteConn)
			io.Copy(remoteConn, localConn)
		}()
	}
}

// ConnectWithPortForward establishes SSH connection and sets up port forwarding
func ConnectWithPortForward(config ssh.ClientConfig, host, port, localPort, remoteHost, remotePort string) error {
	addr := fmt.Sprintf("%s:%s", host, port)
	client, err := ssh.Dial("tcp", addr, &config)
	if err != nil {
		return fmt.Errorf("failed to dial: %v", err)
	}
	defer client.Close()

	return PortForward(client, localPort, remoteHost, remotePort)
}
