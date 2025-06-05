package ssh

import "testing"

func TestSshConfig(t *testing.T) {
		user := "hitsan"
		config := CreateSshConfig(user)
		want := ssh.ClientCongig{
				User: "hitsan",

		}
		if got != want {
			t.Errorf("Add(2, 3) = %d; want %d", got, want)
		}
}
