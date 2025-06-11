package ssh

import (
		"strings"
		"strconv"
		"log"
)

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

func ParsePort(str string) []int {
		//lines := strings.Split(str, "\n")
		ports := []int{9000, 8000}
		return ports
}

