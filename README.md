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
go build win_helper/cmd/win_helper
```

show help `win_helper -h`
```bash
win_helper is a CLI generator for windows service script

Usage:
  win_helper [flags]
  win_helper [command]

Available Commands:
  completion   Generate the autocompletion script for the specified shell
  help         Help about any command
  initLanguage init language directory
  initProject  init project directory
  mklink       windows make link
  obr          obr tools

Flags:
      --config string   config file (default is $HOME/.win_helper.yaml)
  -h, --help            help for win_helper
  -v, --version         version
      --viper           use Viper for configuration (default true)

Use "win_helper [command] --help" for more information about a command.
```

## Architecture
```bash

```