package handlers

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	creds string
	auth  string
)

func init() {
	creds = strings.Join([]string{os.Getenv("UNAME"), os.Getenv("PWORD")}, ":")
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	loginPage(w, r)
}

func loginPage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("html/login.html")
	if err != nil {
		log.Println(err)
	}

	t.Execute(w, nil)
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
	cookie := http.Cookie{Name: "X-Auth-Token", Value: "plex-auth", Expires: expiration, HttpOnly: true}
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
