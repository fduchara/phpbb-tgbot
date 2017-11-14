package parseforum

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

func GetNew(forumurl string, username string, password string, debug bool) string {

	cookies := login(forumurl, username, password, debug)
	//	urlNewMessage := findNew(forumurl, cookies, debug)
	//	if urlNewMessage != "" {
	//		urlNewMessage := forumurl+"/viewtopic.php?"+urlNewMessage
	urlNewMessage := "http://partizanen.su/forum/viewtopic.php?f=29&t=2679"

	bodyStringsHTML := getBodyHTML(urlNewMessage, cookies, debug)

	//	postMessage := parseMessage(bodyStringsHTML)
	//	log.Println(postMessage)

	_ = parseMessage(bodyStringsHTML)

	//	message := cleanMessageQuote(postMessage)
	//	log.Println(message)

	//	}

	return ""
}

func login(forumurl string, username string, password string, debug bool) *http.Cookie {
	urllogin := forumurl + "/ucp.php?mode=login"
	form := url.Values{
		"username": {username},
		"password": {password},
		"login":    {"Login"},
	}
	if debug {
		log.Print("Login forum: " + urllogin + " username: " + username + " password: " + password)
	}
	resp, err := http.PostForm(urllogin, form)
	defer resp.Body.Close()
	if err != nil {
		log.Panic(err)
	}
	if debug {
		log.Println(resp)
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
	urlnew := forumurl + "/search.php?search_id=unreadposts"
	if debug {
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
	if debug {
		log.Println("response:\n", string(body))
	}

	re := regexp.MustCompile("f=([0-9]+)&amp;t=([0-9]+)&amp;view=unread#unread")
	postMessage := re.FindString(string(body))
	postMessage = strings.Replace(postMessage, "amp;", "", -1)
	if postMessage != "" {
		return postMessage
	}
	return ""
}

func getBodyHTML(forumurl string, cookies *http.Cookie, debug bool) []string {
	urlnew := forumurl
	if debug {
		log.Print("get and read new messages: " + urlnew)
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
	if debug {
		log.Println("response:\n", string(body))
	}

	bodyStringsHTML := strings.Split(string(body), "\n")

	return bodyStringsHTML
}

func parseMessage(bodyStringsHTML []string) string {

	reThem := regexp.MustCompile("Re: (.*)</a></div><div style=")
	reAuthor := regexp.MustCompile("<b class=\"postauthor\">(.)+</b>")
	reBody := regexp.MustCompile("<div class=\"postbody\">(.)+</div>")

	postThem := ""
	postAuthor := ""
	postBody := ""

	for _, s := range bodyStringsHTML {

		them := reThem.FindString(s)
		if them != "" {
			them = strings.Replace(them, "</a></div><div style=", "", -1)
//			log.Println(them)
			postThem = them
		}

		author := reAuthor.FindString(s)
		if author != "" {
			author = strings.Replace(author, "<b class=\"postauthor\">", "", -1)
			author = strings.Replace(author, "</b>", "", -1)
//			log.Println(author)
			postAuthor = author
		}

		body := reBody.FindString(s)
		if body != "" {
			if !strings.Contains(body, "<br />_________________<br />") {

				body = strings.Replace(body, "<div class=\"postbody\">", "", -1)
				body = strings.Replace(body, "</div>", "", -1)

				body = cleanMessageQuote(body)
				body = cleanMessageStyle(body)
//				log.Println(body)
				postBody = body
			}
		}
	}

	log.Println(postThem)
	log.Println(postAuthor)
	log.Println(postBody)

	return postThem+postAuthor+postBody
}

func cleanMessageQuote(postMessage string) string {
	if strings.Contains(postMessage, "<div class=\"quotetitle\">") {
		postMessage = strings.Replace(postMessage, "<div class=\"quotetitle\">", "", 1)
		postMessage = strings.Replace(postMessage, "<div class=\"quotecontent\">", "\n", 1)
		postMessage = strings.Replace(postMessage, "```", "", 2)
		postMessage = "```" + postMessage
		postMessage = strings.Replace(postMessage, "<br />", "```\n", 1)
		return cleanMessageQuote(postMessage)
	}
	return postMessage
}

func cleanMessageStyle(postMessage string) string {
	if strings.Contains(postMessage, "<span style=") {
		re := regexp.MustCompile("<span style=\"(.)+\">")
		postMessage = re.ReplaceAllLiteralString(postMessage, "")
		postMessage = strings.Replace(postMessage, "</span>", "", -1)
	}
	if strings.Contains(postMessage, "<br />") {
		postMessage = strings.Replace(postMessage, "<br />", "", -1)
	}
	if strings.Contains(postMessage, "href=\"") {
		postMessage = strings.Replace(postMessage, "<!-- m -->", "", 2)
		re := regexp.MustCompile("<a class=\"postlink\" href=\"(.)+\">")
		postMessage = re.ReplaceAllLiteralString(postMessage, "")
		postMessage = strings.Replace(postMessage, "</a>", " ", 1)
	}

	if strings.Contains(postMessage, "<") {
		if strings.Contains(postMessage, ">") {
			re := regexp.MustCompile("<(.)+>")
			postMessage = re.ReplaceAllLiteralString(postMessage, " ")
		}
	}
	return postMessage
}
