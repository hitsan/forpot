package ssh

import (
	"log"
	"net"
	"strconv"
	"strings"
)

func parseIp(ip string) net.IP {
	ipBytes := []byte{}
	for i := 6; i >= 0; i -= 2 {
		num := ip[i : i+2]
		octet, _ := strconv.ParseInt(num, 16, 0)
		ipBytes = append(ipBytes, byte(octet))
	}
	return net.IP(ipBytes)
}

func parsePort(portHex string) int {
	portI32, err := strconv.ParseInt(portHex, 16, 0)
	if err != nil {
		log.Fatal("Failed to parse port")
	}
	port := int(portI32)
	return port
}

func ParseLineForPort(line string) int {
	token := strings.Split(line, " ")
	address := token[1]
	if address[0:8] != "00000000" {
		return 0
	}
	port_hex := address[9:]
	port_32, err := strconv.ParseInt(port_hex, 16, 0)
	if err != nil {
		log.Println("Faild to parse port")
	}
	port := int(port_32)
	return port
}

//func ParsePort(str string) []int {
//		//lines := strings.Split(str, "\n")
//		ports := []int{9000, 8000}
//		return ports
//}
