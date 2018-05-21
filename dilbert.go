package main

import (
	"github.com/mmcdole/gofeed"
	"github.com/gorilla/feeds"
	"log"
	"strings"
	"fmt"
	"net/http"
	"golang.org/x/net/html"
)

func ConvertDilbertFeedItem(item *gofeed.Item) (outitem feeds.Item) {
	// get the content to find out the link to the image
	log.Println("Converting ", item.Link)
	resp, _ := http.Get(item.Link)
	tokenizer := html.NewTokenizer(resp.Body)

	var comicStripLink string
	var ok bool

	for {
		next := tokenizer.Next()
		switch {
		case next == html.SelfClosingTagToken:
			token := tokenizer.Token()
			isImage := token.Data == "img"
			if isImage {
				for _, a := range token.Attr {
					if a.Key == "src" && strings.HasPrefix(a.Val, "http://assets.amuniversal.com") {
						ok = true
						comicStripLink = a.Val
					}
				}
			}
		case next == html.ErrorToken:
			// end; we are done
			resp.Body.Close()
			if !ok {
				comicStripLink = item.Link
			}
			outitem = feeds.Item{
				Title:       item.Title,
				Link:        &feeds.Link{Href: comicStripLink},
				Description: item.Description,
				Content:     fmt.Sprintf("<img src=\"%s\">", comicStripLink),
				Created:     item.UpdatedParsed.Add(0),
				Updated:     item.UpdatedParsed.Add(0),
			}
			return
		}
	}
}
