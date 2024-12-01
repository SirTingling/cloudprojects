package main

import (
	"flag"
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"net/url"
	"os"
	"sync"
)

var visited = make(map[string]bool)
var mu sync.Mutex
var baseDomain string

func main() {
	// Define CLI flag for URL
	urlFlag := flag.String("url", "", "The URL to start scraping from")
	flag.Parse()

	if *urlFlag == "" {
		fmt.Println("Usage: ./webscraper -url <url>")
		os.Exit(1)
	}

	startURL := *urlFlag
	parsedURL, err := url.Parse(startURL)
	if err != nil {
		fmt.Printf("Invalid URL: %v\n", err)
		os.Exit(1)
	}
	baseDomain = parsedURL.Host

	fmt.Printf("Checking %s for dead links\n", startURL)

	queue := make(chan string)
	var wg sync.WaitGroup

	// Start processing URLs
	wg.Add(1)
	go func() {
		defer wg.Done()
		queue <- startURL
	}()

	// Process the queue of URLs
	go func() {
		wg.Wait()
		close(queue)
	}()

	for link := range queue {
		wg.Add(1)
		go func(link string) {
			defer wg.Done()
			processLink(link, queue)
		}(link)
	}
}

func processLink(link string, queue chan string) {
	// Lock the visited map to prevent race conditions
	mu.Lock()
	if visited[link] {
		mu.Unlock()
		return
	}
	visited[link] = true
	mu.Unlock()

	// Make the HTTP request without following redirects
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Prevent any redirects
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Get(link)
	if err != nil {
		fmt.Printf("Error accessing %s: %v\n", link, err)
		return
	}
	defer resp.Body.Close()

	// Check if the response is a redirect
	if resp.StatusCode >= 300 && resp.StatusCode < 400 {
		fmt.Printf("Redirect found (not followed): %s (Status Code: %d)\n", link, resp.StatusCode)
		return
	}

	// Check if the link is a dead link
	if resp.StatusCode >= 400 {
		fmt.Printf("Dead link found: %s (Status Code: %d)\n", link, resp.StatusCode)
		return
	}

	// Parse the HTML document
	doc, err := html.Parse(resp.Body)
	if err != nil {
		fmt.Printf("Error parsing HTML at %s: %v\n", link, err)
		return
	}

	// Extract links from the page and add them to the queue
	baseURL, _ := url.Parse(link)
	links := extractLinks(doc, baseURL)
	for _, l := range links {
		parsedLink, err := url.Parse(l)
		if err != nil {
			continue
		}
		// Only process links that belong to the base domain
		if parsedLink.Host == "" || parsedLink.Host == baseDomain {
			fullURL := parsedLink.String()
			mu.Lock()
			if !visited[fullURL] {
				visited[fullURL] = true
				mu.Unlock()
				queue <- fullURL
			} else {
				mu.Unlock()
			}
		}
	}
}

func extractLinks(doc *html.Node, baseURL *url.URL) []string {
	var links []string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					link, err := baseURL.Parse(attr.Val)
					if err == nil && (link.Host == "" || link.Host == baseDomain) {
						links = append(links, link.String())
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return links
}
