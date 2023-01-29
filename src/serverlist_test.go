package main

import "testing"

// test that adding server works
func TestAddServer(t *testing.T) {
	s := Serverlist{}
	s.AddServer("test")
	if s.servers[0] != "test" {
		t.Errorf("Expected 'test', got %s", s.servers[0])
	}
}

// test that adding duplicate server doesn't add another server
func TestAddDuplicateServer(t *testing.T) {
	s := Serverlist{}
	s.AddServer("test")
	s.AddServer("test")
	if len(s.servers) != 1 {
		t.Errorf("Expected 1 server, got %d", len(s.servers))
	}
}

// test that removing server works
func TestRemoveServer(t *testing.T) {
	s := Serverlist{}
	s.AddServer("test")
	s.RemoveServer("test")
	if len(s.servers) != 0 {
		t.Errorf("Expected 0 servers, got %d", len(s.servers))
	}
}

// test that removing non-existing server doesn't remove anything
func TestRemoveNonExistingServer(t *testing.T) {
	s := Serverlist{}
	s.AddServer("test")
	s.RemoveServer("test2")
	if len(s.servers) != 1 {
		t.Errorf("Expected 1 server, got %d", len(s.servers))
	}
}

// test that getting servers works
func TestServerlistGetServers(t *testing.T) {
	s := Serverlist{}
	s.AddServer("test")
	if s.GetServers()[0] != "test" {
		t.Errorf("Expected 'test', got %s", s.GetServers()[0])
	}
}
