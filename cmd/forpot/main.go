package main

import (
	"forpot/internal/ssh"
	"net"
)

func main() {
	config := ssh.CreateSshConfig("root", "password")
	addr := net.TCPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 2222,
	}
	remoteAddr := net.ParseIP("127.0.0.1")
	ssh.InitSshSession(config, addr, remoteAddr)
	//fmt.Printf("config %+v", config)
}
