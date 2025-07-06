package ssh

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
)

type ForwardSession struct {
	listener   net.Listener
	remoteAddr string
	once       sync.Once
	done       chan struct{}
}

type SessionMG struct {
	remoteHost string
	client     *ssh.Client
	sessionMap map[int]*ForwardSession
}

type SessionFunc func() error

func createSession(fn SessionFunc, sec time.Duration, done chan struct{}) {
	go func() {
		for {
			select {
			case <-done:
				fmt.Println("close")
				return
			default:
				err := fn()
				if err != nil {
					log.Printf("Failed")
				}
			}
			time.Sleep(sec * time.Second)
		}
	}()
}

func NewForwardSession(localAddr string, remoteAddr string) (*ForwardSession, error) {
	listener, err := net.Listen("tcp", localAddr)
	if err != nil {
		return nil, err
	}
	f := ForwardSession{
		listener:   listener,
		remoteAddr: remoteAddr,
		done:       make(chan struct{}),
	}
	return &f, nil
}

func (f *ForwardSession) Close() {
	f.once.Do(func() {
		close(f.done)
		f.listener.Close()
	})
}

func NewSessionMG(remoteHost string, client *ssh.Client) *SessionMG {
	return &SessionMG{
		remoteHost: remoteHost,
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
		remoteAddr := fmt.Sprintf("%s:%d", s.remoteHost, port)
		fs, err := NewForwardSession(localAddr, remoteAddr)
		if err != nil {
			fmt.Println(err)
			continue
		}
		go fs.forwardPort(s.client)
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

func createMonitorPortsFunc(client *ssh.Client, uid Uid, portChan chan []int) SessionFunc {
	return func() error {
		session, err := client.NewSession()
		if err != nil {
			return errors.New("Failed to create session")
		}
		defer session.Close()
		pn, err := fetchProcNet(session)
		if err != nil {
			return errors.New("Failed to fetch port info")
		}
		ports := FindForwardablePorts(pn, uid)
		portChan <- ports
		return nil
	}
}

func createUpdateForwardingPortSession(smg SessionMG, portChan chan []int) SessionFunc {
	return func() error {
		select {
		case ports := <-portChan:
			go smg.DownPorts(ports)
			go smg.UpPorts(ports)
		default:
		}
		return nil
	}
}

func InitSshSession(config ssh.ClientConfig, addr string, remoteHost string) error {
	client, err := ssh.Dial("tcp", addr, &config)
	if err != nil {
		return err
	}
	defer client.Close()

	uid, err := fetchUid(client)
	if err != nil {
		return err
	}
	done := make(chan struct{})
	portChan := make(chan []int)
	monitorFunc := createMonitorPortsFunc(client, uid, portChan)
	createSession(monitorFunc, 1, done)
	sessionMG := NewSessionMG(remoteHost, client)
	ufp := createUpdateForwardingPortSession(*sessionMG, portChan)
	createSession(ufp, 1, done)
	for {
		time.Sleep(5 * time.Second)
	}
}

func (f *ForwardSession) connect(connChan chan net.Conn, errChan chan error) SessionFunc {
	return func() error {
		conn, err := f.listener.Accept()
		if err != nil {
			err := errors.New("Failed to accept")
			errChan <- err
			return err
		}
		connChan <- conn
		return nil
	}
}

func (f *ForwardSession) handleDataTransport(client *ssh.Client, connChan chan net.Conn, errChan chan error) SessionFunc {
	return func() error {
		select {
		case err := <-errChan:
			fmt.Println(err)
			return nil
		case localConn := <-connChan:
			defer localConn.Close()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			remoteConn, err := client.Dial("tcp", f.remoteAddr)
			if err != nil {
				err := errors.New("Faild to dial")
				fmt.Println(err)
				return err
			}
			defer remoteConn.Close()
			go func() {
				io.Copy(remoteConn, localConn)
				cancel()
			}()
			go func() {
				io.Copy(localConn, remoteConn)
				cancel()
			}()
			<-ctx.Done()
		default:
		}
		return nil
	}
}

func (f *ForwardSession) forwardPort(client *ssh.Client) {
	connChan := make(chan net.Conn)
	errChan := make(chan error)

	conn := f.connect(connChan, errChan)
	createSession(conn, 1, f.done)
	hdt := f.handleDataTransport(client, connChan, errChan)
	createSession(hdt, 1, f.done)
}
