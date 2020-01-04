package main

import (
	"bytes"
	"fmt"
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
	page, err := CrawlPage(start)
	if err != nil {
		panic(err)
	}
	fmt.Println(page)
}

// Page represents a single HTML page
type Page struct {
	Contents []byte
	Links    []*url.URL
}

// CrawlPage will attempt to read the URL, extract links from it and return the page.
func CrawlPage(u *url.URL) (*Page, error) {
	// fetch url
	resp, err := http.DefaultClient.Get(u.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get page %s: %v", u, err)
	}

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read respons body: %v", err)
	}
	if err := resp.Body.Close(); err != nil {
		return nil, fmt.Errorf("failed to close body: %v", err)
	}
	page := &Page{
		Contents: contents,
		Links:    []*url.URL{},
	}

	buf := bytes.NewReader(contents)
	doc, err := html.Parse(buf)
	if err != nil {
		fmt.Println(err)
		return page, nil
	}
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					link, err := url.Parse(a.Val)
					if err != nil {
						fmt.Println(err)
						continue
					}
					page.Links = append(page.Links, link)
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return page, nil
}
