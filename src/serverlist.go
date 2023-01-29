package main

type Serverlist struct {
	servers []string
}

func (s *Serverlist) AddServer(server string) bool {
	if server == "" {
		return false
	}

	for _, v := range s.servers {
		if v == server {
			return false
		}
	}

	s.servers = append(s.servers, server)
	return true
}

func (s *Serverlist) RemoveServer(server string) {
	for i, v := range s.servers {
		if v == server {
			s.servers = append(s.servers[:i], s.servers[i+1:]...)
			break
		}
	}
}

func (s *Serverlist) GetServers() []string {
	return s.servers
}
