package marks

import (
	"slices"

	"golang.org/x/net/html"
)

func findTitle(n *html.Node) string {
	if n.Type == html.ElementNode && n.Data == "title" && n.FirstChild != nil {
		return n.FirstChild.Data
	}

	for node := n.FirstChild; node != nil; node = node.NextSibling {
		if title := findTitle(node); title != "" {
			return title
		}
	}

	return ""
}

func findMeta(n *html.Node, keys []string) string {
	if n.Type == html.ElementNode && n.Data == "meta" {
		var name, content string

		for _, attr := range n.Attr {
			if (attr.Key == "name" || attr.Key == "property") && slices.Contains(keys, attr.Val) {
				name = attr.Val
			}

			if attr.Key == "content" {
				content = attr.Val
			}
		}

		if slices.Contains(keys, name) {
			return content
		}
	}

	for node := n.FirstChild; node != nil; node = node.NextSibling {
		if desc := findMeta(node, keys); desc != "" {
			return desc
		}
	}

	return ""
}
