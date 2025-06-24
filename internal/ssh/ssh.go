package ssh

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

type ForwardSession struct {
	uid    Uid
	distIP net.IP
	port   int
	client *ssh.Client
}

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

func fetchUid(client *ssh.Client) (Uid, error) {
	session, err := client.NewSession()
	if err != nil {
		return "", errors.New("Failed to create Session")
	}
	defer session.Close()

	var output bytes.Buffer
	session.Stdout = &output
	if err := session.Run("cat /proc/self/status | grep Uid"); err != nil {
		return "", errors.New("Failed to get UID")
	}
	items := strings.Fields(output.String())
	uid := Uid(items[1])
	return uid, nil
}

func InitSshSession(config ssh.ClientConfig, addr net.TCPAddr) error {
	client, err := ssh.Dial("tcp", addr.String(), &config)
	if err != nil {
		msg := "Failed to dial:" + addr.String()
		return errors.New(msg)
	}
	defer client.Close()

	uid, err := fetchUid(client)
	if err != nil {
		return err
	}

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
		for _, port := range ports {
			forwardPort(client, addr.IP, port)
		}
		time.Sleep(5 * time.Second)
	}
}

func forwardPort(client *ssh.Client, distIP net.IP, port int) error {
	localAddr := fmt.Sprintf("127.0.0.1:%d", port)
	listener, err := net.Listen("tcp", localAddr)
	if err != nil {
		return err
	}

	dist := fmt.Sprintf("%s:%d", distIP.String(), port)
	go func() {
		for {
			localConn, err := listener.Accept()
			if err != nil {
				fmt.Println("err")
				continue
			}
			remoteConn, err := client.Dial("tcp", dist)
			if err != nil {
				fmt.Println("err")
				continue
			}
			go io.Copy(remoteConn, localConn)
			go io.Copy(localConn, remoteConn)
		}
	}()
	return nil
}
