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

// ParsedText represents extracted text from a single source.
type ParsedText struct {
	Title   string
	Content string
	Lang    string
}

// ParsedBook represents a single book extracted from a multi-book source (like a full Bible JSON).
type ParsedBook struct {
	Number  int
	Title   string
	Content string
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
	// First, try parsing as Bible SuperSearch API format.
	if text, err := parseSuperSearchJSON(data); err == nil && text.Content != "" {
		return text, nil
	}

	// Try parsing as Bible SuperSearch Exported format (manual download).
	// Since this format returns multiple books, we'll just join them for the generic parser.
	if books, _, err := ParseExportedJSON(data); err == nil && len(books) > 0 {
		var buf strings.Builder
		for _, b := range books {
			buf.WriteString(b.Content)
			buf.WriteString("\n\n")
		}
		return &ParsedText{Content: strings.TrimSpace(buf.String())}, nil
	}

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

// ExportedJSON represents the format of manually downloaded Bible SuperSearch files.
type ExportedJSON struct {
	Metadata struct {
		Name      string `json:"name"`
		Shortname string `json:"shortname"`
		Module    string `json:"module"`
		Year      string `json:"year"`
		Language  string `json:"lang"`
	} `json:"metadata"`
	Verses []struct {
		BookName string `json:"book_name"`
		Book     int    `json:"book"`
		Chapter  int    `json:"chapter"`
		Verse    int    `json:"verse"`
		Text     string `json:"text"`
	} `json:"verses"`
}

// ParseExportedJSON parses the "exported" JSON format from manual downloads.
// Returns a slice of ParsedBooks and the language code.
func ParseExportedJSON(data []byte) ([]ParsedBook, string, error) {
	var exported ExportedJSON
	if err := json.Unmarshal(data, &exported); err != nil {
		return nil, "", err
	}

	if len(exported.Verses) == 0 {
		return nil, "", fmt.Errorf("no verses in exported json")
	}

	var books []ParsedBook
	var currentBook *ParsedBook
	var buf strings.Builder
	var lastChapter int
	var lastBookNum int

	for _, v := range exported.Verses {
		if currentBook == nil || v.Book != lastBookNum {
			if currentBook != nil {
				currentBook.Content = strings.TrimSpace(buf.String())
				books = append(books, *currentBook)
			}
			buf.Reset()
			currentBook = &ParsedBook{
				Number: v.Book,
				Title:  v.BookName,
			}
			lastBookNum = v.Book
			lastChapter = v.Chapter
		} else if v.Chapter != lastChapter {
			if buf.Len() > 0 {
				buf.WriteString("\n\n")
			}
			lastChapter = v.Chapter
		} else if buf.Len() > 0 {
			buf.WriteString(" ")
		}
		buf.WriteString(strings.TrimSpace(v.Text))
	}

	if currentBook != nil {
		currentBook.Content = strings.TrimSpace(buf.String())
		books = append(books, *currentBook)
	}

	return books, exported.Metadata.Language, nil
}

// ParseBibleXML parses the common XML format found in many Bible repositories.
// Structure: <bible><testament><book number="X"><chapter number="Y"><verse number="Z">Text</verse>...
func ParseBibleXML(data []byte) ([]ParsedBook, string, error) {
	type XMLVerse struct {
		Number int    `xml:"number,attr"`
		Text   string `xml:",chardata"`
	}
	type XMLChapter struct {
		Number int        `xml:"number,attr"`
		Verses []XMLVerse `xml:"verse"`
	}
	type XMLBook struct {
		Number   int          `xml:"number,attr"`
		Name     string       `xml:"name,attr"`
		Chapters []XMLChapter `xml:"chapter"`
	}
	type XMLTestament struct {
		Books []XMLBook `xml:"book"`
	}
	type XMLBible struct {
		Translation string         `xml:"translation,attr"`
		Testaments  []XMLTestament `xml:"testament"`
		Books       []XMLBook      `xml:"book"`
	}

	var bible XMLBible
	if err := xml.Unmarshal(data, &bible); err != nil {
		return nil, "", err
	}

	var books []ParsedBook
	processBooks := func(xmlBooks []XMLBook) {
		for _, b := range xmlBooks {
			var buf strings.Builder
			for _, c := range b.Chapters {
				for _, v := range c.Verses {
					if buf.Len() > 0 {
						if v.Number == 1 && buf.String()[buf.Len()-1] != '\n' {
							buf.WriteString("\n\n")
						} else if buf.String()[buf.Len()-1] != '\n' {
							buf.WriteString(" ")
						}
					}
					buf.WriteString(strings.TrimSpace(v.Text))
				}
				buf.WriteString("\n\n")
			}
			books = append(books, ParsedBook{
				Number:  b.Number,
				Title:   b.Name,
				Content: strings.TrimSpace(buf.String()),
			})
		}
	}

	for _, t := range bible.Testaments {
		processBooks(t.Books)
	}
	processBooks(bible.Books)

	return books, "", nil
}

// parseSuperSearchJSON handles the specific JSON format from api.biblesupersearch.com.
func parseSuperSearchJSON(data []byte) (*ParsedText, error) {
	var resp struct {
		Results []struct {
			Verses map[string]map[string]map[string]struct {
				Text string `json:"text"`
			} `json:"verses"`
		} `json:"results"`
	}

	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	if len(resp.Results) == 0 {
		return nil, fmt.Errorf("no results in bible supersearch json")
	}

	var buf strings.Builder
	// The API returns chapters and verses as keys. We need to sort them or at least
	// iterate consistently. Since they are strings in JSON, we'll try to reconstruct.
	for _, result := range resp.Results {
		for _, chapters := range result.Verses {
			// We need a stable order for chapters and verses.
			// The keys are "1", "2", "3" etc.
			for c := 1; c <= 200; c++ {
				chapterKey := fmt.Sprintf("%d", c)
				verses, ok := chapters[chapterKey]
				if !ok {
					continue
				}
				for v := 1; v <= 200; v++ {
					verseKey := fmt.Sprintf("%d", v)
					verseData, ok := verses[verseKey]
					if !ok {
						continue
					}
					if buf.Len() > 0 {
						buf.WriteString(" ")
					}
					buf.WriteString(strings.TrimSpace(verseData.Text))
				}
				buf.WriteString("\n\n")
			}
		}
	}

	content := strings.TrimSpace(buf.String())
	if content == "" {
		return nil, fmt.Errorf("extracted content is empty")
	}

	return &ParsedText{Content: content}, nil
}
