package parser

import (
	"fmt"
	"golang.org/x/net/html"
	"regexp"
	"strings"
)

type Chapter struct {
	Number int
	Title  string
	Href   string
	Time   string
}

func ParseHTML(data string, chapter *Chapter) error {
	doc, err := html.Parse(strings.NewReader(data))
	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}

	// Pass the chapter struct down to extract data into it
	findChapterDiv(doc, chapter)
	return nil
}

func findChapterDiv(n *html.Node, chapter *Chapter) {
	if n.Type == html.ElementNode && n.Data == "div" {
		for _, attr := range n.Attr {
			if attr.Key == "class" && strings.Contains(attr.Val, "l-chapter") {
				extractDataFromChapterDiv(n, chapter)
			}
		}
	}

	// Recursively search for the div in child nodes
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		findChapterDiv(c, chapter)
	}
}

// Function to extract data from the div with class "l-chapter"
func extractDataFromChapterDiv(n *html.Node, chapter *Chapter) {

	// Iterate through child nodes of the "l-chapter" div
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "div" {
			for _, attr := range c.Attr {
				if attr.Key == "class" && strings.Contains(attr.Val, "item") {
					extractDataFromItemDiv(c, chapter)
				}
			}
		}
	}

}

// Function to extract data from the div with class "item"
func extractDataFromItemDiv(n *html.Node, chapter *Chapter) {
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

	// Extract the chapter number from the title using regex
	re := regexp.MustCompile(`Chapter\s+(\d+)`)
	match := re.FindStringSubmatch(chapterTitle)
	if len(match) > 1 {
		// Convert the chapter number to an integer
		fmt.Sscanf(match[1], "%d", &chapter.Number)
	}

	// Assign the extracted data to the chapter struct
	chapter.Title = chapterTitle
	chapter.Href = href
	chapter.Time = itemTime
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
