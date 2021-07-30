package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	text string
	url  string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please specify a starting page of the url")
		os.Exit(1)
	}
	//log.SetPriorityString("info")
	log.SetPrefix("crawler")
	URL := os.Args[1]
	Parser(URL)
}

func Parser(URL string) {

	//defer resp.Body.Close()

	pages := dfs(URL, 2)
	for _, page := range pages {
		fmt.Println(page)
	}
}

func LinkReader(resp *http.Response) []Link {
	page := html.NewTokenizer(resp.Body)
	link := []Link{}
	var url string
	var text string
	for {
		tokenType := page.Next()
		if tokenType == html.ErrorToken {
			err := page.Err()
			if err == io.EOF {
				break
			}
			log.Fatalf("error tokenizing HTML: %v", page.Err())
		}
		if tokenType == html.StartTagToken {
			token := page.Token()
			if token.Data == "a" {
				tokenType = page.Next()
				if tokenType == html.TextToken {
					text = strings.TrimSpace(page.Token().Data)
					for i := range token.Attr {
						if token.Attr[i].Key == "href" {
							url = strings.TrimSpace(token.Attr[i].Val)
							//fmt.Println(link)
							link = append(link, Link{text: text, url: url})
						}
					}
				}
			}
		}
	}
	return link
}

func hrefs(links []Link, base string) []string {
	var ret []string
	for _, v := range links {
		switch {
		case strings.HasPrefix(v.url, "/"):
			ret = append(ret, base+v.url)
		case strings.HasPrefix(v.url, "http"):
			ret = append(ret, v.url)
		}
	}
	return ret
}

func get(urlStr string) []string {
	resp, err := http.Get(urlStr)
	if err != nil {
		log.Fatal(err)
	}
	links := LinkReader(resp)
	reqUrl := resp.Request.URL
	fmt.Println(reqUrl)
	baseUrl := &url.URL{
		Scheme: reqUrl.Scheme,
		Host:   reqUrl.Host,
	}
	base := baseUrl.String()
	return filter(base, hrefs(links, base))
}

func filter(base string, links []string) []string {
	var ret []string
	for _, v := range links {
		if strings.HasPrefix(v, base) {
			ret = append(ret, v)
		}
	}
	return ret
}

func dfs(urlStr string, maxDepth int) []string {
	seen := make(map[string]struct{})
	var q map[string]struct{}
	nq := map[string]struct{}{
		urlStr: {},
	}
	for i := 0; i <= maxDepth; i++ {
		q, nq = nq, make(map[string]struct{})
		for url, _ := range q {
			if _, ok := seen[url]; ok {
				continue
			}
			seen[url] = struct{}{}
			for _, link := range get(url) {
				nq[link] = struct{}{}
			}
		}
	}
	var ret []string
	for url, _ := range seen {
		ret = append(ret, url)
	}
	return ret
}
