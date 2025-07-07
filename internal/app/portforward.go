package app

import (
	"fmt"
	"forpot/internal/cli"
	"forpot/internal/config"
	"forpot/internal/ssh"
	"log"
)

func RunPortForwarding(hostArg string, port int) {
	user, host, err := cli.ParseHost(hostArg)
	if err != nil {
		log.Fatalln(err)
	}

	password, err := cli.ReadPassword()
	if err != nil {
		log.Fatalln(err)
	}

	sshConfig := config.CreateSSHConfig(user, password)
	addr := fmt.Sprintf("%s:%d", "127.0.0.1", port)
	remoteHost := host
	err = ssh.InitSshSession(sshConfig, addr, remoteHost)
	if err != nil {
		log.Fatalln(err)
	}
}