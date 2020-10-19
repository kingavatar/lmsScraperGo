package scraper

import (
	"log"

	"github.com/PuerkitoBio/goquery"

	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

//NearEvent Type
type NearEvent struct {
	Title    string
	CourseID string
	EventID  string
	Deadline string
	URL      string
	Name     string
}

//Event Type
type Event struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	CourseID      int    `json:"course.id"`
	Instance      int    `json:"instance"`
	Eventtype     string `json:"eventtype"`
	Timestart     int    `json:"timestart"`
	Timeduration  int    `json:"timeduration"`
	Timesort      int    `json:"timesort"`
	Visible       int    `json:"visible"`
	Timemodified  int    `json:"timemodified"`
	Formattedtime string `json:"formattedtime"`
	URL           string `json:"url"`
}

//Course Type
type Course struct {
	Name       string `json:"fullname"`
	Code       string `json:"shortname"`
	CourseID   int    `json:"id"`
	URL        string `json:"viewurl"`
	Isfavorite bool   `json:"isfavourite"`
}

//Announcement Type
type Announcement struct {
	Date     string
	Name     string
	Info     string
	CourseID int
}

//getEvent() function
func (app *App) getNearEvents() []NearEvent {
	var Events []NearEvent
	calenderEventsURL := baseURL + "calendar/view.php"
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
		eventid, _ := s.Attr("data-event-id")
		title, _ := s.Attr("data-event-title")
		url, _ := s.Find(".description a").Attr("href")
		name := s.Find("name").Text()
		deadline := s.Find(".date a").Text()
		event := NearEvent{
			Name:     name,
			Title:    title,
			URL:      url,
			EventID:  eventid,
			CourseID: courseid,
			Deadline: deadline,
		}
		Events = append(Events, event)
	})
	return Events
}

func (app *App) getCourses() []Course {
	// var Courses []Course
	client := app.Client
	sessKey := app.getSessKey()
	postURL := "https://lms.iiitb.ac.in/moodle/lib/ajax/service.php?sesskey=" + sessKey.SessKey + "&info=core_course_get_enrolled_courses_by_timeline_classification"
	formData := []byte(`[{"index":0,"methodname":"core_course_get_enrolled_courses_by_timeline_classification","args":{"offset":0,"limit":96,"classification":"all","sort":"ul.timeaccess desc"}}]`)
	req, _ := http.NewRequest("POST", postURL, bytes.NewBuffer(formData))
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:81.0) Gecko/20100101 Firefox/81.0")
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
	type CourseResult []struct {
		Data struct {
			Courses []Course `json:"courses"`
		} `json:"data"`
	}

	var result CourseResult

	err := json.Unmarshal(postdocument, &result)
	if err != nil {
		log.Fatalln("Error Reading json Data: ", err)
	}
	return result[0].Data.Courses
}

func (app *App) getEvents() []Event {
	var Events []Event
	client := app.Client
	sessKey := app.getSessKey()
	postURL := "https://lms.iiitb.ac.in/moodle/lib/ajax/service.php?sesskey=" + sessKey.SessKey + "&info=core_calendar_get_action_events_by_timesort"
	formData := []byte(`[{"index":0,"methodname":"core_calendar_get_action_events_by_timesort","args":{"limitnum":6,"timesortfrom":` + strconv.FormatInt(time.Now().Unix(), 10) + `,"timesortto":` + strconv.FormatInt(time.Now().AddDate(0, 6, 0).Unix(), 10) + `,"limittononsuspendedevents":true}}]`)
	req, _ := http.NewRequest("POST", postURL, bytes.NewBuffer(formData))
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:81.0) Gecko/20100101 Firefox/81.0")
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	response, err := client.Do(req)
	if err != nil {
		log.Fatalln("Error fetching Post response. ", err)
	}

	defer response.Body.Close()
	postdocument, _ := ioutil.ReadAll(response.Body)
	type actionEventResult []struct {
		Error bool `json:"error"`
		Data  struct {
			Events []struct {
				Name        string `json:"name"`
				Description string `json:"description"`
				// Modulename string `json:"modulename"`
				Instance     int    `json:"instance"`
				Eventtype    string `json:"eventtype"`
				Timestart    int    `json:"timestart"`
				Timeduration int    `json:"timeduration"`
				Timesort     int    `json:"timesort"`
				Visible      int    `json:"visible"`
				Timemodified int    `json:"timemodified"`
				C            struct {
					ID int `json:"id"`
					// Fullname string `json:"fullname"`
					// Shortname string `json:"shortname"`
					// Startdate int `json:"startdate"`
					// Enddate int `json:"enddate"`
					// Fullnamedisplay string `json:"fullnamedisplay"`
					// Viewurl string `json:"viewurl"`
					// Isfavourite bool `json:"isfavourite"`
					// Hidden bool `json:"hidden"`
				} `json:"course"`
				Formattedtime string `json:"formattedtime"`
				URL           string `json:"url"`
				Action        struct {
					Name          string `json:"name"`
					URL           string `json:"url"`
					Itemcount     int    `json:"itemcount"`
					Actionable    bool   `json:"actionable"`
					Showitemcount bool   `json:"showitemcount"`
				} `json:"action"`
			} `json:"events"`
			Firstid int `json:"firstid"`
			Lastid  int `json:"lastid"`
		} `json:"data"`
	}
	var result actionEventResult
	err = json.Unmarshal(postdocument, &result)
	if err != nil {
		log.Fatalln("Error Reading json Data: ", err)
	}
	for _, events := range result[0].Data.Events {
		formtime := events.Formattedtime[88:108] + " " + events.Formattedtime[113:]
		if events.Formattedtime[107:108] == "<" {
			formtime = events.Formattedtime[88:107] + " " + events.Formattedtime[112:]
		}
		event := Event{
			Name:          events.Name,
			CourseID:      events.C.ID,
			Description:   events.Description,
			Instance:      events.Instance,
			Eventtype:     events.Eventtype,
			Timestart:     events.Timestart,
			Timeduration:  events.Timeduration,
			Timesort:      events.Timesort,
			Visible:       events.Visible,
			Timemodified:  events.Timemodified,
			Formattedtime: formtime,
			URL:           events.URL,
		}
		Events = append(Events, event)
	}
	return Events
}
func (app *App) getAnnouncements() []Announcement {
	var Announcements []Announcement
	for _, course := range app.getCourses() {
		if !course.Isfavorite {
			continue
		}
		coursePageURL := baseURL + "course/view.php?id=" + strconv.FormatInt(int64(course.CourseID), 10)
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
			name := s.Find(".name").Text()
			date := s.Find(".date").Text()
			info := s.Find(".info").Text()
			announcement := Announcement{
				Name:     name,
				Date:     date,
				Info:     info,
				CourseID: course.CourseID,
			}
			Announcements = append(Announcements, announcement)
		})
	}
	return Announcements
}
