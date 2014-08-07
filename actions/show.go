package actions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Zaibon/ircbot"
)

const (
	trakUrl = "http://api.trakt.tv/show/summary.json/ecba50561f7e607edc40ba85b1028bcd/%s"
)

type show struct {
	AirDay        string   `json:"air_day"`
	AirDayUtc     string   `json:"air_day_utc"`
	AirTime       string   `json:"air_time"`
	AirTimeUtc    string   `json:"air_time_utc"`
	Certification string   `json:"certification"`
	Country       string   `json:"country"`
	FirstAired    float64  `json:"first_aired"`
	FirstAiredIso string   `json:"first_aired_iso"`
	FirstAiredUtc float64  `json:"first_aired_utc"`
	Genres        []string `json:"genres"`
	Images        struct {
		Banner string `json:"banner"`
		Fanart string `json:"fanart"`
		Poster string `json:"poster"`
	} `json:"images"`
	ImdbID      string  `json:"imdb_id"`
	LastUpdated float64 `json:"last_updated"`
	Network     string  `json:"network"`
	Overview    string  `json:"overview"`
	People      struct {
		Actors []struct {
			Character string `json:"character"`
			Images    struct {
				Headshot string `json:"headshot"`
			} `json:"images"`
			Name string `json:"name"`
		} `json:"actors"`
	} `json:"people"`
	Poster  string `json:"poster"`
	Ratings struct {
		Hated      float64 `json:"hated"`
		Loved      float64 `json:"loved"`
		Percentage float64 `json:"percentage"`
		Votes      float64 `json:"votes"`
	} `json:"ratings"`
	Runtime float64 `json:"runtime"`
	Stats   struct {
		Checkins         float64 `json:"checkins"`
		CheckinsUnique   float64 `json:"checkins_unique"`
		Collection       float64 `json:"collection"`
		CollectionUnique float64 `json:"collection_unique"`
		Plays            float64 `json:"plays"`
		Scrobbles        float64 `json:"scrobbles"`
		ScrobblesUnique  float64 `json:"scrobbles_unique"`
		Watchers         float64 `json:"watchers"`
	} `json:"stats"`
	Status      string `json:"status"`
	Title       string `json:"title"`
	TopEpisodes []struct {
		FirstAired    float64 `json:"first_aired"`
		FirstAiredIso string  `json:"first_aired_iso"`
		FirstAiredUtc float64 `json:"first_aired_utc"`
		Number        float64 `json:"number"`
		Plays         float64 `json:"plays"`
		Season        float64 `json:"season"`
		Title         string  `json:"title"`
		URL           string  `json:"url"`
	} `json:"top_episodes"`
	TopWatchers []struct {
		About     string  `json:"about"`
		Age       string  `json:"age"`
		Avatar    string  `json:"avatar"`
		FullName  string  `json:"full_name"`
		Gender    string  `json:"gender"`
		Joined    float64 `json:"joined"`
		Location  string  `json:"location"`
		Plays     float64 `json:"plays"`
		Protected bool    `json:"protected"`
		URL       string  `json:"url"`
		Username  string  `json:"username"`
	} `json:"top_watchers"`
	TvdbID   float64 `json:"tvdb_id"`
	TvrageID float64 `json:"tvrage_id"`
	URL      string  `json:"url"`
	Year     float64 `json:"year"`
}

type ShowSummary struct{}

func (s *ShowSummary) Command() []string {
	return []string{".show"}
}

func (s *ShowSummary) Usage() string {
	return ".show <title> : display information about show"
}

func (s *ShowSummary) Do(b *ircbot.IrcBot, m *ircbot.IrcMsg) {
	if len(m.Args) < 2 {
		b.Say(m.Channel, "et si tu me donnais le titer du show ?")
		return
	}

	title := strings.Replace(strings.Join(m.Args[1:], " "), " ", "-", -1)
	resp, err := http.Get(fmt.Sprintf(trakUrl, title))
	if err != nil {
		b.Say(m.Channel, fmt.Sprintf("error : %s", err.Error()))
		return
	}
	if resp.StatusCode != 200 {
		switch resp.StatusCode {
		case http.StatusNotFound:
			b.Say(m.Channel, "not found")
		default:
			b.Say(m.Channel, fmt.Sprintf("error occurs :%s", resp.Status))
		}
		return
	}

	show := show{}
	if err := json.NewDecoder(resp.Body).Decode(&show); err != nil {
		b.Say(m.Channel, fmt.Sprintf("error : %s", err.Error()))
		return
	}

	b.Say(m.Channel, fmt.Sprintf("Title : %s", show.Title))
	b.Say(m.Channel, fmt.Sprintf("Beginning year : %.0f", show.Year))
	b.Say(m.Channel, fmt.Sprintf("Status : %s", show.Status))
	b.Say(m.Channel, fmt.Sprintf("Overview : %s", show.Overview))
	b.Say(m.Channel, fmt.Sprintf("Poster URL : %s", show.Poster))
}
