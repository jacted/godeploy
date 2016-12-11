# Godeploy

A CLI tool for deploying files directly from your computer to a server using SSH and SFTP.

## Usage

1. Create a .json file in your project folder.
2. Edit settings in .json file.
3. Run `godeploy` inside project folder.

godeploy.json example
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

CLI flags:

`-config godeploy.json`