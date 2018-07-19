package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"./getdoc"
	"./parsedoc"
)

type mail struct {
	mid       int
	sender    string
	recievers []string
	isValid   bool
}

func getMailByMid(mid int) (mail, error) {
	var m mail
	for {
		// log.Printf("try loading %v...\n", mid)
		doc, err := getdoc.GetDocumentByMid(mid)
		if err != nil {
			log.Println("!getdoc.GetDocumentByMid ->")
			log.Println(err)
			// retry
			// return m, err
			time.Sleep(3 * time.Second)
			continue
		}
		if !parsedoc.IsNotDenialMessage(doc) || !parsedoc.IsValidMid(doc) || !parsedoc.IsTheMidExists(doc) {
			log.Println("!parsedoc, document not valid")
			m = mail{
				isValid: false,
			}
			return m, nil
		}
		sender, recievers, err := parsedoc.GetSenderAndRecievers(doc)
		if err != nil {
			log.Println("!parsedoc.GetSenderAndRecievers ->")
			log.Println(err)
			// retry
			// return m, err
			continue
		}
		m = mail{
			mid:       mid,
			sender:    sender,
			recievers: recievers,
			isValid:   true,
		}
		if len(m.sender) == 0 {
			// 送信者がいないのは流石におかしい
			// retry
			continue
		}
		break
	}
	return m, nil
}

func convergeDataByMidBetween(startMid int, endMid int) chan []mail {
	var wg sync.WaitGroup
	var err error
	remains := 0
	allcount := endMid - startMid + 1
	de := make([]mail, allcount, allcount)
	for i := startMid; i <= endMid; i++ {
		wg.Add(1)
		remains++
		go func(mid int) {
			defer func() {
				wg.Done()
				remains--
				log.Printf("mid %v ended. %v / %v remains\n", mid, remains, allcount)
			}()
			de[mid-startMid], err = getMailByMid(mid)
			if err != nil {
				log.Println("!getMailByMid ->")
				log.Println(err)
				de[mid-startMid] = mail{mid: mid}
			}
		}(i)
	}
	c := make(chan []mail)
	go func() {
		defer close(c)
		wg.Wait()
		c <- de
		return
	}()
	return c
}

func main() {
	c := convergeDataByMidBetween(603657, 622345)
	log.Println("procedure start")
	arr := <-c
	for m := range arr {
		fmt.Println(arr[m])
	}
}
