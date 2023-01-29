package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"time"
)

func pollMasterlist(serversChannel chan []string, masterlist string) error {
	if masterlist == "" {
		return errors.New("masterlist is required")
	}

	for {
		_, err := GetServers(GetServersOptions{
			Address: masterlist,
			Timeout: 30 * time.Second,
			OnResponse: func(servers []string) {
				serversChannel <- servers
			},
		})

		if err != nil {
			fmt.Println(err, "error polling masterlist")
		}
	}
}

func pollServerStatus(requestServersChannel chan bool, updateServersChannel chan []string) {
	servers := []string{}
	statuses := map[string]ServerStatus{}
	requestServersChannel <- true

	statusChannel := make(chan ServerStatus)

	pollChannel := make(chan bool)
	updateStatusChannel := make(chan bool)

	go func() {
		secondsSinceLastUpdate := 0
		for {
			pollChannel <- true
			time.Sleep(5 * time.Second)
			if secondsSinceLastUpdate >= 30 {
				updateStatusChannel <- true
				secondsSinceLastUpdate = 0
			}
			secondsSinceLastUpdate += 5
		}
	}()

	for {
		select {
		case newServers := <-updateServersChannel:
			servers = newServers
		case <-pollChannel:
			for _, server := range servers {
				go func(server string) {
					status, err := GetStatus(GetStatusOptions{
						Address: server,
						Timeout: 2 * time.Second,
					})
					if err != nil {
						fmt.Println(err, "error polling server"+server)
					}
					statusChannel <- status
				}(server)
			}
		case status := <-statusChannel:
			statuses[status.Info["sv_hostname"]] = status
		case <-updateStatusChannel:
			fmt.Println("writing statuses")
			// write to file
			json, err := json.Marshal(statuses)
			if err != nil {
				fmt.Println(err, "error marshalling statuses")
			}

			err = ioutil.WriteFile("statuses.json", json, 0644)
			if err != nil {
				fmt.Println(err, "error writing statuses")
			}
		}
	}
}

func main() {
	serversChannel := make(chan []string)
	requestServersChannel := make(chan bool)
	updateServersChannel := make(chan []string)
	go pollMasterlist(serversChannel, "etmaster.idsoftware.com:27950")
	go pollServerStatus(requestServersChannel, updateServersChannel)
	serverlist := Serverlist{}

	for {
		select {
		case servers := <-serversChannel:
			addCount := 0
			for _, server := range servers {
				_ = serverlist.AddServer(server)
				addCount++
			}
			if addCount > 0 {
				updateServersChannel <- serverlist.GetServers()
			}
		case <-requestServersChannel:
			updateServersChannel <- serverlist.GetServers()
		}
	}
}
