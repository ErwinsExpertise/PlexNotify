package handlers

import (
	"bufio"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// Log is created to hold the individual lines of the log file
// LogLines is a slice of type string and will hold the individual lines
// inside the log file.
type Log struct {
	LogLines []string
}

var (
	creds string
	auth  string
)

func init() {
	creds = strings.Join([]string{os.Getenv("UNAME"), os.Getenv("PWORD")}, ":")
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

func loginPage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("html/login.html")
	if err != nil {
		log.Println(err)
	}

	t.Execute(w, nil)
}

func logPage(w http.ResponseWriter, r *http.Request) {
	var pLog Log

	pLog.buildLogs(openLogs())

	t, err := template.ParseFiles("html/activity.html")
	if err != nil {
		log.Println(err)
	}

	t.Execute(w, &pLog)
}

func openLogs() *os.File {
	fil, err := os.Open("/var/log/plexnotify/notify.log")
	if err != nil {
		log.Println(err)
	}
	return fil
}

func (l *Log) buildLogs(fil *os.File) {
	scanner := bufio.NewScanner(fil)

	for scanner.Scan() {
		l.LogLines = append(l.LogLines, scanner.Text())
	}
}

func findTop(lines []string) map[string]int {
	top := make(map[string]int)
	for _, lin := range lines {
		if strings.Contains(lin, "media.play") {
			tmp := strings.Split(lin, "|")
			title := tmp[2][8:]
			top[title]++
		}
	}
	return top
}

func checkCookie(r *http.Request) (bool, string) {
	cook, _ := r.Cookie("X-Auth-Token")
	if cook != nil {
		return true, cook.Value
	}
	return false, ""

}

func setCookie(w http.ResponseWriter) {
	expiration := time.Now().Add(1 * 24 * time.Hour)
	cookie := http.Cookie{Name: "X-Auth-Token", Value: "plex-auth", Expires: expiration}
	http.SetCookie(w, &cookie)
}

func unAuthorized(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte("You are not authorized to access this resource."))
	log.Printf("Invalid Authorization Attempt! IP: %s", r.RemoteAddr)
}

func login(w http.ResponseWriter, r *http.Request) string {
	err := r.ParseForm()
	CheckErr(err)

	user := r.FormValue("username")
	pass := r.FormValue("password")
	up := strings.Join([]string{user, pass}, ":")

	return up
}
