package main

import (
	"github.com/mmcdole/gofeed"
	"github.com/gorilla/feeds"
	"fmt"
	"log"
	"golang.org/x/net/html"
	"net/http"
)

func ConvertNedroidFeedItem(item *gofeed.Item) (outitem feeds.Item) {
	// get the content to find out the link to the image
	log.Println("Converting ", item.Link)
	resp, _ := http.Get(item.Link)
	tokenizer := html.NewTokenizer(resp.Body)

	var comicStripLink string
	found := false
	nextToken := false

	for {
		next := tokenizer.Next()
		switch {
		case next == html.StartTagToken:
			token := tokenizer.Token()
			if token.Data == "div" && GetAttr(token, "id") == "comic" {
				log.Println("next")
				nextToken = true
			}
		case next == html.SelfClosingTagToken:
			token := tokenizer.Token()
			log.Println(token.Data, token.Attr)
			if nextToken && token.Data == "img" && !found {
				found = true
				comicStripLink = GetAttr(token, "src")
				log.Println("found", comicStripLink)
			}
		case next == html.ErrorToken:
			// end; we are done
			resp.Body.Close()
			if !found {
				comicStripLink = "not found" //item.Link
			}
			outitem = feeds.Item{
				Title: item.Title,
				Link:  &feeds.Link{Href: comicStripLink},
				//Description: item.Description,
				Content: fmt.Sprintf("<img src=\"%s\">", comicStripLink),
				Created: item.PublishedParsed.Add(0),
				Updated: item.PublishedParsed.Add(0),
			}
			return
		}
	}
}

