package app

import (
	"fmt"
	"forpot/internal/cli"
	"forpot/internal/config"
	"forpot/internal/signal"
	"forpot/internal/ssh"
	"log"
	"sync"
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

	signalHandler := signal.New()
	ctx := signalHandler.Context()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		err = ssh.InitSshSession(ctx, sshConfig, addr, remoteHost)
		if err != nil {
			log.Printf("SSH session error: %v", err)
		}
	}()

	signalHandler.Wait()
	wg.Wait()
	signalHandler.Shutdown()
	fmt.Println("Port forwarding stopped.")
}