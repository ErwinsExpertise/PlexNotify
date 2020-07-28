package handlers

import (
	"encoding/json"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

// Log is created to hold the individual lines of the log file
// LogLines is a slice of type string and will hold the individual lines
// inside the log file.
type Log struct {
	Events []Event
}

type Event struct {
	Timestamp string
	Ev        string
	User      string
	Title     string
	IP        string
}

var eventFile *os.File

func init() {
	eventFile, err := os.OpenFile("activity.json", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		log.Panic(err)
	}

	defer eventFile.Close()
}

// ActivityHandler is the primary route for /activity
// route accepts both GET and POST requests
// This route is password protected to prevent indexing on search engines as well as
// preventing data leaks to scaping bots.
// POST requests are used to authenticate user
// GET  requests are used to display the content of the logs
func ActivityHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	switch r.Method {
	case "POST":
		auth = login(w, r)

		if strings.Compare(auth, creds) == 0 {
			setCookie(w)
			logPage(w, r)
			return
		}
		unAuthorized(w, r)
		return

	case "GET":
		bo, val := checkCookie(r)
		if bo == true {
			if strings.Compare(val, "plex-auth") == 0 {
				logPage(w, r)
				return
			}
			loginPage(w, r)
			return

		}
		loginPage(w, r)
		return

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

}

func activityPage(w http.ResponseWriter, r *http.Request) {
	var pLog Log

	pLog.buildActivities()

	t, err := template.ParseFiles("html/activity.html")
	if err != nil {
		log.Println(err)
	}

	t.Execute(w, &pLog)
}

func openActivityLog() *os.File {
	fil, err := os.Open("activity.json")
	if err != nil {
		log.Println(err)
	}
	return fil
}

func (l *Log) buildActivities() {
	evFile, err := ioutil.ReadFile("activity.json")
	CheckErr(err)
	json.Unmarshal(evFile, &l)
}

func AppendActivity(time, event, user, title, ip string) {

	ev := Event{
		Timestamp: time,
		Ev:        event,
		User:      user,
		Title:     title,
		IP:        ip,
	}

	evJson, err := json.Marshal(&ev)
	CheckErr(err)

	f, err := os.OpenFile("activity.json", os.O_APPEND|os.O_WRONLY, 0666)
	CheckErr(err)

	_, err = io.WriteString(f, string(evJson))
	CheckErr(err)

}
