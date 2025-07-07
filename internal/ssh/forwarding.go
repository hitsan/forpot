package ssh

import (
	"fmt"
	"io"
	"net"
	"sync"

	"golang.org/x/crypto/ssh"
)

type ForwardSession struct {
	listener   net.Listener
	remoteAddr string
	once       sync.Once
	done       chan struct{}
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

func (f *ForwardSession) handleConnection(client *ssh.Client, localConn net.Conn) {
	defer localConn.Close()

	remoteConn, err := client.Dial("tcp", f.remoteAddr)
	if err != nil {
		fmt.Printf("Failed to dial: %v\n", err)
		return
	}
	defer remoteConn.Close()

	go func() {
		io.Copy(remoteConn, localConn)
		localConn.Close()
	}()

	io.Copy(localConn, remoteConn)
}

func (f *ForwardSession) forwardPort(client *ssh.Client) {
	go func() {
		for {
			select {
			case <-f.done:
				return
			default:
				conn, err := f.listener.Accept()
				if err != nil {
					select {
					case <-f.done:
						return
					default:
						continue
					}
				}
				go f.handleConnection(client, conn)
			}
		}
	}()
}