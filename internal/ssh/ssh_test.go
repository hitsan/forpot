package ssh

import (
	"testing"
	"golang.org/x/crypto/ssh"
)

func TestSshConfig(t *testing.T) {
		user := "hitsan"
		got := CreateSshConfig(user)
		want := ssh.ClientConfig{
				User: "hitsan",
				Auth: []ssh.AuthMethod{
						ssh.Password(""),
				},
				HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}
		if got != want {
			t.Errorf("Add(2, 3) = %d; want %d", got, want)
		}
}
