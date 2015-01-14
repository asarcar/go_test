package main

import (
	"fmt"
)

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

type UrlReq struct {
	url   string
	depth int
}

type UrlRes struct {
	urls  []string
	depth int
}

func ParallelCrawl(urlReq UrlReq, fetcher Fetcher, ch chan UrlRes) {
	body, urls, err := fetcher.Fetch(urlReq.url)
	if err != nil {
		ch <- UrlRes{nil, urlReq.depth}
		fmt.Println(err)
		return
	}
	fmt.Printf("found: %s %q\n", urlReq.url, body)
	ch <- UrlRes{urls, urlReq.depth}
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(baseUrl string, depth int, fetcher Fetcher) {
	// 1. Fetch URLs in parallel.
	// 2. Don't fetch the same URL twice.

	// Keep tabs on all url that have been visited
	// Keep a count of every outstanding request to Parallel Crawl
	// Parallel Crawl should reply with either the new URLs or nil
	// (if error in fetching) for every request sent.
	// When the count of outstanding request counts down to zero we are done

	ch := make(chan UrlRes, 25)
	go ParallelCrawl(UrlReq{baseUrl, depth}, fetcher, ch)
	countToParallelCrawl := 1
	urlVisited := map[string]bool{baseUrl: true}

	for countToParallelCrawl > 0 {
		urlRes := <-ch
		countToParallelCrawl--
		for _, url := range urlRes.urls {
			// Terminate further crawl after nesting to full depth
			if urlRes.depth <= 0 {
				continue
			}
			// If URL already visited no need to crawl deeper on this one
			if _, alreadyVisited := urlVisited[url]; alreadyVisited {
				continue
			}

			go ParallelCrawl(UrlReq{url, urlRes.depth - 1}, fetcher, ch)
			countToParallelCrawl++
			urlVisited[url] = true
		}
	}
	return
}

func main() {
	Crawl("http://golang.org/", 4, fetcher)
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"http://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"http://golang.org/pkg/",
			"http://golang.org/cmd/",
		},
	},
	"http://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"http://golang.org/",
			"http://golang.org/cmd/",
			"http://golang.org/pkg/fmt/",
			"http://golang.org/pkg/os/",
		},
	},
	"http://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
	"http://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
}
