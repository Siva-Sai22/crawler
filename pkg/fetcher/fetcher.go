// Package fetcher fetches the data from the given url.
package fetcher

import (
	"context"
	"io"
	"log"
	"net/http"
	"strings"
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

func Fetch(url string, repo *WebsiteRepository, ctx context.Context) ([]string, error) {
	baseURL, path := splitURL(url)

	if checkPathDisallowed(baseURL, path) {
		return []string{}, nil
	}

	resp, err := http.Get(url)
	if err != nil {
		log.Println("Error fetching url: ", url, " Error: ", err)
		return []string{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return []string{}, nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading body: ", err)
		return []string{}, err
	}

	outboundLinks := getOutboundLinks(string(body), baseURL)

	log.Println("Crawled: ", url)
	_, err = repo.Create(ctx, url, string(body))
	if err != nil {
		return []string{}, err
	}

	return outboundLinks, nil
}
