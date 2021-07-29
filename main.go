package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
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
	URL := os.Args[1]
	Parser(URL)
}

func Parser(url string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
	}
	LinkReader(resp)

}

func LinkReader(resp *http.Response) Link {
	page := html.NewTokenizer(resp.Body)
	link := Link{}
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
					text = page.Token().Data
					fmt.Println(text)
					for i := range token.Attr {
						if token.Attr[i].Key == "href" {
							url = strings.TrimSpace(token.Attr[i].Val)
							fmt.Println(url)
						}
					}
				}
			}
		}
	}
	return link
}
