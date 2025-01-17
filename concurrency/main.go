package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"golang.org/x/net/html"
)

type Crawler struct {
	visited   map[string]bool
	queue     chan string
	results   chan string
	wg        sync.WaitGroup
	mutex     sync.Mutex
	maxDepth  int
	rateLimit time.Duration
}

func NewCrawler(maxDepth int, rateLimit time.Duration) *Crawler {
	return &Crawler{
		visited:   make(map[string]bool),
		queue:     make(chan string),
		results:   make(chan string),
		maxDepth:  maxDepth,
		rateLimit: rateLimit,
	}
}

func (c *Crawler) Start(url string) {
	c.wg.Add(1)
	go c.crawl(url, 0)
	go c.processResults()
	c.wg.Wait()
	close(c.results)
}

func (c *Crawler) crawl(url string, depth int) {
	defer c.wg.Done()

	if depth > c.maxDepth {
		return
	}

	c.mutex.Lock()
	if c.visited[url] {
		c.mutex.Unlock()
		return
	}
	c.visited[url] = true
	c.mutex.Unlock()

	time.Sleep(c.rateLimit)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching URL:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: non-200 status code")
		return
	}

	c.results <- url

	doc, err := html.Parse(resp.Body)
	if err != nil {
		fmt.Println("Error parsing HTML:", err)
		return
	}

	var links []string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					links = append(links, a.Val)
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	for _, link := range links {
		c.wg.Add(1)
		go c.crawl(link, depth+1)
	}
}

func (c *Crawler) processResults() {
	for url := range c.results {
		fmt.Println("Crawled URL:", url)
	}
}

func main() {
	crawler := NewCrawler(2, 500*time.Millisecond)
	crawler.Start("https://webscraper.io/test-sites/e-commerce/allinone")
}
