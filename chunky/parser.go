package chunky

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

type Parser struct {
	html string
}

func NewParser(html string) *Parser {
	return &Parser{
		html: html,
	}
}

// ExtractContentBlocks extracts and returns top-level content blocks from the HTML.
func (parser *Parser) ExtractContentBlocks() ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(parser.html))
	if err != nil {
		return nil, err
	}

	var contentBlocks []string
	processedNodes := make(map[*html.Node]bool)

	doc.Find("body").Children().Each(func(i int, s *goquery.Selection) {
		if parser.isTopLevelContent(s) && !parser.isNodeProcessed(s.Nodes[0], processedNodes) {
			htmlStr, err := parser.goquerySelectionToString(s)
			if err == nil {
				contentBlocks = append(contentBlocks, htmlStr)
				parser.markNodeAsProcessed(s.Nodes[0], processedNodes)
			}
		}
	})

	return contentBlocks, nil
}

// isTopLevelContent determines if a selection is a top-level content block.
func (parser *Parser) isTopLevelContent(s *goquery.Selection) bool {
	topLevelTags := []string{"div", "section", "article", "main"}
	for _, tag := range topLevelTags {
		if goquery.NodeName(s) == tag {
			return true
		}
	}
	return false
}

// isNodeProcessed checks if a node or its ancestors are already processed.
func (parser *Parser) isNodeProcessed(node *html.Node, processedNodes map[*html.Node]bool) bool {
	for n := node; n != nil; n = n.Parent {
		if processedNodes[n] {
			return true
		}
	}
	return false
}

// markNodeAsProcessed marks a node and all its children as processed.
func (parser *Parser) markNodeAsProcessed(node *html.Node, processedNodes map[*html.Node]bool) {
	var mark func(*html.Node)
	mark = func(n *html.Node) {
		processedNodes[n] = true
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			mark(c)
		}
	}
	mark(node)
}

// goquerySelectionToString converts a goquery selection to its HTML string representation.
func (parser *Parser) goquerySelectionToString(s *goquery.Selection) (string, error) {
	var buf bytes.Buffer
	for _, node := range s.Nodes {
		if err := html.Render(&buf, node); err != nil {
			return "", err
		}
	}
	return buf.String(), nil
}

// SimplifyContent simplifies the HTML content while preserving essential structure.
func (parser *Parser) SimplifyContent(htmlContent string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		log.Println("Error parsing HTML:", err)
		return ""
	}

	var simplifiedContent strings.Builder
	doc.Find("*").Each(func(i int, s *goquery.Selection) {
		if goquery.NodeName(s) == "#text" {
			simplifiedContent.WriteString(s.Text() + " ")
		} else {
			simplifiedContent.WriteString(fmt.Sprintf("<%s> ", goquery.NodeName(s)))
		}
	})

	return simplifiedContent.String()
}

