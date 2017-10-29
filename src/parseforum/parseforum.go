package parseforum

import (
	"golang.org/x/net/html"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func GetNew(urlLogin string, urlFindNew string, urlMarkRead string, username string, password string, debug bool) string {

	form := url.Values{
		"username": {username},
		"password": {password},
		"login":    {"Login"},
	}

	if debug {
		log.Print("Login forum")
	}
	resp, err := http.PostForm(urlLogin, form)
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

	if debug {
		log.Print("Check new messages")
	}
	req, err := http.NewRequest("GET", urlFindNew, nil)
	if err != nil {
		log.Panic(err)
	}
	req.AddCookie(cookies)
	resp, err = http.DefaultClient.Do(req)
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

			if debug {
				log.Print("Mark read all messages")
			}
			req, err = http.NewRequest("GET", urlMarkRead, nil)
			if err != nil {
				log.Panic(err)
			}
			req.AddCookie(cookies)
			resp, err = http.DefaultClient.Do(req)
			if err != nil {
				log.Panic(err)
			}
			if debug {
				log.Print(resp)
				body, _ := ioutil.ReadAll(resp.Body)
				log.Println("response:\n", string(body))
			}

			return reply

		case tt == html.StartTagToken:
			t := z.Token()
			for _, a := range t.Attr {
				if a.Key == "class" && a.Val == "topictitle" {
					tt := z.Next()
					if tt == html.TextToken {
						t := z.Token()
						reply = reply + "\n " + t.Data
						if debug {
							log.Println(t.Data)
						}
					}
				}
			}
		}
	}
}
