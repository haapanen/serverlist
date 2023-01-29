package main

import (
	"testing"
	"time"
)

func TestGetStatus(t *testing.T) {
	status, err := GetStatus(GetStatusOptions{
		Address: "trickjump.net:27960",
		Timeout: 2 * time.Second,
	})

	if err != nil {
		t.Errorf("Error: %s", err)
	}

	if status.Info["gamename"] != "etjump" {
		t.Errorf("gamename should be etjump")
	}
}
