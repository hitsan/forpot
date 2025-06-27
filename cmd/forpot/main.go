package main

import (
	"errors"
	"fmt"
	"forpot/internal/ssh"
	"github.com/spf13/cobra"
	"log"
	"net"
	"os"
	"os/user"
	"strings"

	"golang.org/x/term"
)

func parseHost(arg string) (string, string, error) {
	items := strings.Split(arg, "@")
	if len(items) > 2 {
		return "", "", errors.New("Illigal connection")
	}
	if len(items) == 2 {
		return items[0], items[1], nil
	}
	u, err := user.Current()
	if err != nil {
		return "", "", errors.New("Cannot get user name")
	}
	return u.Username, items[0], nil
}

func main() {
	var port int
	cmd := &cobra.Command{
		Use:   "forpot [user@]host",
		Short: "Port forwarding app",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			user, host, err := parseHost(args[0])
			if err != nil {
				log.Fatalln(err)
				os.Exit(1)
			}

			fmt.Print("Password:")
			passwordsBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				log.Fatalln(err)
				os.Exit(1)
			}
			fmt.Println()
			password := string(passwordsBytes)

			config := ssh.CreateSshConfig(user, password)
			addr := net.TCPAddr{
				IP:   net.ParseIP(host),
				Port: port,
			}
			remoteHost := "127.0.0.1"
			ssh.InitSshSession(config, addr, remoteHost)
		},
	}
	cmd.Flags().IntVar(&port, "port", 22, "Set ssh port")

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
