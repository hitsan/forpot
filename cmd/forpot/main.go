package main

import (
	//"fmt"
	"forpot/internal/ssh"
)

func main() {
	config := ssh.CreateSshConfig("root", "password")
	ssh.Connect(config, "127.0.0.1", "2222")
	//fmt.Printf("config %+v", config)
}
