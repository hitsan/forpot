package ssh

import (
	"bytes"
	"errors"

	"golang.org/x/crypto/ssh"
)

func fetchProcNet(session *ssh.Session) (*string, error) {
	var output bytes.Buffer
	session.Stdout = &output
	if err := session.Run("cat /proc/net/tcp"); err != nil {
		return nil, err
	}
	p := output.String()
	return &p, nil
}

func createMonitorPortsFunc(client *ssh.Client, uid Uid, portChan chan []int) func() error {
	return func() error {
		session, err := client.NewSession()
		if err != nil {
			return errors.New("Failed to create session")
		}
		defer session.Close()
		
		pn, err := fetchProcNet(session)
		if err != nil {
			return errors.New("Failed to fetch port info")
		}
		
		ports := findForwardablePorts(pn, uid)
		
		select {
		case portChan <- ports:
		default:
		}
		
		return nil
	}
}