package main

import (
	"fmt"
	"os"

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
func createSession(conn *ssh.Client) (session *ssh.Session) {
	session, err := conn.NewSession()
	if err != nil {
		fmt.Println("Error establishing new session")
		os.Exit(0)
	}
	return session
}

// Creates a session and runs a command
func runCommand(cmd string) {
	session := createSession(conn)
	defer session.Close()
	session.Run(cmd)
}
