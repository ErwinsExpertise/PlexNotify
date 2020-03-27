package handlers

import (
	"reflect"
	"testing"
)

func TestBuildField(t *testing.T) {
	t.Run("creating new Field with given values", func(t *testing.T) {
		want := Field{
			Name:   "Title",
			Value:  "Testing",
			Inline: false,
		}

		got := BuildField("Title", "Testing", false)

		if want != got {
			t.Errorf("Wanted:%v Got:%v\n", want, got)
		}

	})

	t.Run("creating new Field with empty values", func(t *testing.T) {
		want := Field{}

		got := BuildField("", "", false)

		if want != got {
			t.Errorf("Wanted:%v Got:%v\n", want, got)
		}

	})
}

func TestBuildEmbed(t *testing.T) {
	t.Run("build new embed object with given values", func(t *testing.T) {
		var want, got Embed

		auth := makeAuthor("Test Server", "https://thumbnail.imgbin.com/19/3/22/imgbin-plex-media-server-computer-icons-media-player-tv-icon-Z5vd0y4dwWu0NNmcvuWCTNnAv_t.jpg")

		e := EventInfo{
			Server: ServerInfo{
				Title: "Test Server",
			},
		}

		want = Embed{
			Title:  "Test content",
			Color:  "4889028",
			Author: auth,
		}

		got.BuildEmbed(&e, auth, "4889028", "Test content", "", e.Metadata["type"])

		res := reflect.DeepEqual(want, got)
		if res == false {
			t.Errorf("Wanted:%v Got:%v\n", want, got)
		}
	})

}
