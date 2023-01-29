package main

import (
	testutilities "haapanen/serverlist/src/testUtilities"
	"testing"
	"time"
)

func TestGetServers(t *testing.T) {
	masterlist := "etmaster.idsoftware.com:27950"

	servers, err := GetServers(GetServersOptions{
		Address: masterlist,
		Timeout: 2 * time.Second,
	})

	if err != nil {
		t.Errorf("Error: %s", err)
	}

	if len(servers) == 0 {
		t.Errorf("No servers found")
	}
}

func TestGetServersCallback(t *testing.T) {
	masterlist := "etmaster.idsoftware.com:27950"

	callCount := 0

	GetServers(GetServersOptions{
		Address: masterlist,
		Timeout: 2 * time.Second,
		OnResponse: func(servers_ []string) {
			callCount++
		},
	})

	if callCount == 0 {
		t.Errorf("Callback should've been called at least once")
	}
}

func TestWithTestServer(t *testing.T) {
	testServer := testutilities.NewTestServer()
	testServer.Start(27950)

	servers, err := GetServers(GetServersOptions{
		Address: "localhost:27950",
		Timeout: 1 * time.Second,
	})

	if err != nil {
		t.Errorf("Error: %s", err)
	}

	if len(servers) == 0 {
		t.Errorf("No servers found")
	}

	testServer.Stop()
}

func TestInvalidGetServersResponse(t *testing.T) {
	testServer := testutilities.NewTestServer()
	testServer.Start(27950)

	testServer.SetGetServersResponse([]byte("invalid"))

	_, err := GetServers(GetServersOptions{
		Address: "localhost:27950",
		Timeout: 1 * time.Second,
	})

	if err == nil {
		t.Errorf("Expected error")
	}

	testServer.Stop()
}
