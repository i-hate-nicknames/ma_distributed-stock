package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// perform a simple GET request to the given warehouse
// address with a given action. Return data read from
// the response
func callWarehouse(address, action string) ([]byte, error) {
	resp, err := http.Get("http://" + address + "/" + action)
	if err != nil {
		return nil, fmt.Errorf("Error getting warehouse items, address: %s, error: %s", address, err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading response data: %s", err)
	}
	return body, nil
}
