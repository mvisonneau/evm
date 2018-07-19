package main

import (
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestConfigureFatalText(t *testing.T) {
	l := &LogConfig{Level: "fatal", Format: "text"}
	l.configure()

	if log.GetLevel() != log.FatalLevel {
		t.Fatalf("Expected log.Level to be 'fatal' but got %s", log.GetLevel())
	}
}

func TestConfigureLoggingDefault(t *testing.T) {
	l := &LogConfig{Level: "fatal", Format: "default"}
	err := l.configure()

	if err == nil {
		t.Fatal("Expected function to return an error, got nil")
	}
}

func TestConfigureLoggingJson(t *testing.T) {
	l := &LogConfig{Level: "debug", Format: "json"}
	err := l.configure()

	if err != nil {
		t.Fatalf("Function is not expected to return an error, got '%s'", err.Error())
	}
}

func TestConfigureLoggingInvalidLogFormat(t *testing.T) {
	l := &LogConfig{Level: "foo", Format: "default"}
	err := l.configure()

	if err == nil {
		t.Fatal("Expected function to return an error, got nil")
	}
}
