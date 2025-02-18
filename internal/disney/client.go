package disney

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"github.com/PuerkitoBio/goquery"
)

type AttractionInfo struct {
	URL string
	Name string
	WaitTime string
}

func GetAttractionInfoList(park string) (attractionInfoList []AttractionInfo) {
	base, _ := url.Parse("https://tokyodisneyresort.info/realtime.php")
	reference, _ := url.Parse("?park=" + park)
	endpoint := base.ResolveReference(reference).String()

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	var client *http.Client = &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
	}

	doc.Find(".realtime-attr a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists {
			href = "N/A" // `href` がない場合のフォールバック
		}

		// `<div class="realtime-attr-name">` のテキストを取得
		name := cleanText(s.Find(".realtime-attr-name").Text())

		// `<div class="realtime-attr-condition">` のテキストを取得
		waitTime := cleanText(s.Find(".realtime-attr-condition").Text())

		attractionInfo := AttractionInfo{
			URL: href,
			Name: name,
			WaitTime: waitTime,
		}

		attractionInfoList = append(attractionInfoList, attractionInfo)
	})

	return attractionInfoList
}

func cleanText(text string) string {
	trimmed := strings.TrimSpace(text)     // 前後のスペース・改行を削除
	words := strings.Fields(trimmed)       // 空白で区切ってスライスに変換（余分なスペースを除去）
	return strings.Join(words, " ")        // スペース1つで結合
}