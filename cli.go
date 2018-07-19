package main

import (
	"github.com/urfave/cli"
)

var version = "<devel>"

// runCli : Generates cli configuration for the application
func runCli() (c *cli.App) {
	c = cli.NewApp()
	c.Name = "evm"
	c.Version = version
	c.Usage = "Mount and configure an EBS volume within and onto an EC2 instance"
	c.EnableBashCompletion = true

	c.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "log-level",
			EnvVar:      "EVM_LOG_LEVEL",
			Usage:       "log level (debug,info,warn,fatal,panic)",
			Value:       "info",
			Destination: &logConfig.Level,
		},
		cli.StringFlag{
			Name:        "log-format",
			EnvVar:      "EVM_LOG_FORMAT",
			Usage:       "log format (json,text)",
			Value:       "text",
			Destination: &logConfig.Format,
		},
	}

	c.Commands = []cli.Command{
		{
			Name:      "mount",
			Usage:     "mount and configure an EBS volume",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "block-device-name, b",
					EnvVar: "EVM_BLOCK_DEVICE_NAME",
					Usage:  "name of the block device to use locally (required)",
				},
				cli.StringFlag{
					Name:   "filesystem-type, f",
					EnvVar: "EVM_FILESYSTEM_TYPE",
					Usage:  "type of filesystem to use for the volume",
					Value:  "ext4",
				},
				cli.StringFlag{
					Name:   "mount-point, m",
					EnvVar: "EVM_MOUNT_POINT",
					Usage:  "location to mount the volume (required)",
				},
				cli.StringFlag{
					Name:   "volume-name, v",
					EnvVar: "EVM_VOLUME_NAME",
					Usage:  "name of the EBS volume to attach to this instance (required)",
				},
			},
			Action: mount,
		},
	}

	return
}
