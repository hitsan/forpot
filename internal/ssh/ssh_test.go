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
	if got.User != want.User {
		t.Errorf("User: got %v; want %v", got.User, want.User)
	}
}
