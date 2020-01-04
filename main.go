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

// Wrap up the html.Parse function to encapsulate the dependency.
type htmlParser struct{}

func (h *htmlParser) Parse(r io.Reader) (*html.Node, error) {
	return html.Parse(r)
}

// Interfaces with user
func main() {
	input := os.Args[1]
	start, err := url.Parse(input)
	if err != nil {
		panic(err)
	}

	// Populate the real dependencies.
	pageReader := &PageReader{
		getter: http.DefaultClient,
	}
	linkExtractor := &LinkExtractor{
		parser: &htmlParser{},
	}

	page, err := CrawlPage(start, pageReader, linkExtractor)
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

type getter interface {
	Get(string) (*http.Response, error)
}

// PageReader is a struct that holds dependencies for reading pages.
type PageReader struct {
	getter getter
}

// ReadPage reads a web page.
func (p *PageReader) ReadPage(u *url.URL) (io.ReadCloser, error) {
	// fetch url
	resp, err := p.getter.Get(u.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get page %s: %v", u, err)
	}
	return resp.Body, nil
}

type parser interface {
	Parse(io.Reader) (*html.Node, error)
}

// LinkExtractor is a struct to hold dependencies for extracting links.
type LinkExtractor struct {
	parser parser
}

// ExtractLinks extracts urls from from html contents.
func (e *LinkExtractor) ExtractLinks(page io.Reader) []*url.URL {
	links := make([]*url.URL, 0)
	doc, err := e.parser.Parse(page)
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

// These shouldn't be reusable, they are designed specifically for CrawlPage

type pageReader interface {
	ReadPage(*url.URL) (io.ReadCloser, error)
}

type linkExtractor interface {
	ExtractLinks(io.Reader) []*url.URL
}

// CrawlPage will attempt to read the URL, extract links from it and return the page.
func CrawlPage(u *url.URL, pageReader pageReader, linkExtractor linkExtractor) (*Page, error) {
	body, err := pageReader.ReadPage(u)

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
	links := linkExtractor.ExtractLinks(buf)
	page.Links = links
	return page, nil
}
