package parser

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"strings"

	"golang.org/x/net/html"
)

// ParsedText represents extracted text from a source.
type ParsedText struct {
	Title   string
	Content string
	Lang    string
}

// Parse auto-detects format from contentType and explicit hint, then extracts text.
func Parse(data []byte, contentType string, format string) (*ParsedText, error) {
	f := strings.ToLower(format)
	if f == "" {
		ct := strings.ToLower(contentType)
		switch {
		case strings.Contains(ct, "html"):
			f = "html"
		case strings.Contains(ct, "xml"):
			f = "xml"
		case strings.Contains(ct, "json"):
			f = "json"
		default:
			trimmed := bytes.TrimSpace(data)
			if len(trimmed) > 0 {
				switch trimmed[0] {
				case '<':
					prefix := trimmed[:clamp(50, len(trimmed))]
					if bytes.Contains(bytes.ToLower(prefix), []byte("html")) {
						f = "html"
					} else {
						f = "xml"
					}
				case '{', '[':
					f = "json"
				default:
					f = "txt"
				}
			} else {
				f = "txt"
			}
		}
	}

	switch f {
	case "html":
		return parseHTML(data)
	case "xml":
		return parseXML(data)
	case "json":
		return parseJSON(data)
	default:
		return parseTXT(data)
	}
}

func clamp(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// skipTags lists HTML elements whose content should be suppressed.
var skipTags = map[string]bool{
	"script": true, "style": true, "nav": true,
	"header": true, "footer": true, "noscript": true,
}

// parseHTML strips tags and returns visible text.
func parseHTML(data []byte) (*ParsedText, error) {
	doc, err := html.Parse(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("html parse: %w", err)
	}

	var title string
	var buf strings.Builder
	var skipDepth int

	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode {
			tag := strings.ToLower(n.Data)
			if tag == "title" && title == "" {
				if n.FirstChild != nil && n.FirstChild.Type == html.TextNode {
					title = strings.TrimSpace(n.FirstChild.Data)
				}
			}
			if skipTags[tag] {
				skipDepth++
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					walk(c)
				}
				skipDepth--
				return
			}
			if isBlock(tag) {
				buf.WriteString("\n")
			}
		}
		if n.Type == html.TextNode && skipDepth == 0 {
			text := strings.TrimSpace(n.Data)
			if text != "" {
				buf.WriteString(text)
				buf.WriteString(" ")
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
		if n.Type == html.ElementNode {
			tag := strings.ToLower(n.Data)
			if isBlock(tag) {
				buf.WriteString("\n")
			}
		}
	}
	walk(doc)

	return &ParsedText{
		Title:   title,
		Content: strings.TrimSpace(buf.String()),
	}, nil
}

func isBlock(tag string) bool {
	switch tag {
	case "p", "div", "br", "h1", "h2", "h3", "h4", "h5", "h6",
		"li", "tr", "td", "th", "blockquote", "pre", "article", "section":
		return true
	}
	return false
}

// parseTXT handles plain text.
func parseTXT(data []byte) (*ParsedText, error) {
	return &ParsedText{Content: string(data)}, nil
}

// xmlNode is used for generic XML parsing.
type xmlNode struct {
	XMLName  xml.Name
	Attrs    []xml.Attr `xml:",any,attr"`
	Content  string     `xml:",chardata"`
	Children []xmlNode  `xml:",any"`
}

// parseXML handles XML (OSIS, USFM, and other Bible XML formats).
func parseXML(data []byte) (*ParsedText, error) {
	var root xmlNode
	if err := xml.Unmarshal(data, &root); err != nil {
		// Try wrapping as fragment.
		wrapped := append([]byte("<root>"), append(data, []byte("</root>")...)...)
		if err2 := xml.Unmarshal(wrapped, &root); err2 != nil {
			return nil, fmt.Errorf("xml parse: %w", err)
		}
	}

	var buf strings.Builder
	var extractText func(xmlNode)
	extractText = func(n xmlNode) {
		text := strings.TrimSpace(n.Content)
		if text != "" {
			buf.WriteString(text)
			buf.WriteString("\n")
		}
		for _, child := range n.Children {
			extractText(child)
		}
	}
	extractText(root)

	return &ParsedText{Content: strings.TrimSpace(buf.String())}, nil
}

// parseJSON traverses a JSON value tree and collects all string leaves.
func parseJSON(data []byte) (*ParsedText, error) {
	dec := json.NewDecoder(bytes.NewReader(data))
	var buf strings.Builder

	var extract func(v interface{})
	extract = func(v interface{}) {
		switch val := v.(type) {
		case map[string]interface{}:
			for _, mv := range val {
				extract(mv)
			}
		case []interface{}:
			for _, item := range val {
				extract(item)
			}
		case string:
			if t := strings.TrimSpace(val); t != "" {
				buf.WriteString(t)
				buf.WriteString("\n")
			}
		}
	}

	var obj interface{}
	if err := dec.Decode(&obj); err != nil && err != io.EOF {
		return nil, fmt.Errorf("json parse: %w", err)
	}
	if obj != nil {
		extract(obj)
	}

	return &ParsedText{Content: strings.TrimSpace(buf.String())}, nil
}
