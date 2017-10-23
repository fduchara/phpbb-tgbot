package parseforum

import (
	"golang.org/x/net/html"
	"log"
	"net/http"
	"net/url"
)

func GetNew(urlLogin string, urlFindNew string, urlMarkRead string, username string, password string) string {

	form := url.Values{
		"username": {username},
		"password": {password},
		"login":    {"Login"},
	}

	log.Print("Login forum")
	resp, err := http.PostForm(urlLogin, form)
	if err != nil {
		log.Panic(err)
	}

	defer resp.Body.Close()

	cookis := resp.Cookies()[len(resp.Cookies())-1]

	log.Print("Check new messages")
	req, err := http.NewRequest("GET", urlFindNew, nil)
	if err != nil {
		log.Panic(err)
	}

	req.AddCookie(cookis)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Panic(err)
	}

	reply := ""
	z := html.NewTokenizer(resp.Body)
	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			log.Print("Mark read all messages")
			req, err = http.NewRequest("GET", urlMarkRead, nil)
			if err != nil {
				log.Panic(err)
			}

			req.AddCookie(cookis)
			resp, err = http.DefaultClient.Do(req)
			if err != nil {
				log.Panic(err)
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
						log.Println(t.Data)
					}
				}
			}
		}
	}
}
