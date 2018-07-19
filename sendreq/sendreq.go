package sendreq

import (
	"net/http"
)

const (
	limit = 1
)

var (
	remains = 0
)

// Resp is custom response
type Resp struct {
	Resp *http.Response
	Err  error
}

// DoRequest limits request count
func DoRequest(client *http.Client, req *http.Request) chan Resp {
	// remains++
	// log.Printf("%v request remains", remains)
	for {
		if remains < limit {
			// log.Printf("%v request remains(waiting)", remains)
			break
		}
	}
	// log.Printf("wait ended")
	c := make(chan Resp)
	go func() {
		defer close(c)
		resp, err := client.Do(req)
		c <- Resp{
			Resp: resp,
			Err:  err,
		}
		// remains--
		// log.Printf("%v request remains(ended)", remains)
	}()
	return c
}
