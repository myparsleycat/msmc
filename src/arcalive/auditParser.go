package arcalive

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type AuditAction string

const (
	AuditActionDelete AuditAction = "delete"
	AuditActionEdit   AuditAction = "edit"
	AuditActionBan    AuditAction = "ban"
	AuditActionOther  AuditAction = "other"
)

type Audit struct {
	AuditID   int         `json:"auditId"`
	CreatedAt time.Time   `json:"createdAt"`
	Author    string      `json:"author"`
	Action    AuditAction `json:"action"`
	PostID    int         `json:"postId,omitempty"`
	Detail    string      `json:"detail,omitempty"`
}

func GetAudits(baseURL string) ([]Audit, error) {
	resp, err := Fetch(FetcherProps{
		URL:    baseURL,
		Method: "GET",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %v", err)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %v", err)
	}
	defer resp.Body.Close()

	var audits []Audit

	doc.Find(".board-audit-list li").Each(func(_ int, s *goquery.Selection) {
		// 체크박스 스킵
		if s.Find(".batch-check-all").Length() > 0 {
			return
		}

		auditIDStr := strings.TrimPrefix(s.AttrOr("id", ""), "audit-")
		auditID, err := strconv.Atoi(auditIDStr)
		if err != nil {
			return
		}

		createdAtStr, exists := s.Find("time").First().Attr("datetime")
		if !exists {
			return
		}
		createdAt, err := time.Parse(time.RFC3339, createdAtStr)
		if err != nil {
			return
		}

		author := s.Find(".user-info a").First().AttrOr("data-filter", "")

		actionEl := s.Find("i")
		actionText := strings.TrimSpace(actionEl.Text())

		var action AuditAction
		switch {
		case strings.Contains(actionText, "삭제"):
			action = AuditActionDelete
		case strings.Contains(actionText, "편집"):
			action = AuditActionEdit
		case strings.Contains(actionText, "차단"):
			action = AuditActionBan
		default:
			action = AuditActionOther
		}

		var postID int
		postLink := actionEl.Find("a").FilterFunction(func(_ int, s *goquery.Selection) bool {
			href := s.AttrOr("href", "")
			return strings.Contains(href, "/b/")
		})
		if postLink.Length() > 0 {
			href := postLink.AttrOr("href", "")
			postIDStr := extractPostID(href)
			postID, _ = strconv.Atoi(postIDStr)
		}

		detail := ""
		if actionEl.Text() != "" {
			parts := strings.Split(actionEl.Text(), ")")
			if len(parts) > 1 {
				detail = strings.TrimSpace(strings.TrimPrefix(parts[len(parts)-1], "("))
				detail = strings.TrimSuffix(detail, ")")
			}
		}

		audit := Audit{
			AuditID:   auditID,
			CreatedAt: createdAt,
			Author:    author,
			Action:    action,
			PostID:    postID,
			Detail:    detail,
		}

		audits = append(audits, audit)
	})

	return audits, nil
}
