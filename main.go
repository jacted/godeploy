package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"

	"gopkg.in/urfave/cli.v1"

	"github.com/pkg/sftp"

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
	var err error

	app := cli.NewApp()
	app.Name = "godeploy"
	app.Usage = "Deploy files to server using SSH and SFTP"
	app.Version = "0.0.2"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Value: "godeploy.json",
			Usage: "json file with settings",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "ssh",
			Usage: "deploys using SSH and SFTP",
			Action: func(c *cli.Context) error {
				// Flags
				configPath := c.GlobalString("config")

				// Get settings
				settings, err = getConfig(configPath)
				if err != nil {
					return err
				}

				// Get SSH session
				conn, err = connectToSSH()
				if err != nil {
					return err
				}
				defer conn.Close()

				// Create sftp
				sftpClient, err = sftp.NewClient(conn)
				if err != nil {
					return err
				}
				defer sftpClient.Close()

				// Pre deployment
				if settings.PreDeploy != "" {
					err = uploadAndRunScript(settings.PreDeploy)
					if err != nil {
						return err
					}
				}

				// Random tmp path
				tmpPath = "/tmp/godeploy_" + randToken()

				// Delete tmpPath + create it
				err = sftpClient.Mkdir(tmpPath)
				if err != nil {
					return err
				}

				// Walk local folder
				err = filepath.Walk(settings.Files.Include, walkFiles)
				if err != nil {
					return err
				}

				// Backup
				if settings.Files.Backup != "" {
					t := time.Now()
					backupName := t.Format("godeploy_2006-01-02T150405")
					backupPath := path.Join(settings.Files.Backup, backupName)
					err = runCommand("mv " + settings.Files.Dist + " " + backupPath)
					if err != nil {
						return err
					}
				} else {
					err = runCommand("rm -rf " + settings.Files.Dist)
					if err != nil {
						return err
					}
				}

				// Move tmpPath to Dist
				err = runCommand("mv " + path.Join(tmpPath, settings.Files.Include) + " " + settings.Files.Dist)
				if err != nil {
					return err
				}

				// Clean up tmpPath
				err = runCommand("rm -rf " + tmpPath)
				if err != nil {
					return err
				}

				// Post deployment
				if settings.PostDeploy != "" {
					err := uploadAndRunScript(settings.PostDeploy)
					if err != nil {
						return err
					}
				}

				fmt.Println("Done deploying!")

				// Return nil
				return nil
			},
		},
	}

	app.Run(os.Args)
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
