package main

import (
	"testing"
)

func TestRunCli(t *testing.T) {
	c := runCli()
	if c.Name != "evm" {
		t.Fatalf("Expected c.Name to be evm, got '%v'", c.Name)
	}
}
