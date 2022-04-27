package main

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/gocolly/colly/v2"
)

func contains(elems []string, v string) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}
func ExtractLink(uri string) {
	// links := make([]string, 0)
	u := uriValid(uri)
	c := colly.NewCollector(
		colly.AllowedDomains(u[:]...),
	)
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))
		fmt.Println("\t\t", link)
		if strings.Contains(link, "=") && !contains(LinksParam, link) {
			LinksParam = append(LinksParam, link)
		}
		// links = append(links, link)
		// fmt.Printf("%s\n", strings.Repeat(" ", 5))
		c.Visit(link)
	})
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Usr-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")
		r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
		r.Headers.Set("Accept-Language", "en-US,en;q=0.8")
		r.Headers.Set("Accept-Encoding", "gzip")
		// fmt.Println(*r.Headers)
		fmt.Println(r.URL.String())
	})
	c.OnResponse(func(r *colly.Response) {
		// fmt.Println(string(r.Body))
	})

	c.Visit(uri)
	// return links
}

func uriValid(uri string) []string {
	ur, _ := url.Parse(uri)
	if strings.HasPrefix(ur.Host, "www.") {
		rep := strings.Replace(ur.Host, ur.Host[:4], "", 1)
		// fmt.Println(rep)
		return []string{rep, ur.Host}
	} else {
		return []string{ur.Host, "www." + ur.Host}
	}

}
