package scraper

import (
	// "fmt"
	"net/http"
	// "net/http/cookiejar"
	"log"
	"net/url"
	// "io/ioutil"
	"github.com/PuerkitoBio/goquery"
)

const (
	baseURL = "https://lms.iiitb.ac.in/moodle/"
)
//App is the base http Client
type App struct {
	Client *http.Client
	isLoggedIn bool
}
type loginToken struct {
	Token string
}
//SessKey Type
type SessKey struct {
	SessKey string
}
var (
	username = "XXXXXX"
	password = "XXXXXX"
)

func (app *App) getToken() loginToken {
	loginURL := baseURL + "login/index.php"
	client := app.Client

	response, err := client.Get(loginURL)

	if err != nil {
		log.Fatalln("Error fetching response. ", err)
	}

	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal("Error loading HTTP response body. ", err)
	}

	token, _ := document.Find("input[name='logintoken']").Attr("value")

	loginToken := loginToken{
		Token: token,
	}
	
	return loginToken
}
func (app *App) getSessKey() SessKey {
	loginURL := baseURL + "my/"
	client := app.Client

	response, err := client.Get(loginURL)

	if err != nil {
		log.Fatalln("Error fetching response. ", err)
	}

	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal("Error loading HTTP response body. ", err)
	}

	sessKey,_:=document.Find("input[type=hidden][name=sesskey]").Attr("value")

	SessKey := SessKey{
		SessKey:sessKey,
	}

	return SessKey
}
func  (app *App) login() (string,error){
	client := app.Client

	loginToken := app.getToken()

	loginURL := baseURL + "login/index.php"

	data := url.Values{
		"logintoken": {loginToken.Token},
		"username":        {username},
		"password":     {password},
	}

	response, err := client.PostForm(loginURL, data)

	if err != nil {
		log.Fatalln(err)
		return "Error Login",err
	}

	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)
	loginerror:=document.Find("#loginerrormessage").Text()
	if err != nil {
		log.Fatal("Error loading HTTP response body. ", err)
		return "Error Login",err
	}
	if *Debug{
		log.Println("Response Logging in: ",response.Status)
	}
	if len(loginerror)==0 {
		app.isLoggedIn = true
	}
	return loginerror,nil
}

func (app *App) logout(){
	client := app.Client

	sessKey := app.getSessKey()

	logoutURL := baseURL + "login/logout.php?sesskey="+sessKey.SessKey

	response, err := client.Get(logoutURL)

	if err != nil {
		log.Fatalln("Error fetching response. ", err)
	}

	defer response.Body.Close()
	if *Debug{
		log.Println("Response Status Logging Out: ",response.Status)
	}
	if app.isLoggedIn {
		app.isLoggedIn = false
	}

}