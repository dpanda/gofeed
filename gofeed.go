package main

import (
	"fmt"

	"github.com/mmcdole/gofeed"
	"github.com/gorilla/feeds"
	"log"
	"context"
	"github.com/aws/aws-lambda-go/events"
	"golang.org/x/net/html"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(HandleRequest)
	//parseFeed("awkardyeti")
	//parseFeed("dilbert")
	//parseFeed("nedroid")
	//parseFeed("stefanotartarotti")
}

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	rss, err := parseFeed(request.QueryStringParameters["feed"])
	headers := make(map[string]string)
	headers["Content-Type"] = "application/xml"
	return events.APIGatewayProxyResponse{Body: rss, Headers: headers, StatusCode: 200}, err
}

func parseFeed(feedName string) (string, error) {

	var link string
	switch feedName {
	case "dilbert":
		link = "http://www.dilbert.com/feed"
	case "awkardyeti":
		link = "http://theawkwardyeti.com/feed/"
	case "nedroid":
		link = "http://nedroid.com/feed/"
	case "stefanotartarotti":
		link = "https://www.ilpost.it/stefanotartarotti/feed/"
	}

	log.Println("Starting")
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL(link)

	fmt.Println("Parsing feed", feed.Title, feed.Link)

	outfeed := &feeds.Feed{
		Title:       feed.Title,
		Link:        &feeds.Link{Href: feed.Link},
		Description: feed.Description,
		Created:     feed.UpdatedParsed.Add(0),
		Updated:     feed.UpdatedParsed.Add(0),
	}

	c := convertFeed(feedName, feed)

	for i := 0; i < len(feed.Items); i++ {
		outitem := <-c
		log.Println("done item ", outitem.Title)
		outfeed.Items = append(outfeed.Items, &outitem)
	}

	rss, err := outfeed.ToRss()
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	log.Println("Done, returning xml")
	log.Println(rss)
	return rss, nil
}

func convertFeed(id string, feed *gofeed.Feed) <-chan feeds.Item {
	c := make(chan feeds.Item)

	for i := 0; i < len(feed.Items); i++ {
		item := feed.Items[i]
		switch id {
		case "dilbert":
			go func() { c <- ConvertDilbertFeedItem(item) }()
		case "awkardyeti":
			go func() { c <- ConvertAwkardyetiFeedItem(item) }()
		case "nedroid":
			go func() { c <- ConvertNedroidFeedItem(item) }()
		case "stefanotartarotti":
			go func() { c <- ConvertStefanoTartarottiFeedItem(item) }()
		}
	}
	return c
}

func GetAttr(token html.Token, attr string) string {
	for _, a := range token.Attr {
		if a.Key == attr {
			return a.Val
		}
	}
	return ""
}
