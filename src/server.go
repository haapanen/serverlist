package main

import (
	"errors"
	"net"
	"strings"
	"time"
)

type GetStatusOptions struct {
	Address string
	Timeout time.Duration
}

var defaultGetStatusOptions = GetStatusOptions{
	Timeout: 2 * time.Second,
}

type ServerStatus struct {
	Info    map[string]string
	Players []string
}

/**
 * Get the status of a server.
 *
 * @param options
 * @return
 */
func GetStatus(options GetStatusOptions) (ServerStatus, error) {
	address := options.Address
	if address == "" {
		return ServerStatus{}, errors.New("address is required")
	}
	timeout := options.Timeout
	if timeout == 0 {
		timeout = defaultGetStatusOptions.Timeout
	}

	conn, err := net.Dial("udp", address)
	if err != nil {
		return ServerStatus{}, err
	}
	defer conn.Close()

	_, err = conn.Write([]byte("\xff\xff\xff\xffgetstatus"))
	if err != nil {
		return ServerStatus{}, err
	}

	buf := make([]byte, 1024)
	conn.SetReadDeadline(time.Now().Add(timeout))
	// read until timeout
	n, err := conn.Read(buf)
	if err != nil {
		return ServerStatus{}, err
	}

	return parseGetStatusResponse(buf[:n])
}

func parseGetStatusResponse(statusResponseBuffer []byte) (ServerStatus, error) {
	statusResponseString := string(statusResponseBuffer)
	lines := strings.Split(statusResponseString, "\n")
	if len(lines) < 2 {
		return ServerStatus{}, errors.New("invalid response from server")
	}

	players := []string{}
	if len(lines) > 2 {
		players = parsePlayers(lines[2 : len(lines)-1])
	}

	return ServerStatus{
		Info:    parseInfo(lines[1]),
		Players: players,
	}, nil
}

func parseInfo(infoLine string) map[string]string {
	info := make(map[string]string)
	infoParts := strings.Split(infoLine, "\\")
	for i := 1; i < len(infoParts)-1; i += 2 {
		info[infoParts[i]] = infoParts[i+1]
	}
	return info
}

func parsePlayers(playerLines []string) []string {
	players := make([]string, len(playerLines))
	for i, playerLine := range playerLines {
		players[i] = strings.Split(playerLine, " ")[1]
	}
	return players
}
