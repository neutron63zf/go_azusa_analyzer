package main

import (
	"fmt"

	"./getdoc"
	"./parsedoc"
)

func main() {

	doc, err := getdoc.GetDocumentByMid(622592)
	if err != nil {
		fmt.Println("!getdoc.GetDocumentByMid ->")
		fmt.Println(err)
		return
	}
	html, _ := doc.Find(".azusa-log").Html()
	fmt.Println(html)
	// isDenialMessage := parsedoc.IsDenialMessage(doc)
	// fmt.Println(isDenialMessage)
	// timeStr := parsedoc.GetTimeStr(doc)
	// fmt.Println(timeStr)
	sender, _ := parsedoc.GetSender(doc)
	fmt.Println(sender)
}
