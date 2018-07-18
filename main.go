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

func createHTTPRequestWithMid(mid int) (*http.Request, error) {
	values := url.Values{}
	values.Add("mid", strconv.Itoa(mid))
	base := "https://www.a103.net/azusa/form_oth/view_favs.cgi"
	url := base + "?" + values.Encode()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func readContentString(resp *http.Response) (string, error) {
	buffer, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(buffer), err
}

func isDenialMessage(content string) bool {
	r := regexp.MustCompile(`不正な操作が行われたか`)
	return r.MatchString(content)
}

func main() {
	req, err := createHTTPRequestWithMid(622592)
	if err != nil {
		fmt.Println("!createHTTPRequestWithMid ->")
		fmt.Println(err)
	}
	req = appendSessionCookie(req)

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("!client.Do ->")
		fmt.Println(err)
	}

	content, err := readContentString(resp)
	if err != nil {
		fmt.Println("!readContentString ->")
		fmt.Println(err)
	}

	fmt.Println(content)

	fmt.Println(isDenialMessage(content))

}
