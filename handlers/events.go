package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ErwinsExpertise/PlexNotify/metrics"
)

type EventInfo struct {
	Event    string            `json:"event, omitempty"`
	User     string            `json:"user, omitempty"`
	Owner    bool              `json:"owner, omitempty`
	Account  AccountInfo       `json:"Account, omitempty"`
	Server   ServerInfo        `json:"Server, omitempty"`
	Player   PlayerInfo        `json:"Player, omitempty"`
	Metadata map[string]string `json:"metadata, omitempty"`
}

type AccountInfo struct {
	ID    string `json:"id, omitempty"`
	Thumb string `json:"thumb, omitempty"`
	Title string `json:"title, omitempty"`
}

type ServerInfo struct {
	Title string `json:"title, omitempty"`
	UUID  string `json:"uuid, omitempty"`
}

type PlayerInfo struct {
	Local         string `json:"local, omitempty"`
	PublicAddress string `json:"publicAddress, omitempty"`
	Title         string `json:"title, omitempty"`
	UUID          string `json:"uuid, omitempty"`
}

type Payload struct {
	Content  string  `json:"content, omitempty"`
	Username string  `json:"username, omitempty"`
	Avatar   string  `json:"avatar_url, omitempty"`
	TTS      bool    `json:"tts, omitempty"`
	Embeds   []Embed `json:"embeds, omitempty"`
}

type Embed struct {
	Color       string            `json:"color, omitempty"`
	Author      map[string]string `json:"author, omitempty"`
	Title       string            `json:"title, omitempty"`
	URL         string            `json:"url, omitempty"`
	Description string            `json:"description, omitempty"`
	Fields      []Field           `json:"fields, omitempty"`
}

type Field struct {
	Name   string `json:"name, omitempty"`
	Value  string `json:"value, omitempty"`
	Inline bool   `json:"inline, omitempty"`
}

var (
	discordURL string
)

func init() {
	discordURL = os.Getenv("DISCORDURL")
}

func EventHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var eventLoad EventInfo
	defer r.Body.Close()

	payload := ExtractPayload(eventLoad, "payload", r)
	json.Unmarshal(payload, &eventLoad)

	resp, err := eventLoad.identifyEvent()
	if err != nil {
		log.Println(err)
		log.Printf("Payload: %v\n", eventLoad)
	}

	//Checks if event type is empty
	if strings.Compare(resp, "") == 0 {
		w.WriteHeader(500)
		return
	}

	metrics.EventCollector(int64(time.Since(start) / 1000000))

	log.Printf("Processing event type: %s | User: %s | Title: %s | IP: %s", resp, eventLoad.Account.Title, eventLoad.CheckTitle(), eventLoad.Player.PublicAddress)
	AppendActivity(time.Now().Format("2006-01-02 15:04:05"), resp, eventLoad.Account.Title, eventLoad.CheckTitle(), eventLoad.Player.PublicAddress)
	w.WriteHeader(200)

}

func (e *EventInfo) identifyEvent() (string, error) {
	var payload string
	switch e.Event {
	case "library.on.deck": // A new item is added that appears in the user’s On Deck. A poster is also attached to this event.
		return e.Event, nil
	case "library.new": // A new item is added to a library to which the user has access. A poster is also attached to this event.
		return e.Event, nil
	case "media.pause": // Media playback pauses.
		payload = fmt.Sprintf("%s paused %s", e.Account.Title, e.CheckTitle())
		e.SendHook(payload)
		return e.Event, nil
	case "media.play": // Media starts playing. An appropriate poster is attached.
		payload = fmt.Sprintf("%s started playing %s", e.Account.Title, e.CheckTitle())
		e.SendHook(payload)
		return e.Event, nil
	case "media.rate": // Media is rated. A poster is also attached to this event.
		payload = fmt.Sprintf("%s rated %s", e.Account.Title, e.CheckTitle())
		e.SendHook(payload)
		return e.Event, nil
	case "media.resume": // Media playback resumes.
		payload = fmt.Sprintf("%s resumed %s", e.Account.Title, e.CheckTitle())
		e.SendHook(payload)
		return e.Event, nil
	case "media.scrobble": // Media is viewed (played past the 90% mark).
		payload = fmt.Sprintf("%s has finished %s", e.Account.Title, e.CheckTitle())
		e.SendHook(payload)
		return e.Event, nil
	case "media.stop": // Media playback stops.
		payload = fmt.Sprintf("%s stopped %s", e.Account.Title, e.CheckTitle())
		e.SendHook(payload)
		return e.Event, nil
	case "admin.database.backup": // A database backup is completed successfully via Scheduled Tasks.
		payload = fmt.Sprintf("Database backup completed successfully.")
		e.SendHook(payload)
		return e.Event, nil
	case "admin.database.corrupted": // Corruption is detected in the server database.
		payload = fmt.Sprintf("Server database corrupt!")
		e.SendHook(payload)
		return e.Event, nil
	case "device.new": // A device accesses the owner’s server for any reason, which may come from background connection testing and doesn’t necessarily indicate active browsing or playback.
		payload = fmt.Sprintf("New Device ( %s ) accessed owner server from %s", e.Player.Title, e.Player.PublicAddress)
		e.SendHook(payload)
		return e.Event, nil
	case "playback.started": //Playback is started by a shared user for the server. A poster is also attached to this event.
		payload = fmt.Sprintf("%s started playback for %s", e.Account.Title, e.CheckTitle())
		e.SendHook(payload)
		return e.Event, nil
	default:
		payload = "Unknown event attempted"
		e.SendHook(payload)
		return "", errors.New("Unknown event type( " + e.Event + " ) recieved from " + e.Account.Title + "!")

	}

}

func (e *EventInfo) SendHook(content string) {

	payload := e.BuildPayload(content)

	requestBody, err := json.Marshal(payload)
	CheckErr(err)

	conType := "application/json"
	_, err = http.Post(discordURL, conType, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Println(err)
		log.Printf("RequestBody: %v\n", string(requestBody))
	}
}

func (e *EventInfo) BuildPayload(content string) Payload {
	var payload Payload
	var emb Embed

	field1 := BuildField("Player", e.Player.Title, true)
	field2 := BuildField("IP Address", e.Player.PublicAddress, true)

	auth := makeAuthor(e.Server.Title, "https://thumbnail.imgbin.com/19/3/22/imgbin-plex-media-server-computer-icons-media-player-tv-icon-Z5vd0y4dwWu0NNmcvuWCTNnAv_t.jpg")

	emb.BuildEmbed(e, auth, "4889028", content, "", e.Metadata["type"], field1, field2)

	payload.AddEmbed(emb)
	return payload

}

func (e *EventInfo) CheckTitle() string {
	if e.Metadata["grandparentTitle"] != "" {
		return fmt.Sprintf(e.Metadata["grandparentTitle"] + " - " + e.Metadata["title"])
	}
	return e.Metadata["title"]

}

func (p *Payload) AddEmbed(emb ...Embed) {
	for _, em := range emb {
		p.Embeds = append(p.Embeds, em)
	}
}

func (emb *Embed) AddFields(f Field) {
	emb.Fields = append(emb.Fields, f)
}

func (emb *Embed) BuildEmbed(e *EventInfo, author map[string]string, color, title, url, description string, fields ...Field) {
	emb.Author = make(map[string]string)
	for key, val := range author {
		emb.Author[key] = val
	}
	emb.Title = title
	emb.Color = color
	emb.Description = description
	emb.URL = url

	for _, f := range fields {
		emb.AddFields(f)
	}

}

func makeAuthor(name, icon_url string) map[string]string {
	m := make(map[string]string)
	m["name"] = name
	m["icon_url"] = icon_url

	return m
}

func ExtractPayload(sType interface{}, partName string, r *http.Request) []byte {

	mpr, err := r.MultipartReader()
	if err != nil {
		log.Println(err)
	}

	for {
		part, err := mpr.NextPart()
		if err == io.EOF {
			break
		}
		CheckErr(err)

		if part.FormName() == partName {
			decoder := json.NewDecoder(part)
			decoder.Decode(&sType)
			break
		}
	}

	bSlice, err := json.Marshal(&sType)
	CheckErr(err)

	return bSlice
}

func BuildField(name, value string, inline bool) Field {
	newField := Field{
		Name:   name,
		Value:  value,
		Inline: inline,
	}
	return newField
}

func CheckErr(err error) {
	if err != nil {
		log.Println(err)
	}
}
