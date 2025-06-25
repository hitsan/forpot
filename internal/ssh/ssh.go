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
	listener   net.Listener
	remoteAddr string
	isClosed   bool
}

type SessionMG struct {
	ip         string
	client     *ssh.Client
	sessionMap map[int]*ForwardSession
}

func NewForwardSession(localAddr string, remoteAddr string) (*ForwardSession, error) {
	listener, err := net.Listen("tcp", localAddr)
	if err != nil {
		return nil, err
	}
	f := ForwardSession{
		listener:   listener,
		remoteAddr: remoteAddr,
		isClosed:   false,
	}
	return &f, nil
}

func (f *ForwardSession) Close() error {
	if !f.isClosed {
		return errors.New("Already closed")
	}
	f.isClosed = true
	return nil
}

func NewSessionMG(ip *net.IP, client *ssh.Client) *SessionMG {
	return &SessionMG{
		ip:         ip.String(),
		sessionMap: make(map[int]*ForwardSession),
		client:     client,
	}
}

func (s *SessionMG) UpPorts(ports []int) []int {
	var up []int
	for _, port := range ports {
		_, ok := s.sessionMap[port]
		if ok {
			continue
		}
		localAddr := fmt.Sprintf("127.0.0.1:%d", port)
		remoteAddr := fmt.Sprintf("%s:%d", s.ip, port)
		fs, err := NewForwardSession(localAddr, remoteAddr)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fs.forwardPort(s.client)
		s.sessionMap[port] = fs
	}
	return up
}

func (s *SessionMG) DownPorts(ports []int) []int {
	pm := make(map[int]struct{})
	for _, port := range ports {
		pm[port] = struct{}{}
	}
	var down []int
	for port := range s.sessionMap {
		_, ok := pm[port]
		if ok {
			continue
		}
		session := s.sessionMap[port]
		session.Close()
		delete(s.sessionMap, port)
	}
	return down
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

func fetchProcNet(session *ssh.Session) (*string, error) {
	var output bytes.Buffer
	session.Stdout = &output
	if err := session.Run("cat /proc/net/tcp"); err != nil {
		return nil, err
	}
	p := output.String()
	return &p, nil
}

func InitSshSession(config ssh.ClientConfig, addr net.TCPAddr, remoteIP net.IP) error {
	client, err := ssh.Dial("tcp", addr.String(), &config)
	if err != nil {
		return err
	}
	defer client.Close()

	uid, err := fetchUid(client)
	if err != nil {
		return err
	}

	sessionMG := NewSessionMG(&remoteIP, client)
	for {
		session, err := client.NewSession()
		if err != nil {
			log.Fatal(err)
			continue
		}
		defer session.Close()
		pn, err := fetchProcNet(session)
		if err != nil {
			log.Fatal(err)
			continue
		}
		ports := FindForwardablePorts(*pn, uid)
		go sessionMG.DownPorts(ports)
		go sessionMG.UpPorts(ports)
		time.Sleep(5 * time.Second)
	}
}

func (f *ForwardSession) forwardPort(client *ssh.Client) error {
	go func() {
		for {
			localConn, err := f.listener.Accept()
			if err != nil {
				fmt.Println(err)
				continue
			}
			remoteConn, err := client.Dial("tcp", f.remoteAddr)
			if err != nil {
				fmt.Println(err)
				continue
			}
			go io.Copy(remoteConn, localConn)
			go io.Copy(localConn, remoteConn)
		}
	}()
	return nil
}
