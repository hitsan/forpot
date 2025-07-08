package ssh

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

type SessionManager struct {
	remoteHost string
	client     *ssh.Client
	sync       *SessionSynchronizer
}

type SessionFunc func() error

func NewSessionManager(remoteHost string, client *ssh.Client) *SessionManager {
	return &SessionManager{
		remoteHost: remoteHost,
		client:     client,
		sync:       NewSessionSynchronizer(),
	}
}

func (s *SessionManager) getPortMap() map[int]struct{} {
	return s.sync.GetAll()
}

func (s *SessionManager) selectUpdatePorts(portChan chan []int, upPortChan chan int, downPortChan chan int) {
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

func (s *SessionManager) UpdateForwardingSession(upPortChan chan int, downPortChan chan int) {
	go func() {
		for {
			select {
			case port := <-downPortChan:
				s.sync.Delete(port)
			case port := <-upPortChan:
				localAddr := fmt.Sprintf("127.0.0.1:%d", port)
				remoteAddr := fmt.Sprintf("%s:%d", s.remoteHost, port)
				fs, err := NewForwardSession(localAddr, remoteAddr)
				if err != nil {
					fmt.Println(err)
					continue
				}
				go fs.forwardPort(s.client)
				s.sync.Set(port, fs)
			default:
			}
		}
	}()
}

func (s *SessionManager) Close() {
	s.sync.Close()
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
	go func() {
		ticker := time.NewTicker(1000 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				if err := monitorFunc(); err != nil {
					fmt.Printf("Monitor function failed: %v\n", err)
				}
			}
		}
	}()

	sessionMgr := NewSessionManager(remoteHost, client)
	defer sessionMgr.Close()

	ufp := createUpdateForwardingPortSession(sessionMgr, portChan)
	go func() {
		ticker := time.NewTicker(1000 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				if err := ufp(); err != nil {
					fmt.Printf("Update forwarding function failed: %v\n", err)
				}
			}
		}
	}()

	select {
	case <-ctx.Done():
		close(done)
		return nil
	}
}

func createUpdateForwardingPortSession(smg *SessionManager, portChan chan []int) SessionFunc {
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
