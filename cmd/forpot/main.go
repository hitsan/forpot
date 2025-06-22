package main

import (
	//"fmt"
	"forpot/internal/ssh"
)

func main() {
	config := ssh.CreateSshConfig("root", "password")
	ssh.InitSshConnect(config, "2222")
	//fmt.Printf("config %+v", config)
}
