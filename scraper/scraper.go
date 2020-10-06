package scraper

import (
	      // "fmt"
        // "log"
        // "os"
        "errors"
        "bytes"
        "strings"
        "net/http/cookiejar"
        "net/http"
)
var app App

//Debug to set debugging mode 
var Debug *bool
//Start Scraper 
func Start() error {
  // file, errlog := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
  // if errlog != nil {
        // log.Fatal(errlog)
    // }

  // log.SetOutput(file)
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
  if app.isLoggedIn{
    app.logout()
  } else{
    return errors.New("Service was not Started to Stop so ... did not stop the service")
  }
  return nil
}

//GetEventsTerm to get Events byte array for the Terminal
func GetEventsTerm() []byte{
  events:=app.getEvents()
  var buffer bytes.Buffer
  courseID:=0
  var setColor bool
  for _,event := range events{
    if courseID!=event.CourseID{
      courseID = event.CourseID
      buffer.WriteString("\\e[38;2;250;169;22m")
      setColor = true
    }
    buffer.WriteString(event.Name)
    buffer.WriteString(" ")
    buffer.WriteString(event.Formattedtime)
    buffer.WriteString("\\n")
    if setColor {
      setColor = false
      buffer.WriteString("\\e[38;2;251;255;254m")
    }
  }
  buffer.WriteString("\\e[0m\n")
  return buffer.Bytes()
}


//GetEvents to get Events byte array
func GetEvents() []byte{
  events:=app.getEvents()
  var buffer bytes.Buffer
  for _,event := range events{
    buffer.WriteString(event.Name)
    buffer.WriteString(" ")
    buffer.WriteString(event.Formattedtime)
    buffer.WriteString("\n")
  }
  return buffer.Bytes()
}

//GetEventsConky to get Events byte array
func GetEventsConky() []byte{
  events:=app.getEvents()
  var buffer bytes.Buffer
  courseID:=0
  var setColor bool
  for _,event := range events{
    if courseID!=event.CourseID{
      courseID = event.CourseID
      buffer.WriteString("${#FAA916}")
      setColor = true
    }
    buffer.WriteString(event.Name)
    buffer.WriteString(" ")
    buffer.WriteString(event.Formattedtime)
    buffer.WriteString("\n")
    if setColor {
      setColor = false
      buffer.WriteString("${#FBFFFE}")
    }
  }
  return buffer.Bytes()
}

//GetAnnouncementsTerm to get Announcements byte array
func GetAnnouncementsTerm() []byte{
  announcements:=app.getAnnouncements()
  var buffer bytes.Buffer
  courseID:=0
  var setColor bool
  for _,announcement := range announcements{
    if courseID!=announcement.CourseID{
      courseID = announcement.CourseID
      buffer.WriteString("\\e[38;2;250;169;22m")
      setColor = true
    }
    buffer.WriteString(announcement.Date)
    buffer.WriteString("   ")
    buffer.WriteString(announcement.Name)
    buffer.WriteString("   ")
    info:=strings.ReplaceAll(announcement.Info,"#","\\#")
    buffer.WriteString(info)
    buffer.WriteString("\n")
    if setColor {
      setColor = false
      buffer.WriteString("\\e[38;2;251;255;254m")
    }
  }
  buffer.WriteString("\\e[0m\n")
  return buffer.Bytes()
}

//GetAnnouncements to get Announcements byte array
func GetAnnouncements() []byte{
  announcements:=app.getAnnouncements()
  var buffer bytes.Buffer
  for _,announcement := range announcements{
    buffer.WriteString(announcement.Date)
    buffer.WriteString("   ")
    buffer.WriteString(announcement.Name)
    buffer.WriteString("   ")
    info:=strings.ReplaceAll(announcement.Info,"#","\\#")
    buffer.WriteString(info)
    buffer.WriteString("\n")
  }
  return buffer.Bytes()
}
//GetAnnouncementsConky to get Announcements byte array
func GetAnnouncementsConky() []byte{
  announcements:=app.getAnnouncements()
  var buffer bytes.Buffer
  courseID:=0
  var setColor bool
  for _,announcement := range announcements{
    if courseID!=announcement.CourseID{
      courseID = announcement.CourseID
      buffer.WriteString("${#FAA916}")
      setColor = true
    }
    buffer.WriteString(announcement.Date)
    buffer.WriteString("   ")
    buffer.WriteString(announcement.Name)
    buffer.WriteString("   ")
    info:=strings.ReplaceAll(announcement.Info,"#","\\#")
    buffer.WriteString(info)
    buffer.WriteString("\n")
    if setColor {
      setColor = false
      buffer.WriteString("${#FBFFFE}")
    }
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

//SetDebugMode to set debug Mode
func SetDebugMode(debug *bool){
    Debug=debug
}
