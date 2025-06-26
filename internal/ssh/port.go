package ssh

import (
	"errors"
	"strconv"
	"strings"
)

type Uid string

func canListen(status string) bool {
	return status == "0A"
}

func equalsUid(uid Uid, targetUid Uid) bool {
	return uid == targetUid
}

func isLocalhost(ip string) bool {
	return ip == "00000000"
}

func isWellKnownPort(portStr string) bool {
	i, err := strconv.ParseInt(portStr, 16, 64)
	port := int(i)
	if err != nil {
		return true
	}
	if port < 1024 {
		return true
	}
	return false
}

func canPortForward(line string, uid Uid) bool {
	items := strings.Fields(line)
	address := items[1]
	isLocalhostIp := isLocalhost(address[:8])
	iswkp := isWellKnownPort(address[9:])
	canLitesned := canListen(items[3])
	targetUid := Uid(items[7])
	isEquqlsUid := equalsUid(uid, targetUid)
	return isLocalhostIp && canLitesned && isEquqlsUid && (!iswkp)
}

func parsePort(portHex string) (int, error) {
	portI64, err := strconv.ParseInt(portHex, 16, 64)
	if err != nil {
		return 0, errors.New("Failed to parse port")
	}
	return int(portI64), nil
}

func parseLine(line string, uid Uid) (int, error) {
	cpf := canPortForward(line, uid)
	if !cpf {
		return 0, errors.New("This port is not forwardable")
	}
	items := strings.Fields(line)
	address := items[1]
	portHex := address[9:]
	port, err := parsePort(portHex)
	if err != nil {
		return 0, err
	}
	return port, nil
}

func FindForwardablePorts(netInfo *string, uid Uid) []int {
	splitedNetInfo := strings.Split(*netInfo, "\n")
	length := len(splitedNetInfo)
	lines := splitedNetInfo[1 : length-1]
	var ports []int
	for _, line := range lines {
		port, err := parseLine(line, uid)
		if err != nil {
			continue
		}
		ports = append(ports, port)
	}
	return ports
}
