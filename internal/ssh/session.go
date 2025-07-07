package ssh

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
)

type SessionMG struct {
	remoteHost string
	client     *ssh.Client
	sessionMap map[int]*ForwardSession
	mu         sync.RWMutex
}

type SessionFunc func() error

func NewSessionMG(remoteHost string, client *ssh.Client) *SessionMG {
	return &SessionMG{
		remoteHost: remoteHost,
		sessionMap: make(map[int]*ForwardSession),
		client:     client,
	}
}

func (s *SessionMG) getPortMap() map[int]struct{} {
	s.mu.Lock()
	defer s.mu.Unlock()

	pm := make(map[int]struct{})
	for port := range s.sessionMap {
		pm[port] = struct{}{}
	}
	return pm
}

func (s *SessionMG) selectUpdatePorts(portChan chan []int, upPortChan chan int, downPortChan chan int) {
	go func() {
		for {
			select {
			case ports := <-portChan:
				pm := make(map[int]struct{})
				for _, port := range ports {
					pm[port] = struct{}{}
				}
				portMap := s.getPortMap()
				for port := range portMap {
					_, ok := pm[port]
					if ok {
						continue
					}
					downPortChan <- port
				}

				for _, port := range ports {
					_, ok := portMap[port]
					if ok {
						continue
					}
					upPortChan <- port
				}
			default:
			}
		}
	}()
}

func (s *SessionMG) UpdateForwardingSession(upPortChan chan int, downPortChan chan int) {
	go func() {
		for {
			select {
			case port := <-downPortChan:
				session := s.sessionMap[port]
				session.Close()
				delete(s.sessionMap, port)
			case port := <-upPortChan:
				localAddr := fmt.Sprintf("127.0.0.1:%d", port)
				remoteAddr := fmt.Sprintf("%s:%d", s.remoteHost, port)
				fs, err := NewForwardSession(localAddr, remoteAddr)
				if err != nil {
					fmt.Println(err)
					continue
				}
				go fs.forwardPort(s.client)
				s.sessionMap[port] = fs
			default:
			}
		}
	}()
}

func (s *SessionMG) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, session := range s.sessionMap {
		session.Close()
	}
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

func createSession(fn SessionFunc, ms time.Duration, done chan struct{}) {
	go func() {
		ticker := time.NewTicker(time.Duration(ms) * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-done:
				fmt.Println("close")
				return
			case <-ticker.C:
				err := fn()
				if err != nil {
					fmt.Printf("Session function failed: %v\n", err)
				}
			}
		}
	}()
}

func InitSshSession(ctx context.Context, config ssh.ClientConfig, addr string, remoteHost string) error {
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
	portChan := make(chan []int, 100)

	monitorFunc := createMonitorPortsFunc(client, uid, portChan)
	createSession(monitorFunc, 1000, done)

	sessionMG := NewSessionMG(remoteHost, client)
	defer sessionMG.Close()

	ufp := createUpdateForwardingPortSession(sessionMG, portChan)
	createSession(ufp, 1000, done)

	select {
	case <-ctx.Done():
		close(done)
		return nil
	}
}

func createUpdateForwardingPortSession(smg *SessionMG, portChan chan []int) SessionFunc {
	return func() error {
		select {
		case <-portChan:
			upPortChan := make(chan int)
			downPortChan := make(chan int)

			smg.selectUpdatePorts(portChan, upPortChan, downPortChan)
			smg.UpdateForwardingSession(upPortChan, downPortChan)
		default:
		}
		return nil
	}
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
