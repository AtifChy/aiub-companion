package fetcher

import (
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

var reTag = regexp.MustCompile(`(?i)<[^>]+>`)

// GetAttr retrieves the value of the specified attribute from the HTML node.
func GetAttr(n *html.Node, key string) string {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

// GetAttrInt retrieves the integer value of the specified attribute from the HTML node.
func GetAttrInt(n *html.Node, attr string, defaultValue int) int {
	for _, a := range n.Attr {
		if a.Key == attr {
			if val, err := strconv.Atoi(a.Val); err == nil {
				return val
			}
		}
	}
	return defaultValue
}

// GetInnerText retrieves the concatenated text content of the HTML node and its descendants.
func GetInnerText(n *html.Node) string {
	var sb strings.Builder
	var walk func(*html.Node)
	walk = func(node *html.Node) {
		if node.Type == html.TextNode {
			sb.WriteString(node.Data)
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(n)
	return strings.TrimSpace(sb.String())
}

// GetInnerHTML retrieves the HTML content of the HTML node and its descendants as a string.
func GetInnerHTML(n *html.Node) string {
	var sb strings.Builder
	var walk func(*html.Node)
	walk = func(node *html.Node) {
		switch node.Type {
		case html.TextNode:
			sb.WriteString(node.Data)
		case html.ElementNode:
			sb.WriteString("<")
			sb.WriteString(node.Data)
			for _, attr := range node.Attr {
				sb.WriteString(" ")
				sb.WriteString(attr.Key)
				sb.WriteString(`="`)
				sb.WriteString(attr.Val)
				sb.WriteString(`"`)
			}
			sb.WriteString(">")
			for child := node.FirstChild; child != nil; child = child.NextSibling {
				walk(child)
			}
			sb.WriteString("</")
			sb.WriteString(node.Data)
			sb.WriteString(">")
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		walk(c)
	}
	return sb.String()
}

// StripTags removes all HTML tags from the input string and returns the plain text.
func StripTags(s string) string {
	return reTag.ReplaceAllString(s, "")
}

// FindNodeByTag searches for the first HTML node with the specified tag name in the subtree rooted at the given node.
func FindNodeByTag(n *html.Node, tagName string) *html.Node {
	if n.Type == html.ElementNode && n.Data == tagName {
		return n
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if result := FindNodeByTag(c, tagName); result != nil {
			return result
		}
	}
	return nil
}

// FindNodeByID searches for the first HTML node with the specified ID in the subtree rooted at the given node.
func FindNodeByID(n *html.Node, id string) *html.Node {
	if n.Type == html.ElementNode && GetAttr(n, "id") == id {
		return n
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if result := FindNodeByID(c, id); result != nil {
			return result
		}
	}
	return nil
}

// FindNodesByClass searches for all HTML nodes with the specified tag name and class name in the subtree rooted at the given node.
func FindNodesByClass(n *html.Node, tagName, className string) []*html.Node {
	var result []*html.Node
	if n.Type == html.ElementNode && (tagName == "" || n.Data == tagName) {
		for _, attr := range n.Attr {
			if attr.Key != "class" {
				continue
			}
			for c := range strings.FieldsSeq(attr.Val) {
				if c == className {
					result = append(result, n)
					break
				}
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result = append(result, FindNodesByClass(c, tagName, className)...)
	}
	return result
}
