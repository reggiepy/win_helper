# golang gen windows helper

[![go version](https://img.shields.io/github/go-mod/go-version/reggiepy/win_helper?color=success&filename=go.mod&style=flat)](https://github.com/reggiepy/win_helper)
[![release](https://img.shields.io/github/v/tag/reggiepy/win_helper?color=success&label=release)](https://github.com/reggiepy/win_helper)
[![build status](https://img.shields.io/badge/build-pass-success.svg?style=flat)](https://github.com/reggiepy/win_helper)
[![License](https://img.shields.io/badge/license-GNU%203.0-success.svg?style=flat)](https://github.com/reggiepy/win_helper)
[![Go Report Card](https://goreportcard.com/badge/github.com/reggiepy/win_helper)](https://goreportcard.com/report/github.com/reggiepy/win_helper)

## Installation

```bash
git clone https://github.com/reggiepy/win_helper.git
cd win_helper
go mod tidy
```

## Usage

```bash
go run win_helper/cmd/win_helper
go generate win_helper/cmd/win_helper
go build win_helper/cmd/win_helper
```

show help `win_helper -h`
```bash
Use "win_helper [command] --help" for more information about a command.
```

generic service
```bash
win_helper.exe winserver-gen --name minio --executable minio.exe --description minio --start-arguments "server minio"
win_helper.exe winserver-gen --name bbs-go --executable bbs-go.exe --description bbs-go
win_helper.exe winserver-gen --name gitea --executable gitea.exe --description gitea --start-arguments "web"
win_helper.exe winserver-gen --name frpc-remote --executable frpc.exe --description frpc --start-arguments "-c frpcremote.ini"
win_helper.exe winserver-gen --name supervisord --executable supervisord.exe --description supervisord --start-arguments "-c supervisord.conf"
win_helper.exe winserver-gen --name nsqlookupd --executable nsqlookupd.exe --description nsqlookupd --working-directory bin
win_helper.exe winserver-gen --name nsq-auth --executable nsq-auth.exe --description nsq-auth --start-arguments "--secret %n&yFA2JD85z^g --auth-http-address 127.0.0.1:1325" --working-directory bin
win_helper.exe winserver-gen --name nsqd --executable nsqd.exe --description nsqd --start-arguments "--lookupd-tcp-address=127.0.0.1:4160 --auth-http-address "127.0.0.1:1325"" --working-directory bin
win_helper.exe winserver-gen --name nsqadmin --executable nsqadmin.exe --description nsqadmin --start-arguments "--lookupd-http-address=127.0.0.1:4161 -u "127.0.0.1:1325"" --working-directory bin

win_helper.exe winserver-gen --name minio --executable minio.exe --description minio --start-arguments "server minio"
```
## Architecture
```bash

```

