// Package fetcher fetches the data from the given url.
package fetcher

import (
	"io"
	"log"
	"net/http"
	"strings"
)

func splitUrl(url string) (string, string) {
	urlParts := strings.SplitAfterN(url, "/", 4)

	baseUrl := urlParts[0] + urlParts[1] + urlParts[2]
	baseUrl = strings.TrimRight(baseUrl, "/")

	var path string
	if len(urlParts) == 4 {
		path = "/" + urlParts[3]
	}
	return baseUrl, path
}

func getRobotTxt(baseUrl string) (string, bool) {
	url := baseUrl + "/robots.txt"

	resp, err := http.Get(url)
	if err != nil {
		log.Println("Error fetching robots.txt: ", url, " Error: ", err)
		return "", false
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return "", false
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading robots.txt: ", err)
		return "", false
	}

	return string(body), true
}

func getDisallowedPaths(baseUrl string) []string {
	robotsTxt, ok := getRobotTxt(baseUrl)
	if !ok {
		return []string{}
	}

	lines := strings.Split(robotsTxt, "\n")
	disallowedPaths := make([]string, 0)
	for _, line := range lines {
		if strings.HasPrefix(line, "Disallow:") {
			disallowedPaths = append(disallowedPaths, strings.TrimSpace(line[9:]))
		}
	}
	return disallowedPaths
}

func checkPathDisallowed(baseUrl, path string) bool {
	disallowedPaths := getDisallowedPaths(baseUrl)

	for _, disallowedPath := range disallowedPaths {
		if strings.HasPrefix(path, disallowedPath) {
			return true
		}
	}
	return false
}

func getOutboundLinks(body string, baseUrl string) []string {
	links := make([]string, 0)
	for _, link := range strings.Split(body, "href=\"") {
		if strings.HasPrefix(link, "http") {
			links = append(links, strings.TrimSpace(link))
		} else if strings.HasPrefix(link, "/") {
			links = append(links, baseUrl+link)
		}
	}
	return links
}

func Fetch(url string) (string, error) {
	baseUrl, path := splitUrl(url)

	if checkPathDisallowed(baseUrl, path) {
		return "", nil
	}

	resp, err := http.Get(url)
	if err != nil {
		log.Println("Error fetching url: ", url, " Error: ", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return "", nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading body: ", err)
		return "", err
	}

	log.Println("Crawled: ", url)
	return string(body), nil
}
