package parser

import (
	"fmt"
	"golang.org/x/net/html"
	"regexp"
	"strings"
)

func ParseHTML(data string) {
	doc, err := html.Parse(strings.NewReader(data))
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	findChapterDiv(doc)
}

func findChapterDiv(n *html.Node) {
	if n.Type == html.ElementNode && n.Data == "div" {
		for _, attr := range n.Attr {
			if attr.Key == "class" && strings.Contains(attr.Val, "l-chapter") {
				extractDataFromChapterDiv(n)
			}
		}
	}

	// Recursively search for the div in child nodes
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		findChapterDiv(c)
	}
}

// Function to extract data from the div with class "l-chapter"
func extractDataFromChapterDiv(n *html.Node) {

	// Iterate through child nodes of the "l-chapter" div
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "div" {
			for _, attr := range c.Attr {
				if attr.Key == "class" && strings.Contains(attr.Val, "item") {
					extractDataFromItemDiv(c)
				}
			}
		}
	}

}

// Function to extract data from the div with class "item"
func extractDataFromItemDiv(n *html.Node) {
	var chapterTitle, href, itemTime string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "div" {
			for _, attr := range c.Attr {
				if attr.Key == "class" && attr.Val == "item-value" {
					chapterTitle, href, _ = extractLinkData(c)
				} else if attr.Key == "class" && attr.Val == "item-time" {
					itemTime = extractTextData(c)
				}
			}
		}
	}
	re := regexp.MustCompile(`Chapter\s+(\d+)`)

	// Step 2: Find the number after "Chapter"
	match := re.FindStringSubmatch(chapterTitle)

	fmt.Println("Chapter Number: ", match[1])
	fmt.Printf("Chapter Title: %s\n", chapterTitle)
	fmt.Printf("href: %s\n", href)
	fmt.Printf("Item Time: %s\n", itemTime)
}

// Helper function to extract data from the <a> tag within "item-value" div
func extractLinkData(n *html.Node) (chapterTitle, href, title string) {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "a" {
			for _, attr := range c.Attr {
				if attr.Key == "href" {
					href = attr.Val
				} else if attr.Key == "title" {
					title = attr.Val
				}
			}
			chapterTitle = extractTextData(c)
		}
	}
	return
}

// Helper function to extract text content from a node
func extractTextData(n *html.Node) string {
	if n.Type == html.TextNode {
		return strings.TrimSpace(n.Data)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		text := extractTextData(c)
		if text != "" {
			return text
		}
	}
	return ""
}
