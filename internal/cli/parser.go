package cli

import (
	"errors"
	"os/user"
	"strings"
)

func ParseHost(arg string) (string, string, error) {
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