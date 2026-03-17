// Package fetcher fetches the data from the given url.
package fetcher

import (
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/Siva-Sai22/crawler/urlqueue"
)

func splitURL(url string) (string, string) {
	urlParts := strings.SplitAfterN(url, "/", 4)

	baseURL := urlParts[0] + urlParts[1] + urlParts[2]
	baseURL = strings.TrimRight(baseURL, "/")

	var path string
	if len(urlParts) == 4 {
		path = "/" + urlParts[3]
	}
	return baseURL, path
}

func getRobotTxt(baseURL string) (string, bool) {
	url := baseURL + "/robots.txt"

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

func getDisallowedPaths(baseURL string) []string {
	robotsTxt, ok := getRobotTxt(baseURL)
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

func checkPathDisallowed(baseURL, path string) bool {
	disallowedPaths := getDisallowedPaths(baseURL)

	for _, disallowedPath := range disallowedPaths {
		if strings.HasPrefix(path, disallowedPath) {
			return true
		}
	}
	return false
}

func getOutboundLinks(body string, baseURL string) []string {
	links := make([]string, 0)
	for link := range strings.SplitSeq(body, "href=\"") {
		if strings.HasPrefix(link, "http") {
			links = append(links, strings.TrimSpace(link))
		} else if strings.HasPrefix(link, "/") {
			links = append(links, baseURL+link)
		}
	}
	return links
}

func Fetch(url string, urlQueue *urlqueue.URLQueue) (string, error) {
	baseURL, path := splitURL(url)

	if checkPathDisallowed(baseURL, path) {
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

	outboundLinks := getOutboundLinks(string(body), baseURL)
	urlQueue.InsertLinks(outboundLinks)

	log.Println("Crawled: ", url)
	return string(body), nil
}
