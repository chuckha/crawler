Problems with this version:

The biggest thing that jumps out at me is that the `CrawlPage` function does
more than one thing. It first fetches a URL and then parses the content looking
for links.

The overall idea is that we want to take a web page and return all the contents
and all of the links within that page. One thing should be responsible for
getting the contents of the web page. Another separate thing should be
responsible for getting the links within the page. Here we are building
low-level abstractions to support some higher level goals. Our higher level goal
right now is to write a single web page scraper. By the end we will use this
high-level module to build an even higher level module and have a full blow web
crawler written using the SOLID design principles in mind.