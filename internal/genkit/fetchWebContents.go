package genkit

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// fetchWebContent fetches and processes the content from the provided URL
func fetchWebContent(url string) (string, error) {
	// Fetch the content from the provided URL
	res, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer res.Body.Close()

	// Read the HTML content
	html, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// Load the HTML content into goquery for parsing
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(html))
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML: %w", err)
	}

	// Remove unnecessary elements
	doc.Find("script, style, noscript").Each(func(i int, s *goquery.Selection) {
		s.Remove()
	})

	// Prefer 'article' content, fallback to 'body' if not available
	article := doc.Find("article").Text()
	if article != "" {
		return article, nil
	}

	body := doc.Find("body").Text()

	var result []string
	scanner := bufio.NewScanner(strings.NewReader(body))

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			result = append(result, line)
		}
	}

	return strings.Join(result, "\n"), nil
}
