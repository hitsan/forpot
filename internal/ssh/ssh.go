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
		fmt.Println(ports)
		time.Sleep(5 * time.Second)
	}
}
