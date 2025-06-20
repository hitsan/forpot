package ssh

import (
	"errors"
	"os/user"
	"strconv"
	"strings"
)

func canListen(status string) bool {
	return status == "0A"
}

func equalsUid(uid string, targetUid string) bool {
	return uid == targetUid
}

func isLocalhost(ip string) bool {
	return ip == "00000000"
}

func canPortForward(line string, uid string) bool {
	items := strings.Fields(line)
	address := items[1]
	isLocalhostIp := isLocalhost(address[:8])
	canLitesned := canListen(items[3])
	isEquqlsUid := equalsUid(uid, items[7])
	return isLocalhostIp && canLitesned && isEquqlsUid
}

func parsePort(portHex string) string {
	portI64, _ := strconv.ParseInt(portHex, 16, 64)
	port := strconv.FormatInt(portI64, 10)
	return port
}

func GetUid() (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", errors.New("Failed to get user info")
	}
	uid := user.Uid
	return uid, nil
}

func parseLine(line string, uid string) (string, error) {
	cpf := canPortForward(line, uid)
	if !cpf {
		return "", errors.New("This port is not forwardable")
	}
	items := strings.Fields(line)
	address := items[1]
	portHex := address[9:]
	port := parsePort(portHex)
	return port, nil
}

func FindForwardablePorts(netInfo string, uid string) []string {
	splitedNetInfo := strings.Split(netInfo, "\n")
	lines := splitedNetInfo[1:]
	var ports []string
	for _, line := range lines {
		port, err := parseLine(line, uid)
		if err == nil {
			ports = append(ports, port)
		}
	}
	return ports
}
