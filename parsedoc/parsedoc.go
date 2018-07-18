package parsedoc

import (
	"log"
	"regexp"

	"github.com/PuerkitoBio/goquery"
)

// IsDenialMessage returns whether doc is reject message
func IsDenialMessage(doc *goquery.Document) bool {
	text := doc.Find("body").Text()
	r1 := regexp.MustCompile(`不正な操作が行われたか、cookieが有効になっていません。`)
	r2 := regexp.MustCompile(`ログイン操作をやり直してください。`)
	return r1.MatchString(text) && r2.MatchString(text)
}

// GetTimeStr returns time string of message
func GetTimeStr(doc *goquery.Document) string {
	return doc.Find(".azusa-log i").Text()
}

// GetSender return sender name by string
func GetSender(doc *goquery.Document) (bool, error) {
	// 正規表現で組むしか無い
	html, err := doc.Find(".azusa-log").Html()
	if err != nil {
		log.Println("!doc.Find->")
		return false, err
	}
	// ここのスペースはどうやらただの半角スペースでは 無いらしい
	r1 := regexp.MustCompile(`\[#\d+\]  (.*?) → (.*?)<br/>`)
	result := r1.FindAllStringSubmatch(html, -1)
	if !r1.MatchString(html) {
		return false, nil
	}
	log.Println(result[0][2])
	log.Println([]int{1, 2, 3})
	return r1.MatchString(html), nil
}
