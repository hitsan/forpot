package config

import (
	"forpot/internal/ssh"
	gossh "golang.org/x/crypto/ssh"
)

func CreateSSHConfig(user, password string) gossh.ClientConfig {
	return ssh.CreateSshConfig(user, password)
}