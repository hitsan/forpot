package ssh

import (
	"log"
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

func canPortForward(line string) bool {
	items := strings.Fields(line)
	address := items[1]
	isLocalhostIp := isLocalhost(address[:8])
	canLitesned := canListen(items[3])
	isEquqlsUid := equalsUid("0", items[7])
	return isLocalhostIp && canLitesned && isEquqlsUid
}

func ParseLineForPort(line string) int {
	token := strings.Split(line, " ")
	address := token[1]
	if address[0:8] != "00000000" {
		return 0
	}
	portHex := address[9:]
	portI32, err := strconv.ParseInt(portHex, 16, 0)
	if err != nil {
		log.Println("Faild to parse port")
	}
	port := int(portI32)
	return port
}

//func ParsePort(str string) []int {
//		//lines := strings.Split(str, "\n")
//		ports := []int{9000, 8000}
//		return ports
//}
