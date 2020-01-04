package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"golang.org/x/net/html"
)

// Interfaces with user
func main() {
	input := os.Args[1]
	start, err := url.Parse(input)
	if err != nil {
		panic(err)
	}
	page, err := CrawlPage(start, PageReader, LinkExtractor)
	if err != nil {
		panic(err)
	}
	for _, link := range page.Links {
		fmt.Println(link)
	}
}

// Page represents a single HTML page
type Page struct {
	Contents []byte
	Links    []*url.URL
}

// PageReader reads a web page.
func PageReader(u *url.URL) (io.ReadCloser, error) {
	// fetch url
	resp, err := http.DefaultClient.Get(u.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get page %s: %v", u, err)
	}
	return resp.Body, nil
}

// LinkExtractor extracts urls from from html contents
func LinkExtractor(page io.Reader) []*url.URL {
	links := make([]*url.URL, 0)
	doc, err := html.Parse(page)
	if err != nil {
		fmt.Printf("error parsing page: %v\n", err)
		return links
	}
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					link, err := url.Parse(a.Val)
					if err != nil {
						fmt.Printf("failed to parse url %q: %v", a.Val, err)
						continue
					}
					links = append(links, link)
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

// CrawlPage will attempt to read the URL, extract links from it and return the page.
func CrawlPage(u *url.URL, pageReader func(*url.URL) (io.ReadCloser, error), linkExtractor func(io.Reader) []*url.URL) (*Page, error) {
	body, err := pageReader(u)

	contents, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %v", err)
	}
	if err := body.Close(); err != nil {
		return nil, fmt.Errorf("failed to close body: %v", err)
	}
	page := &Page{
		Contents: contents,
	}

	buf := bytes.NewReader(contents)
	links := linkExtractor(buf)
	page.Links = links
	return page, nil
}
