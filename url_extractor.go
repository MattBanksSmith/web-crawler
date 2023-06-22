package main

import (
	"bytes"
	"golang.org/x/net/html"
	"net/url"
)

func extractURLs(data []byte) (map[string]struct{}, error) {
	reader := bytes.NewReader(data)
	htmlData, err := html.Parse(reader)
	if err != nil {
		return nil, err
	}
	res := make(map[string]struct{})
	checkNode(htmlData, res)
	return res, nil
}

func checkNode(node *html.Node, urls map[string]struct{}) {
	if node.Type == html.ElementNode && node.Data == "a" {
		for _, attribute := range node.Attr {
			if attribute.Key == "href" {
				if parsed, err := url.Parse(attribute.Val); err != nil || !parsed.IsAbs() {
					continue
				}

				if _, ok := urls[attribute.Val]; !ok {
					urls[attribute.Val] = struct{}{}
				}
			}
		}
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		checkNode(child, urls)
	}
}
