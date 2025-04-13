package preview

import (
	"errors"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

// OGData 保存提取的 Open Graph 数据
type OGData struct {
	Title       string
	Description string
	Image       string
	URL         string
}

// FetchOGTags 从指定 URL 获取网页并提取 OG 标签信息
func FetchOGTags(targetURL string) (*OGData, error) {
	resp, err := http.Get(targetURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch page: " + resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	getMeta := func(property string) string {
		content := ""
		doc.Find("meta").Each(func(i int, s *goquery.Selection) {
			if val, exists := s.Attr("property"); exists && val == property {
				if c, ok := s.Attr("content"); ok {
					content = c
				}
			}
		})
		return content
	}

	return &OGData{
		Title:       getMeta("og:title"),
		Description: getMeta("og:description"),
		Image:       getMeta("og:image"),
		URL:         getMeta("og:url"),
	}, nil
}
