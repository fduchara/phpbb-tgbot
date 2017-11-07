package parseforum

import (
	"golang.org/x/net/html"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func GetNew(forumurl string, username string, password string, debug bool) string {

	cookies := login(forumurl, username, password, debug)
	return findNew(forumurl, cookies, debug)
}

func login(forumurl string, username string, password string, debug bool) *http.Cookie {
	urllogin := forumurl+"/ucp.php?mode=login"
	form := url.Values{
		"username": {username},
		"password": {password},
		"login":    {"Login"},
	}
	if debug {
		log.Print("Login forum: "+urllogin+" username: "+username+" password: "+password)
	}
	resp, err := http.PostForm(urllogin, form)
	defer resp.Body.Close()
	if err != nil {
		log.Panic(err)
	}
	if debug {
		log.Print(resp)
		body, _ := ioutil.ReadAll(resp.Body)
		log.Println("response:\n", string(body))
	}
	cookies := resp.Cookies()[len(resp.Cookies())-1]
	if debug {
		log.Print("Cookies")
		log.Print(cookies)
	}
	return cookies
}

func findNew(forumurl string, cookies *http.Cookie, debug bool) string {
	urlnew := forumurl+"/search.php?search_id=unreadposts"
	if debug {
		log.Print("Check new messages: "+urlnew)
	}
	req, err := http.NewRequest("GET", urlnew, nil)
	if err != nil {
		log.Panic(err)
	}
	req.AddCookie(cookies)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Panic(err)
	}
	if debug {
		log.Print(resp)
		body, _ := ioutil.ReadAll(resp.Body)
		log.Println("response:\n", string(body))
	}

	reply := ""
	z := html.NewTokenizer(resp.Body)
	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:

			markRead(forumurl, cookies, debug)
			return reply

		case tt == html.StartTagToken:
			t := z.Token()
			for _, a := range t.Attr {
				if a.Key == "class" && a.Val == "topictitle" {
					tt := z.Next()
					if tt == html.TextToken {
						t := z.Token()
						reply += "\n " + t.Data
						if debug {
							log.Println(t.Data)
						}
					}
				}
			}
		}
	}
}

func markRead(forumurl string, cookies *http.Cookie, debug bool) {
	urlmarkread := forumurl+"/index.php?hash=376996bd&mark=forums"
	if debug {
		log.Print("Mark read all messages: "+urlmarkread)
	}
	req, err := http.NewRequest("GET", urlmarkread, nil)
	if err != nil {
		log.Panic(err)
	}
	req.AddCookie(cookies)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Panic(err)
	}
	if debug {
		log.Print(resp)
		body, _ := ioutil.ReadAll(resp.Body)
		log.Println("response:\n", string(body))
	}
}
