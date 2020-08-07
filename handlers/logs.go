package handlers

import (
	"bufio"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

type Logs struct {
	LogLines []string
}

var logPath string

// LogHandler is the primary route for /logs
// route accepts both GET and POST requests
// This route is password protected to prevent indexing on search engines as well as
// preventing data leaks to scaping bots.
// POST requests are used to authenticate user
// GET  requests are used to display the content of the logs
func LogHandler(w http.ResponseWriter, r *http.Request) {
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

func logPage(w http.ResponseWriter, r *http.Request) {
	var pLog Logs

	pLog.buildLogs(openLogs())

	t, err := template.ParseFiles("html/log.html")
	if err != nil {
		log.Println(err)
	}

	t.Execute(w, &pLog)
}

func openLogs() *os.File {
	fil, err := os.Open(logPath + "notify.log")
	if err != nil {
		log.Println(err)
	}
	return fil
}

func (l *Logs) buildLogs(fil *os.File) {
	scanner := bufio.NewScanner(fil)

	for scanner.Scan() {
		l.LogLines = append(l.LogLines, scanner.Text())
	}
}

func CreateLogs() string {
	if _, err := os.Stat("/var/log/plexnotify"); os.IsNotExist(err) {
		err := os.Mkdir("/var/log/plexnotify", 0755)
		if err != nil {
			logPath = ""
		} else {
			logPath = "/var/log/plexnotify/"
		}
	}

	return logPath

}
