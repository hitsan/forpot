package ssh

import (
	"golang.org/x/crypto/ssh"
	"testing"
)

func TestSshConfig(t *testing.T) {
	user := "hitsan"
	got := CreateSshConfig(user, "password")
	want := ssh.ClientConfig{
		User: "hitsan",
		Auth: []ssh.AuthMethod{
			ssh.Password("password"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	if got.User != want.User {
		t.Errorf("User: got %v; want %v", got.User, want.User)
	}
}
