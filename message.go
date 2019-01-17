package main

import (
	"bytes"
	"strings"

	"golang.org/x/net/html"
)

func ExtractMessage(msg string) (string, error) {
	body, err := html.Parse(strings.NewReader(msg))
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := extractMessage(body, &buf); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func extractMessage(node *html.Node, buf *bytes.Buffer) error {
	if node.Data == "p" && buf.Len() > 0 {
		if _, err := buf.WriteString("\n"); err != nil {
			return err
		}
	}

	switch node.Type {
	case html.TextNode:
		if len(node.Data) > 0 {
			if _, err := buf.WriteString(node.Data); err != nil {
				return err
			}
		}
	case html.ElementNode:
		if node.Data == "br" {
			if _, err := buf.WriteString("\n"); err != nil {
				return err
			}
		}
	default:
		// do nothing
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if err := extractMessage(child, buf); err != nil {
			return err
		}
	}

	return nil
}
