package util

import (
	"bytes"
	"log"
	"strings"

	"golang.org/x/net/html"
)

func extract(node *html.Node, buff *bytes.Buffer) {
	if node.Type == html.TextNode {
		data := strings.Trim(node.Data, "\r\n ")
		if data != "" {
			buff.WriteString("\n")
			buff.WriteString(data)
		}
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		extract(c, buff)
	}
}

//Text Text
func Text(data []byte) string {
	reader := bytes.NewReader(data)
	var buffer bytes.Buffer
	doc, err := html.Parse(reader)
	if err != nil {
		log.Fatal(err)
	}
	extract(doc, &buffer)
	return buffer.String()
}
