package preview

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type PreviewInfo struct {
	Title       string
	Description string
	Image       string
	URL         string
}

// FetchOGTags 爬取网页并提取OG标签
func FetchOGTagss(url string) (*PreviewInfo, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if !strings.Contains(resp.Header.Get("Content-Type"), "text/html") {
		return nil, errors.New("不是HTML页面")
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	preview := &PreviewInfo{URL: url}

	doc.Find("meta").Each(func(i int, s *goquery.Selection) {
		prop, _ := s.Attr("property")
		content, _ := s.Attr("content")

		switch prop {
		case "og:title":
			preview.Title = content
		case "og:description":
			preview.Description = content
		case "og:image":
			preview.Image = content
		}
	})

	// fallback title
	if preview.Title == "" {
		preview.Title = doc.Find("title").Text()
	}

	return preview, nil
}

type OGDataa struct {
	Title          string
	Description    string
	Image          string
	URL            string
	CanonicalURL   string
	TwitterTitle   string
	TwitterDesc    string
	TwitterImage   string
	PageTitle      string
	ContentSnippet string   // 页面正文简要
	ImageList      []string // 所有图片地址
}

func FetchOGTagsA(targetURL string) (*OGDataa, error) {
	resp, err := http.Get(targetURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("请求失败: %s", resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	getMeta := func(attr, key string) string {
		content := ""
		doc.Find("meta").Each(func(i int, s *goquery.Selection) {
			if val, exists := s.Attr(attr); exists && val == key {
				if c, ok := s.Attr("content"); ok {
					content = c
				}
			}
		})
		return content
	}

	// 获取 <link rel="canonical">
	canonical := ""
	doc.Find("link[rel='canonical']").Each(func(i int, s *goquery.Selection) {
		if href, ok := s.Attr("href"); ok {
			canonical = href
		}
	})

	// 获取页面 <title>
	pageTitle := doc.Find("title").Text()

	// 获取正文前 1~2 段内容
	snippet := ""
	doc.Find("article, p").EachWithBreak(func(i int, s *goquery.Selection) bool {
		text := strings.TrimSpace(s.Text())
		if len(text) > 50 {
			snippet = text
			return false // 终止循环
		}
		return true
	})

	// 提取所有 <img src>
	images := []string{}
	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		if src, ok := s.Attr("src"); ok && strings.HasPrefix(src, "http") {
			images = append(images, src)
		}
	})

	return &OGDataa{
		Title:          getMeta("property", "og:title"),
		Description:    getMeta("property", "og:description"),
		Image:          getMeta("property", "og:image"),
		URL:            getMeta("property", "og:url"),
		CanonicalURL:   canonical,
		TwitterTitle:   getMeta("name", "twitter:title"),
		TwitterDesc:    getMeta("name", "twitter:description"),
		TwitterImage:   getMeta("name", "twitter:image"),
		PageTitle:      pageTitle,
		ContentSnippet: snippet,
		ImageList:      images,
	}, nil
}
