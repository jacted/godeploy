package main

import (
	"fmt"

	"golang.org/x/crypto/ssh"
)

func connectToSSH() (con *ssh.Client, err error) {
	sshConfig := &ssh.ClientConfig{
		User: settings.SSH.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(settings.SSH.Password),
		},
	}
	sshPath := fmt.Sprintf("%v:%v", settings.SSH.Host, settings.SSH.Port)
	connection, err := ssh.Dial("tcp", sshPath, sshConfig)
	if err != nil {
		return nil, err
	}
	return connection, nil
}

// Creates and returns a session
func createSession(conn *ssh.Client) (session *ssh.Session, err error) {
	newSession, err := conn.NewSession()
	if err != nil {
		return nil, err
	}
	return newSession, nil
}

// Creates a session and runs a command
func runCommand(cmd string) error {
	session, err := createSession(conn)
	if err != nil {
		return err
	}
	defer session.Close()
	err = session.Run(cmd)
	if err != nil {
		return err
	}
	return nil
}
