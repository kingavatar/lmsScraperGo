package scraper

import (
        "log"
        "os"
        "errors"
        "bytes"
        "net/http/cookiejar"
        "net/http"
)
var app App
//Start Scraper 
func Start() error {
  file, errlog := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
  if errlog != nil {
        log.Fatal(errlog)
    }

  log.SetOutput(file)
  jar, _ := cookiejar.New(nil)
	app = App{
		Client: &http.Client{Jar: jar},
	}
  loginstatus,_:=app.login()
  if(len(loginstatus)>0){
    return errors.New(loginstatus)
  }
  return nil 
}

//Stop Scraper
func Stop() error{
  app.logout()
  return nil
}

//GetEvents to get Events byte array
func GetEvents() []byte{
  events:=app.getEvents()
  var buffer bytes.Buffer
  for _,event := range events{
    buffer.WriteString(event.Title)
    buffer.WriteString(" ")
    buffer.WriteString(event.Deadline)
    buffer.WriteString("\n")
  }
  return buffer.Bytes()
}

//GetUnknownResponse for Unkown Args
func GetUnknownResponse() []byte{
    var buffer bytes.Buffer
    buffer.WriteString("Unknown Command\n")
    return buffer.Bytes()
}

//SetUserPass sets username and password
func SetUserPass(user string,pass string){
  username = user
  password = pass
}

//GetUserName returns username
func GetUserName() string{
  return username
}