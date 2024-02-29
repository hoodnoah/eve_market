package idcache

import "net/http"

func getClient() *http.Client {
	newClient := http.Client{}

	return &newClient
}
