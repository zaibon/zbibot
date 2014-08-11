package actions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Zaibon/ircbot"
)

const (
	apiRoot = "http://zaibon.be:3010/"
)

type LetMeKnow struct{}

func (l *LetMeKnow) Command() []string {
	return []string{".lmk"}
}

func (l *LetMeKnow) Usage() string {
	return fmt.Sprintf(".lmk list|add")
}

func (l *LetMeKnow) Do(b *ircbot.IrcBot, msg *ircbot.IrcMsg) {
	if len(msg.Args) < 2 {
		b.Say(msg.Channel, "quelle commande?")
		return
	}

	cmd := msg.Args[1]
	apiResp := &APIResp{}
	switch cmd {
	case "list":
		url := apiRoot + "shows/list"
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("error : ", err)
			return
		}
		defer resp.Body.Close()

		if err := json.NewDecoder(resp.Body).Decode(apiResp); err != nil {
			fmt.Println("list error : ", err)
			return
		}
		if err := checkApiRespError(apiResp, b, msg); err != nil {
			return
		}

		if apiResp.Status == "ok" {
			shows := showsList{}
			if err := json.Unmarshal(apiResp.Payload, &shows); err != nil {
				fmt.Println("list error :", err)
				return
			}
			for _, title := range shows {
				b.Say(msg.Channel, title.Title)
			}
		}

	case "search":
		if len(msg.Args) < 3 {
			b.Say(msg.Channel, "que cherche tu ?")
			return
		}
		url := apiRoot + "shows/search/" + strings.Join(msg.Args[2:], " ")
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("search error get :", err)
			return
		}
		defer resp.Body.Close()

		if err := json.NewDecoder(resp.Body).Decode(apiResp); err != nil {
			fmt.Println("search error decode :", err)
			return
		}

		if err := checkApiRespError(apiResp, b, msg); err != nil {
			return
		}

		if apiResp.Status == "ok" {
			resp := showsSearchResp{}
			if err := json.Unmarshal(apiResp.Payload, &resp); err != nil {
				fmt.Println("search decode error :", err)
				return
			}

			for _, show := range resp {
				b.Say(msg.Channel, fmt.Sprintf("title : %s", show.Title))
				b.Say(msg.Channel, fmt.Sprintf("Overview : %s", show.Overview))
			}
		}

	case "add":
		if len(msg.Args) < 3 {
			b.Say(msg.Channel, "donne le nom de la série")
			return
		}
		url := apiRoot + "shows/add/" + strings.Join(msg.Args[2:], "-")
		resp, err := http.Post(url, "text/html", nil)
		if err != nil {
			fmt.Println("error post: ", err)
			return
		}
		defer resp.Body.Close()

		if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
			fmt.Println("error decoding : ", err)
			return
		}

		if err := checkApiRespError(apiResp, b, msg); err != nil {
			fmt.Println("error decode :", err)
			return
		}

		if apiResp.Status == "ok" {
			var addMsg string
			if err := json.Unmarshal(apiResp.Payload, &addMsg); err != nil {
				fmt.Println("error decode : ", err)
				return
			}
			b.Say(msg.Channel, addMsg)
		}
	}
}

type APIResp struct {
	Status  string          `json:"status"`
	Payload json.RawMessage `json:"msg"`
}

type showsList []struct {
	Title string `json:"title"`
}

type showsSearchResp []struct {
	ID struct {
		ID string `json:"$oid"`
	} `json:"_id"`
	BannerURL string  `json:"banner_url"`
	BeginYear float64 `json:"begin_year"`
	CreatedAt string  `json:"created_at"`
	Overview  string  `json:"overview"`
	PosterURL string  `json:"poster_url"`
	Slug      string  `json:"slug"`
	Title     string  `json:"title"`
	UpdatedAt string  `json:"updated_at"`
}

type showsEp struct {
	ID struct {
		ID string `json:"$oid"`
	} `json:"_id"`
	CreatedAt   string      `json:"created_at"`
	DownloadURL interface{} `json:"download_url"`
	Note        float64     `json:"note"`
	Number      float64     `json:"number"`
	Overview    string      `json:"overview"`
	ReleasedOn  string      `json:"released_on"`
	Season      float64     `json:"season"`
	ShowID      struct {
		ID string `json:"$oid"`
	} `json:"show_id"`
	Title     string `json:"title"`
	UpdatedAt string `json:"updated_at"`
}

func checkApiRespError(apiResp *APIResp, b *ircbot.IrcBot, m *ircbot.IrcMsg) error {
	if apiResp.Status == "error" {
		var errMsg string
		if err := json.Unmarshal(apiResp.Payload, &errMsg); err != nil {
			fmt.Println("error decode :", err)
			return err
		}
		b.Say(m.Channel, errMsg)
	}
	return nil
}
