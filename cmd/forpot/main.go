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
	ssh.InitSshSession(config, addr)
	//fmt.Printf("config %+v", config)
}
