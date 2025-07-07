package cli

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

func ReadPassword() (string, error) {
	fmt.Print("Password:")
	passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", err
	}
	fmt.Println()
	return string(passwordBytes), nil
}