package parser

import (
	"io"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

func ExtractLinks(base *url.URL, body io.Reader) ([]string, error) {
	tokenizer := html.NewTokenizer(body)
	var links []string

	for {
		switch tokenizer.Next() {
		case html.ErrorToken:
			if tokenizer.Err() == io.EOF {
				return links, nil
			}
			return links, tokenizer.Err()

		case html.StartTagToken, html.SelfClosingTagToken:
			token := tokenizer.Token()
			if token.Data != "a" {
				continue
			}
			for _, attr := range token.Attr {
				if attr.Key != "href" {
					continue
				}
				if abs := resolve(base, attr.Val); abs != "" {
					links = append(links, abs)
				}
			}
		}
	}
}

func resolve(base *url.URL, href string) string {
	href = strings.TrimSpace(href)
	if href == "" {
		return ""
	}

	ref, err := url.Parse(href)
	if err != nil {
		return ""
	}

	abs := base.ResolveReference(ref)
	if abs.Scheme != "http" && abs.Scheme != "https" {
		return ""
	}

	abs.Fragment = ""
	return abs.String()
}
