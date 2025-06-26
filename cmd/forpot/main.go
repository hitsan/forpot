package main

import (
	"flag"
	"fmt"
	"forpot/internal/ssh"
	"net"
)

func main() {
	flag.Parse()
	fmt.Println(flag.Args())
	config := ssh.CreateSshConfig("root", "password")
	addr := net.TCPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 2222,
	}
	remoteAddr := net.ParseIP("127.0.0.1")
	ssh.InitSshSession(config, addr, remoteAddr)
	//fmt.Printf("config %+v", config)
}
