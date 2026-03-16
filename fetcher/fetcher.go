// Package fetcher fetches the data from the given url.
package fetcher

import (
	"io"
	"log"
	"net/http"
)

func Fetch(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Error fetching url: ", url, " Error: ", err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading body: ", err)
		return "", err
	}

	log.Println("Crawled: ", url)
	return string(body), nil
}
