# evm

[![GoDoc](https://godoc.org/github.com/mvisonneau/evm?status.svg)](https://godoc.org/github.com/mvisonneau/evm)
[![Go Report Card](https://goreportcard.com/badge/github.com/mvisonneau/evm)](https://goreportcard.com/report/github.com/mvisonneau/evm)
[![Docker Pulls](https://img.shields.io/docker/pulls/mvisonneau/evm.svg)](https://hub.docker.com/r/mvisonneau/evm/)
[![Build Status](https://travis-ci.org/mvisonneau/evm.svg?branch=master)](https://travis-ci.org/mvisonneau/evm)
[![Coverage Status](https://coveralls.io/repos/github/mvisonneau/evm/badge.svg?branch=master)](https://coveralls.io/github/mvisonneau/evm?branch=master)

This projects aims to ease the mount and configuration of EBS volumes onto ASG backed EC2 instances. It is inspired of [monzo/etcd3-bootstrap](https://github.com/monzo/etcd3-bootstrap) project.

## TL;DR

```
~$ wget https://github.com/mvisonneau/evm/releases/download/0.0.1/evm_linux_amd64 -O /usr/local/bin/evm; chmod +x /usr/local/bin/evm
~$ evm mount -b /dev/xvdf -f ext4 -m /mnt/foo -v my_ebs_volume_name
```

## Usage

```
~$ evm
NAME:
   evm - Mount and configure an EBS volume within and onto an EC2 instance

USAGE:
   evm [global options] command [command options] [arguments...]

VERSION:
   <devel>

COMMANDS:
     mount    mount and configure an EBS volume
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --log-level value   log level (debug,info,warn,fatal,panic) (default: "info") [$EVM_LOG_LEVEL]
   --log-format value  log format (json,text) (default: "text") [$EVM_LOG_FORMAT]
   --help, -h          show help
   --version, -v       print the version
```

## Develop

If you have docker locally, you can use the following command in order to quickly get a development env ready: `make dev-env`. You can also have a look onto the [Makefile](/Makefile) in order to see all available options:

```
~$ make
all                            Test, builds and ship package for all supported platforms
build                          Build the binary
clean                          Remove binary if it exists
coverage                       Generates coverage report
deps                           Fetch all dependencies
dev-env                        Build a local development environment using Docker
fmt                            Format source code
help                           Displays this help
imports                        Fixes the syntax (linting) of the codebase
install                        Build and install locally the binary (dev purpose)
lint                           Run golint and go vet against the codebase
publish-github                 Publish the compiled binaries onto the GitHub release API
publish-goveralls              Publish coverage stats on goveralls
setup                          Install required libraries/tools
test                           Run the tests against the codebase
```

## Contribute

Contributions are more than welcome! Feel free to submit a [PR](https://github.com/mvisonneau/evm/pulls).
