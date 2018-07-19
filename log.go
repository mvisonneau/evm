package main

import (
	"errors"
	"os"

	log "github.com/sirupsen/logrus"
)

type LogConfig struct {
	Level  string
	Format string
}

var logConfig LogConfig

func (c *LogConfig) configure() error {
	parsedLevel, _ := log.ParseLevel(c.Level)
	log.SetLevel(parsedLevel)

	formatter := &log.TextFormatter{
		FullTimestamp: true,
	}
	log.SetFormatter(formatter)

	switch c.Format {
	case "text":
	case "json":
		log.SetFormatter(&log.JSONFormatter{})
	default:
		return errors.New("Invalid log format")
	}

	log.SetOutput(os.Stdout)

	return nil
}
