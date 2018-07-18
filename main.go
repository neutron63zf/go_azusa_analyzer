package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"

	"./env"
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

func createHTTPRequestWithMid(mid int) *http.Request {
	values := url.Values{}
	values.Add("mid", strconv.Itoa(mid))
	base := "https://www.a103.net/azusa/form_oth/view_favs.cgi"
	url := base + "?" + values.Encode()
	req, _ := http.NewRequest("GET", url, nil)
	return req
}

func readContentString(resp *http.Response) string {
	buffer, _ := ioutil.ReadAll(resp.Body)
	return string(buffer)
}

func isDenialMessage(content string) bool {
	r := regexp.MustCompile(`不正な操作が行われたか`)
	return r.MatchString(content)
}

func main() {
	req := createHTTPRequestWithMid(622592)
	req = appendSessionCookie(req)

	client := new(http.Client)
	resp, _ := client.Do(req)

	content := readContentString(resp)

	fmt.Println(content)

	fmt.Println(isDenialMessage(content))

}
