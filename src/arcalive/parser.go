package arcalive

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Post struct {
	PostID    int       `json:"postId"`
	Category  string    `json:"category"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"createdAt"`
}

type ParseResult struct {
	Posts       []Post `json:"posts"`
	NextPageURL string `json:"nextPageUrl"`
}

func GetPosts(baseURL string, page int, before string) (ParseResult, error) {
	var fetchURL string
	if before != "" {
		fetchURL = fmt.Sprintf("%s?before=%s", baseURL, before)
	} else {
		fetchURL = fmt.Sprintf("%s?p=%d", baseURL, page)
	}

	resp, err := Fetch(FetcherProps{
		URL:    fetchURL,
		Method: "GET",
	})
	if err != nil {
		return ParseResult{}, fmt.Errorf("failed to fetch URL: %v", err)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return ParseResult{}, fmt.Errorf("failed to parse HTML: %v", err)
	}
	defer resp.Body.Close()

	var result ParseResult

	doc.Find(".pagination-wrapper .pagination .page-item.active").Next().Find("a.page-link").Each(func(_ int, s *goquery.Selection) {
		if href, exists := s.Attr("href"); exists {
			result.NextPageURL = "https://arca.live" + href
		}
	})

	doc.Find(".list-table.table .vrow.column").Each(func(_ int, s *goquery.Selection) {
		if s.HasClass("notice") {
			return
		}

		href, exists := s.Attr("href")
		if !exists {
			return
		}

		postIDStr := extractPostID(href)
		postID, err := strconv.Atoi(postIDStr)
		if err != nil {
			return
		}

		category := s.Find(".badge").First().Text()
		category = strings.TrimSpace(category)

		titleSel := s.Find(".vcol.col-title .title")
		titleSel.Children().Remove()
		title := strings.TrimSpace(titleSel.Text())

		author, exists := s.Find(".user-info [data-filter]").Attr("data-filter")
		if !exists {
			return
		}

		createdAtStr, exists := s.Find("time").Attr("datetime")
		if !exists {
			return
		}

		createdAt, err := time.Parse(time.RFC3339, createdAtStr)
		if err != nil {
			return
		}

		post := Post{
			PostID:    postID,
			Category:  category,
			Title:     title,
			Author:    author,
			CreatedAt: createdAt,
		}

		result.Posts = append(result.Posts, post)
	})

	return result, nil
}

func extractPostID(href string) string {
	parts := strings.Split(href, "/")
	if len(parts) < 4 {
		return ""
	}

	lastPart := parts[len(parts)-1]
	postID := strings.Split(lastPart, "?")[0]

	return postID
}
