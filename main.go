package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"

	"./getdoc"
	"./parsedoc"
)

const (
	// MaxLimit 最大同時接続数
	MaxLimit = 150
)

type mail struct {
	mid       int
	sender    string
	recievers []string
	isValid   bool
}

func getMailByMid(mid int) (mail, error) {
	var m mail
	var interval int64
	interval = 2
	for {
		// log.Printf("try loading %v...\n", mid)
		doc, err := getdoc.GetDocumentByMid(mid)
		if err != nil {
			log.Println("!getdoc.GetDocumentByMid ->")
			log.Println(err)
			// retry
			// return m, err
			time.Sleep(time.Duration(rand.Int63n(interval)) * time.Second)
			interval = interval * 2
			continue
		}
		if !parsedoc.IsNotDenialMessage(doc) || !parsedoc.IsValidMid(doc) || !parsedoc.IsTheMidExists(doc) {
			log.Printf("!parsedoc, document not valid on mid %v", mid)
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
	allcount := endMid - startMid + 1
	remains := allcount
	concuurent := 0
	// 同時接続数を制限
	limitter := make(chan bool, MaxLimit)
	de := make([]mail, allcount, allcount)
	for i := startMid; i <= endMid; i++ {
		wg.Add(1)
		concuurent++
		limitter <- true
		go func(mid int) {
			// limitterの値を捨てる
			defer func() { <-limitter }()
			defer func() {
				wg.Done()
				remains--
				concuurent--
				log.Printf("mid %v ended. %v / %v remains. %v process working\n", mid, remains, allcount, concuurent)
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
	file, err := os.OpenFile("test.txt", os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Println("err opening file")
		log.Println(err)
		return
	}
	defer file.Close()
	log.Println("procedure start")
	start := 603657
	// end := 622345
	end := 603660
	c := convergeDataByMidBetween(start, end)
	arr := <-c
	var m mail
	for midx := range arr {
		m = arr[midx]
		// fmt.Println(m)
		fmt.Fprintf(file, "%v,%v,%v\n", m.mid, m.sender, strings.Join(m.recievers, ","))
	}
}
