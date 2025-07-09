package ssh

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"syscall"

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

func SetupPortForwarding(client *ssh.Client, remoteHost string, port int, sync *SessionSynchronizer) {
	remoteAddr := fmt.Sprintf("%s:%d", remoteHost, port)
	
	for count := 0; count < 10; count++ {
		localAddr := fmt.Sprintf("127.0.0.1:%d", port+count)
		fs, err := NewForwardSession(localAddr, remoteAddr)
		if err != nil {
			if errors.Is(err, syscall.EADDRINUSE) {
				fmt.Println(count)
				continue
			}
			fmt.Println(err)
			return
		}
		go fs.forwardPort(client)
		sync.Set(port, fs)
		fmt.Println("forward port: ", port+count)
		return
	}
}
