package main

import (
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

func removeDuplicatesUnordered(elements []string) []string {
	encountered := map[string]bool{}

	// Create a map of all unique elements.
	for v := range elements {
		encountered[elements[v]] = true
	}

	// Place all keys from the map into a slice.
	result := []string{}
	for key := range encountered {
		result = append(result, key)
	}
	return result
}

func cleanLink(link string) string {
	index := strings.Index(link, "#")
	if index > 0 {
		link = link[:index]
	}
	index = strings.Index(link, "?")
	if index > 0 {
		link = link[:index]
	}
	if link[len(link)-1:] == "/" {
		link = link[0 : len(link)-1]
	}
	if link[:5] == "https" {
		link = "http" + link[5:]
	}
	return link
}

// Helper function to pull the href attribute from a Token
func getHref(t html.Token) (ok bool, href string) {
	// Iterate over all of the Token's attributes until we find an "href"
	for _, a := range t.Attr {
		if a.Key == "href" {
			href = a.Val
			ok = true
		}
	}

	// "bare" return will return the variables (ok, href) as defined in
	// the function definition
	return
}

// Thanks to http://schier.co/blog/2015/04/26/a-simple-web-scraper-in-go.html
// Extract all http** links from a given webpage
func scrapeLinks(page string) (links []string) {
	resp, err := http.Get(page)

	if err != nil {
		fmt.Println("ERROR: Failed to crawl \"" + page + "\"")
		return
	}

	b := resp.Body
	defer b.Close() // close Body when the function returns

	z := html.NewTokenizer(b)

	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			// End of the document, we're done
			return
		case tt == html.StartTagToken:
			t := z.Token()

			// Check if the token is an <a> tag
			isAnchor := t.Data == "a"
			if !isAnchor {
				continue
			}

			// Extract the href value, if there is one
			ok, url := getHref(t)
			if !ok {
				continue
			}

			// Make sure the url begines in http**
			hasProto := strings.Index(url, "http") == 0
			if hasProto {
				links = append(links, cleanLink(url))
			} else {
				if strings.Index(url, "/") == 0 {
					parts := strings.Split(page, "/")
					links = append(links, cleanLink(parts[0]+"//"+parts[2]+url))
				}
			}
		}
		links = removeDuplicatesUnordered(links)
	}
}
