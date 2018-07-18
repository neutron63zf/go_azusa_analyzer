package getdoc

import (
	"log"
	"net/http"
	"net/url"
	"strconv"

	"../env"
	"github.com/PuerkitoBio/goquery"
)

func appendSessionCookie(req *http.Request) *http.Request {
	session := &http.Cookie{
		Name:  "session",
		Value: env.Session,
	}
	a103session := &http.Cookie{
		Name:  "a103session",
		Value: env.A103Session,
	}
	req.AddCookie(session)
	req.AddCookie(a103session)
	return req
}

func createHTTPRequestWithMid(mid int) (*http.Request, error) {
	values := url.Values{}
	values.Add("mid", strconv.Itoa(mid))
	base := "https://www.a103.net/azusa/form_oth/view_favs.cgi"
	url := base + "?" + values.Encode()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("!http.NewRequest ->")
		return nil, err
	}
	return req, nil
}

// GetDocumentByMid get azusa message goquery document by mid
func GetDocumentByMid(mid int) (*goquery.Document, error) {
	req, err := createHTTPRequestWithMid(mid)
	if err != nil {
		log.Println("!createHTTPRequestWithMid ->")
		return nil, err
	}
	req = appendSessionCookie(req)

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		log.Println("!client.Do ->")
		return nil, err
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Println("!goquery.NewDocumentFromReader ->")
		return nil, err
	}
	return doc, nil
}
