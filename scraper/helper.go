package scraper

import(
	"log"
	"github.com/PuerkitoBio/goquery"
	// "fmt"
	"strconv"
	"encoding/json"
	 "io/ioutil"
	"net/http"
	"bytes"
)

//Event Type
type Event struct {
	Title string
	CourseID string
	EventID string
	Deadline string
	URL string
	Name string
}
//Course Type
type Course struct{
	Name string `json:"fullname"`
	Code string `json:"shortname"`
	CourseID int `json:"id"`
	URL string `json:"viewurl"`
	Isfavorite bool `json:"isfavourite"`
}

//Announcement Type
type Announcement struct{
	Date string
	Name string
	Info string
	CourseID int
}
//getEvent() function
func (app *App) getEvents() []Event{
	var Events []Event
	calenderEventsURL:=baseURL+"calendar/view.php"
	client := app.Client

	response, err := client.Get(calenderEventsURL)

	if err != nil {
		log.Fatalln("Error fetching response. ", err)
	}

	defer response.Body.Close()
	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal("Error loading HTTP response body. ", err)
	}
	document.Find(".event").Each(func(i int, s *goquery.Selection) {
    courseid, _ := s.Attr("data-course-id")
	eventid,_ := s.Attr("data-event-id")
	title,_ :=s.Attr("data-event-title")
	url,_:= s.Find(".description a").Attr("href")
	name:=s.Find("name").Text()
	deadline:=s.Find(".date a").Text()
	event:= Event{
		Name:name,
		Title:title,
		URL:url,
		EventID:eventid,
		CourseID:courseid,
		Deadline:deadline,
	}
	Events = append(Events,event)
})
	return Events
}

func (app *App) getCourses() []Course{
	// var Courses []Course
	client:=app.Client
	sessKey:=app.getSessKey()
	postURL:="https://lms.iiitb.ac.in/moodle/lib/ajax/service.php?sesskey="+sessKey.SessKey+"&info=core_course_get_enrolled_courses_by_timeline_classification"
	formData:=[]byte(`[{"index":0,"methodname":"core_course_get_enrolled_courses_by_timeline_classification","args":{"offset":0,"limit":96,"classification":"all","sort":"ul.timeaccess desc"}}]`)
	req, _:= http.NewRequest("POST", postURL, bytes.NewBuffer(formData))
	req.Header.Set("User-Agent","Mozilla/5.0 (X11; Linux x86_64; rv:80.0) Gecko/20100101 Firefox/80.0")
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	postResponse, postErr := client.Do(req)
	if postErr != nil {
		log.Fatalln("Error fetching Post response. ", postErr)
	}
	defer postResponse.Body.Close()
	postdocument, _ := ioutil.ReadAll(postResponse.Body)
	// fmt.Println(string(postdocument))
	type CourseResult []struct{
		Data struct{
			Courses []Course `json:"courses"`
		} `json:"data"`
	}
	
	var result CourseResult

	err := json.Unmarshal(postdocument,&result)
	 if err != nil {
        log.Fatalln("Error Reading json Data: ",err)
	}
	return result[0].Data.Courses
}

func (app *App) getAnnouncements() []Announcement{
	var Announcements []Announcement
	for _,course :=range app.getCourses(){
			if !course.Isfavorite {continue}
			coursePageURL:=baseURL+"course/view.php?id="+strconv.FormatInt(int64(course.CourseID), 10)
			client := app.Client

			response, err := client.Get(coursePageURL)

			if err != nil {
				log.Fatalln("Error fetching response. ", err)
			}

			defer response.Body.Close()
			document, err := goquery.NewDocumentFromReader(response.Body)
			if err != nil {
				log.Fatal("Error loading HTTP response body. ", err)
			}
			document.Find(".post").Each(func(i int, s *goquery.Selection) {
			name:=s.Find(".name").Text()
			date:=s.Find(".date").Text()
			info:=s.Find(".info").Text()
			announcement:= Announcement{
				Name:name,
				Date:date,
				Info:info,
				CourseID:course.CourseID,
			}
			Announcements = append(Announcements,announcement)
		})
	}
	return Announcements
}