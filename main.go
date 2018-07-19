package main

import (
	"fmt"

	"./getdoc"
	"./parsedoc"
)

func printSR(mid int) {
	doc, err := getdoc.GetDocumentByMid(mid)
	if err != nil {
		fmt.Println("!getdoc.GetDocumentByMid ->")
		fmt.Println(err)
		return
	}
	sender, recievers, err := parsedoc.GetSenderAndRecievers(doc)
	fmt.Println(mid)
	fmt.Println(sender)
	fmt.Println(recievers)
	fmt.Println()
}

func main() {
	for i := 622000; i < 622600; i++ {}
		// go をつけると、並列先のが終わる前にメインのスレッドが終わってしまい、何も出ない
		printSR(i)
	}
}
