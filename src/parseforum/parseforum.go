package parseforum

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

func GetNew(forumurl string, username string, password string, debug int) string {

	cookies := login(forumurl, username, password, debug)
	urlNewMessage := findNew(forumurl, cookies, debug)

	var reply, replyThem, replyAuthor, replyMessage = "", "", "", ""

	if urlNewMessage != "" {
		urlNewMessage := forumurl + "/viewtopic.php?" + urlNewMessage
		bodyStringsHTML := getBodyHTML(urlNewMessage, cookies, debug)

		replyThem, replyAuthor, replyMessage = parseMessage(bodyStringsHTML, debug)
		replyMessage = cleanMessageQuote(replyMessage, debug)
		replyMessage = cleanMessageStyle(replyMessage, debug)
		reply = "Тема: " + replyThem + "\nОт: " + replyAuthor + "\n" + replyMessage + "\n\nURL: " + urlNewMessage
	}
	return reply
}

func login(forumurl string, username string, password string, debug int) *http.Cookie {
	urllogin := forumurl + "/ucp.php?mode=login"
	form := url.Values{
		"username": {username},
		"password": {password},
		"login":    {"Login"},
	}
	if debug > 0 {
		log.Print("Login forum: " + urllogin)
	}
	if debug > 1 {
		log.Print("Login username: " + username + " password: " + password)
	}
	resp, err := http.PostForm(urllogin, form)
	defer resp.Body.Close()
	if err != nil {
		log.Panic(err)
	}
	if debug > 2 {
		log.Println(resp)
		body, _ := ioutil.ReadAll(resp.Body)
		log.Println("response:\n", string(body))
	}
	cookies := resp.Cookies()[len(resp.Cookies())-1]
	if debug > 0 {
		log.Print("Cookies")
		log.Print(cookies)
	}
	return cookies
}

func findNew(forumurl string, cookies *http.Cookie, debug int) string {
	urlnew := forumurl + "/search.php?search_id=unreadposts"
	if debug > 0 {
		log.Print("Check new messages: " + urlnew)
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
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panic(err)
	}
	if debug > 2 {
		log.Println("Response:\n", string(body))
	}

	re := regexp.MustCompile("f=([0-9]+)&amp;t=([0-9]+)&amp;view=unread#unread")
	postMessage := re.FindString(string(body))
	if postMessage != "" {
		postMessage = strings.Replace(postMessage, "amp;", "", -1)
		return postMessage
	}
	return ""
}

func getBodyHTML(forumurl string, cookies *http.Cookie, debug int) []string {
	if debug > 0 {
		log.Print("Get and read new messages: " + forumurl)
	}
	req, err := http.NewRequest("GET", forumurl, nil)
	if err != nil {
		log.Panic(err)
	}
	req.AddCookie(cookies)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Panic(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panic(err)
	}
	if debug > 2 {
		log.Println("Response:\n", string(body))
	}

	bodyStringsHTML := strings.Split(string(body), "\n")

	return bodyStringsHTML
}

func parseMessage(bodyStringsHTML []string, debug int) (postThem string, postAuthor string, postBody string) {

	reThem := regexp.MustCompile("Re: (.*)</a></div><div style=")
	reAuthor := regexp.MustCompile("<b class=\"postauthor\">(.)+</b>")
	reBody := regexp.MustCompile("<div class=\"postbody\">(.)+</div>")

	postThem = ""
	postAuthor = ""
	postBody = ""

	for _, s := range bodyStringsHTML {

		if postThem == "" {
			them := reThem.FindString(s)
			if them != "" {
				them = strings.Replace(them, "</a></div><div style=", "", -1)
				them = strings.Replace(them, "Re: ", "", -1)
				postThem = them
			}
		}

		author := reAuthor.FindString(s)
		if author != "" {
			author = strings.Replace(author, "<b class=\"postauthor\">", "", -1)
			author = strings.Replace(author, "</b>", "", -1)
			postAuthor = author
		}

		body := reBody.FindString(s)
		if body != "" {
			if !strings.Contains(body, "<br />_________________<br />") {

				body = strings.Replace(body, "<div class=\"postbody\">", "", -1)
				postBody = body
			}
		}
	}

	return postThem, postAuthor, postBody
}

func cleanMessageQuote(postMessage string, debug int) string {
	if strings.Contains(postMessage, "<div class=\"quotetitle\">") {
		if debug > 1 {
			log.Print("cleanMessageQuote IN: " + postMessage)
		}
		postMessage = strings.Replace(postMessage, "<div class=\"quotetitle\">", " ", 1)
		postMessage = strings.Replace(postMessage, "</div><div class=\"quotecontent\">", "\n", 1)
		postMessage = strings.Replace(postMessage, "```", "", 2)
		postMessage = "```" + postMessage
		postMessage = strings.Replace(postMessage, "</div><br />", "```\n", 1)
		if debug > 1 {
			log.Print("cleanMessageQuote OUT: " + postMessage)
		}
		return cleanMessageQuote(postMessage, debug)
	}
	return postMessage
}

func cleanMessageStyle(postMessage string, debug int) string {
	if strings.Contains(postMessage, "<span style=") {
		if debug > 1 {
			log.Print("cleanMessageStyle IN: " + postMessage)
		}
		re := regexp.MustCompile("<span style=\"(.)+\">")
		postMessage = re.ReplaceAllLiteralString(postMessage, "")
		postMessage = strings.Replace(postMessage, "</span>", "", -1)
		if debug > 1 {
			log.Print("cleanMessageStyle OUT: " + postMessage)
		}
	}

	if strings.Contains(postMessage, "href=\"") {
		if debug > 1 {
			log.Print("cleanMessageStyle href IN: " + postMessage)
		}
		postMessage = strings.Replace(postMessage, "<!-- m -->", "", 2)
		re := regexp.MustCompile("<a class=\"postlink\" href=\"(.)+\">")
		postMessage = re.ReplaceAllLiteralString(postMessage, "")
		postMessage = strings.Replace(postMessage, "</a>", " ", 1)
		if debug > 1 {
			log.Print("cleanMessageStyle href OUT: " + postMessage)
		}
	}

	if strings.Contains(postMessage, "<") {
		if strings.Contains(postMessage, ">") {
			if debug > 1 {
				log.Print("cleanMessageStyle tag IN: " + postMessage)
			}
			re := regexp.MustCompile("<(.)+>")
			postMessage = re.ReplaceAllLiteralString(postMessage, " ")
			if debug > 1 {
				log.Print("cleanMessageStyle tag OUT: " + postMessage)
			}
		}
	}
	return postMessage
}
