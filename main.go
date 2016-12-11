package main

import (
	"crypto/rand"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/pkg/sftp"

	"path/filepath"

	"time"

	"path"

	"golang.org/x/crypto/ssh"
)

type (
	// Files struct
	Files struct {
		Include string `json:"include"`
		Dist    string `json:"dist"`
		Backup  string `json:"backup"`
	}
	// SSH struct
	SSH struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Username string `json:"username"`
		Password string `json:"password"`
	}
	// Settings struct containing deploy settings
	Settings struct {
		PreDeploy  string `json:"pre-deploy"`
		PostDeploy string `json:"post-deploy"`
		Files      Files  `json:"files"`
		SSH        SSH    `json:"ssh"`
	}
)

var conn *ssh.Client
var sftpClient *sftp.Client
var settings *Settings
var tmpPath string

func main() {
	// Parse CLI Flags
	configPath := flag.String("config", "godeploy.json", "Config to parse.")
	flag.Parse()

	var err error

	// Get settings
	settings, err = getConfig(*configPath)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}

	// Get SSH session
	conn, err = connectToSSH()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	defer conn.Close()

	// Create sftp
	sftpClient, err = sftp.NewClient(conn)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	defer sftpClient.Close()

	// Pre deployment
	if settings.PreDeploy != "" {
		err := uploadFile(settings.PreDeploy, settings.PreDeploy)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
		runCommand("bash " + settings.PreDeploy)
	}

	// Random tmp path
	tmpPath = "/tmp/godeploy_" + randToken()

	// Delete tmpPath + create it
	sftpClient.Mkdir(tmpPath)

	// Walk local folder
	err = filepath.Walk(settings.Files.Include, walkFiles)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}

	// Backup
	if settings.Files.Backup != "" {
		t := time.Now()
		backupName := t.Format("godeploy_2006-01-02T150405")
		backupPath := path.Join(settings.Files.Backup, backupName)
		runCommand("mv " + settings.Files.Dist + " " + backupPath)
	} else {
		runCommand("rm -rf " + settings.Files.Dist)
	}

	// Move tmpPath to Dist
	runCommand("mv " + path.Join(tmpPath, settings.Files.Include) + " " + settings.Files.Dist)

	// Clean up tmpPath
	runCommand("rm -rf " + tmpPath)

	// Post deployment
	if settings.PostDeploy != "" {
		err := uploadFile(settings.PostDeploy, settings.PostDeploy)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
		runCommand("bash " + settings.PostDeploy)
	}

	fmt.Println("Done deploying!")
}

// Get and parse config
func getConfig(filePath string) (settings *Settings, err error) {
	configFile, err := os.Open(filePath)
	defer configFile.Close()
	if err != nil {
		return nil, err
	}
	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&settings); err != nil {
		return nil, err
	}
	return settings, nil
}

// Random token
func randToken() string {
	b := make([]byte, 8)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
