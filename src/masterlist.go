package main

import (
	"errors"
	"fmt"
	"net"
	"time"
)

type GetServersOptions struct {
	/**
	 * The address of the masterlist
	 */
	Address string
	/**
	 * How long will the function wait for more responses from the
	 * masterlist before returning
	 */
	Timeout time.Duration
	/**
	 * Called when a response is received from the masterlist.
	 * The response is a list of servers. If you wish to get the updates
	 * as they come in instead of waiting for the timeout, implement
	 * this function
	 */
	OnResponse func([]string)
}

var defaultOptions = GetServersOptions{
	Timeout:    10 * time.Second,
	OnResponse: func(servers []string) {},
}

/**
 * Get a list of servers from a masterlist.
 *
 * @param options
 * @return
 */
func GetServers(options GetServersOptions) ([]string, error) {
	timeout := options.Timeout
	if timeout == 0 {
		timeout = defaultOptions.Timeout
	}
	address := options.Address
	if address == "" {
		return nil, errors.New("address is required")
	}
	onResponse := options.OnResponse
	if onResponse == nil {
		onResponse = defaultOptions.OnResponse
	}

	if timeout.Abs() == 0 {
		timeout = 30 * time.Second
	}

	conn, err := net.Dial("udp", address)

	if err != nil {
		return nil, err
	}

	defer conn.Close()

	_, err = conn.Write([]byte("\xff\xff\xff\xffgetservers 84 full empty"))
	if err != nil {
		return nil, err
	}

	buf := make([]byte, 1024)
	conn.SetReadDeadline(time.Now().Add(timeout))

	var servers []string
	for {
		n, err := conn.Read(buf)
		if err, ok := err.(net.Error); ok && err.Timeout() {
			break
		} else if err != nil {
			return nil, err
		}

		receivedServers, err := parseResponse(buf[:n])
		if err != nil {
			return nil, err
		}
		servers = append(servers, receivedServers...)
		onResponse(receivedServers)
	}

	return servers, nil
}

/**
 * Parse the response from the masterlist to a list of ip:port strings
 *
 * @param buf
 * @return
 */
func parseResponse(buf []byte) ([]string, error) {
	var servers []string

	if len(buf) < 22 {
		return nil, errors.New("invalid response")
	}

	message := buf[22:]
	for i := 0; i < len(message); i = i + 8 {
		if i+8 >= len(message) {
			break
		}
		ip := net.IPv4(message[i+1], message[i+2], message[i+3], message[i+4])
		port := int(message[i+5])<<8 + int(message[i+6])

		servers = append(servers, fmt.Sprintf("%s:%d", ip, port))
	}

	return servers, nil
}
