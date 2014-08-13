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
	if len(msg.Trailing) < 2 {
		b.Say(msg.Channel(), "list | search :title | ep :season :episode :title | add :title")
		return
	}

	cmd := msg.Trailing[1]
	switch cmd {
	case "list":
		doShowsList(b, msg)

	case "search":
		doShowsSearch(b, msg)

	case "ep":
		doShowsSearchEp(b, msg)

	case "add":
		doShowsAdd(b, msg)
	}
}

func doShowsAdd(b *ircbot.IrcBot, msg *ircbot.IrcMsg) error {
	if len(msg.Trailing) < 3 {
		b.Say(msg.Channel(), "donne le nom de la série")
		return nil
	}
	url := apiRoot + "shows/add/" + strings.Join(msg.Trailing[2:], "-")
	resp, err := http.Post(url, "text/html", nil)
	if err != nil {
		fmt.Println("error post: ", err)
		return err
	}
	defer resp.Body.Close()

	apiResp := &APIResp{}
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		fmt.Println("error decoding : ", err)
		return err
	}

	if err := checkApiRespError(apiResp, b, msg); err != nil {
		fmt.Println("error decode :", err)
		return err
	}

	if apiResp.Status == "ok" {
		var addMsg string
		if err := json.Unmarshal(apiResp.Payload, &addMsg); err != nil {
			fmt.Println("error decode : ", err)
			return err
		}
		b.Say(msg.Channel(), addMsg)
	}
	return nil
}

type APIResp struct {
	Status  string          `json:"status"`
	Payload json.RawMessage `json:"msg"`
}

type showsList []struct {
	Title string `json:"title"`
}

func doShowsList(b *ircbot.IrcBot, msg *ircbot.IrcMsg) error {
	url := apiRoot + "shows/list"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("error : ", err)
		return nil
	}
	defer resp.Body.Close()

	apiResp := &APIResp{}
	if err := json.NewDecoder(resp.Body).Decode(apiResp); err != nil {
		fmt.Println("list error : ", err)
		return err
	}
	if err := checkApiRespError(apiResp, b, msg); err != nil {
		return err
	}

	if apiResp.Status == "ok" {
		shows := showsList{}
		if err := json.Unmarshal(apiResp.Payload, &shows); err != nil {
			fmt.Println("list error :", err)
			return err
		}
		for _, title := range shows {
			b.Say(msg.Channel(), title.Title)
		}
	}
	return nil
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

func doShowsSearch(b *ircbot.IrcBot, msg *ircbot.IrcMsg) error {
	if len(msg.Trailing) < 3 {
		b.Say(msg.Channel(), "que cherche tu ?")
		return nil
	}
	url := apiRoot + "shows/search/" + strings.Join(msg.Trailing[2:], " ")
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("search error get :", err)
		return err
	}
	defer resp.Body.Close()

	apiResp := &APIResp{}
	if err := json.NewDecoder(resp.Body).Decode(apiResp); err != nil {
		fmt.Println("search error decode :", err)
		return err
	}

	if err := checkApiRespError(apiResp, b, msg); err != nil {
		return err
	}

	if apiResp.Status == "ok" {
		resp := showsSearchResp{}
		if err := json.Unmarshal(apiResp.Payload, &resp); err != nil {
			fmt.Println("search decode error :", err)
			return err
		}

		for _, show := range resp {
			b.Say(msg.Channel(), fmt.Sprintf("title : %s", show.Title))
			b.Say(msg.Channel(), fmt.Sprintf("Overview : %s", show.Overview))
		}
	}
	return nil
}

type showsSearchEpResp struct {
	ID struct {
		_Oid string `json:"$oid"`
	} `json:"_id"`
	CreatedAt   string      `json:"created_at"`
	DownloadURL interface{} `json:"download_url"`
	Note        float64     `json:"note"`
	Number      float64     `json:"number"`
	Overview    string      `json:"overview"`
	ReleasedOn  string      `json:"released_on"`
	Season      float64     `json:"season"`
	ShowID      struct {
		_Oid string `json:"$oid"`
	} `json:"show_id"`
	Title     string `json:"title"`
	UpdatedAt string `json:"updated_at"`
}

func doShowsSearchEp(b *ircbot.IrcBot, msg *ircbot.IrcMsg) error {
	if len(msg.Trailing) < 5 {
		b.Say(msg.Channel(), "pas assez de paramètres")
		return nil
	}
	season := msg.Trailing[2]
	number := msg.Trailing[3]
	title := strings.Join(msg.Trailing[4:], "-")
	url := fmt.Sprintf("%s%s/%s/%s/%s/%s", apiRoot, "shows/search", title, "episodes", season, number)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("search error get :", err)
		return err
	}
	defer resp.Body.Close()

	apiResp := &APIResp{}
	if err := json.NewDecoder(resp.Body).Decode(apiResp); err != nil {
		fmt.Println("search error decode :", err)
		return err
	}

	if err := checkApiRespError(apiResp, b, msg); err != nil {
		return err
	}

	if apiResp.Status == "ok" {
		resp := showsSearchEpResp{}
		if err := json.Unmarshal(apiResp.Payload, &resp); err != nil {
			fmt.Println("search decode error :", err)
			return err
		}

		b.Say(msg.Channel(), fmt.Sprintf("title : %s", resp.Title))
		b.Say(msg.Channel(), fmt.Sprintf("Season : %.0f", resp.Season))
		b.Say(msg.Channel(), fmt.Sprintf("Episode : %.0f", resp.Number))
		b.Say(msg.Channel(), fmt.Sprintf("Overview : %s", resp.Overview))

	}
	return nil
}

func checkApiRespError(apiResp *APIResp, b *ircbot.IrcBot, m *ircbot.IrcMsg) error {
	if apiResp.Status == "error" {
		var errMsg string
		if err := json.Unmarshal(apiResp.Payload, &errMsg); err != nil {
			fmt.Println("error decode :", err)
			return err
		}
		b.Say(m.Channel(), errMsg)
	}
	return nil
}
