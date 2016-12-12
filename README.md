# Godeploy

[![Go Report Card](https://goreportcard.com/badge/github.com/jacted/godeploy)](https://goreportcard.com/report/github.com/jacted/godeploy)
[![MIT licensed](https://img.shields.io/github/license/jacted/godeploy.svg?maxAge=2592000)](https://github.com/jacted/godeploy/blob/master/LICENSE)

A CLI tool for deploying files directly from your computer to a server using SSH and SFTP.

## Installation

1. Download [latest release](https://github.com/jacted/godeploy/releases)
2. Put it in a folder and make sure its executeable
3. Add the folder to system PATH
4. You're set, you can now run it using `godeploy` in terminal

## Usage

1. Create a .json file in your project folder.
2. Edit settings in .json file.
3. Run `godeploy` inside project folder.

godeploy.json example:
``` json
{
  "files": {
    "include": "dist/",
    "dist": "/var/www/html",
    "backup": "/var/backup"
  },
  "ssh": {
    "host": "localhost",
    "port": 22,
    "username": "root",
    "password": "1234"
  },
  "pre-deploy": "pre-deploy.sh",
  "post-deploy": ""
}
```

## Commands:

|`godeploy <command>`|Description|
|------------------|-----------|
|`help`|Shows help description.|
|`ssh`|Deploys using SSH and SFTP.|

## CLI Flags

`--config PATH_JSON_FILE`

## Questions and issues

Use the Github issue tracker for any questions, issues or feature suggestions.

## License

MIT licensed.