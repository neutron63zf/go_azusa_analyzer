package parsedoc

import (
	"log"
	"regexp"
	"strings"

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

// GetSenderAndReciever returns sender recievers slice
func GetSenderAndRecievers(doc *goquery.Document) (string, []string, error) {
	// azusaのメッセージの部分のマークアップは存在しないため、正規表現しかない
	html, err := doc.Find(".azusa-log").Html()
	if err != nil {
		log.Println("!doc.Find->")
		return "", nil, err
	}
	// ここのスペースはどうやらただの半角スペースでは 無いらしい
	// おそらく&nbsp;の変換先なのだが...
	r := regexp.MustCompile(`\[#\d+\]  (.*?) → (.*?)<br/>`)
	if !r.MatchString(html) {
		return "", nil, nil
	}
	result := r.FindAllStringSubmatch(html, -1)
	senderAndReciever := result[0][1:]
	sender := senderAndReciever[0]
	// "、"で区切られているので...
	recievers := strings.Split(senderAndReciever[1], "、")
	return sender, recievers, nil
}
