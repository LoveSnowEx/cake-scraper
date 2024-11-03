package htmlparser

import (
	"bytes"
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

func extractText(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	var buf bytes.Buffer
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		buf.WriteString(extractText(c))
	}
	return buf.String()
}

// Parse extracts text from an HTML string
func Parse(htmlString string) string {
	doc, err := html.Parse(strings.NewReader(htmlString))
	if err != nil {
		fmt.Println("Error parsing HTML:", err)
		return ""
	}
	return extractText(doc)
}
