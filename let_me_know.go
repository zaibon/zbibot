package zbibot

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Zaibon/ircbot"
	db "github.com/Zaibon/ircbot/database"
)

const (
	apiRoot = "http://lmk.zaibon.be/"
)

// var tokens map[string]string

type LetMeKnow struct {
	dbConn *db.DB

	tokens map[string]string
}

func NewLetMeKnow(bot *ircbot.IrcBot) *LetMeKnow {
	conn, err := bot.DBConnection()
	if err != nil {
		panic(err)
	}

	initDB(conn)
	return &LetMeKnow{
		dbConn: conn,
		tokens: map[string]string{},
	}
}

func initDB(db *db.DB) {
	if err := db.Exec(`CREATE TABLE IF NOT EXISTS lmk_tokens(
		id INTEGER CONSTRAINT line_PK PRIMARY KEY,
		prefix string,
		token,
		timestamp INTEGER)`); err != nil {

		panic(err)
	}
}

func (l *LetMeKnow) Command() []string {
	return []string{".lmk"}
}

func (l *LetMeKnow) Usage() string {
	return fmt.Sprintf(".lmk list|add")
}

func (l *LetMeKnow) Do(b *ircbot.IrcBot, msg *ircbot.IrcMsg) {
	if len(msg.Trailing) < 2 {
		b.Say(msg.Channel(), "shows cmd:")
		b.Say(msg.Channel(), "list | search :title | ep :season :episode :title | add :title")
		b.Say(msg.Channel(), "users cmd:")
		b.Say(msg.Channel(), "sign_up :mail :username :password | sign_in :username :password")

		return
	}

	cmd := msg.Trailing[1]

	//check if token is available if the cmd requiest it
	var (
		token string
		err   error
	)
	if cmd != "signup" && cmd != "signin" {
		token, err = l.getToken(msg.Prefix)
		if err != nil {
			if err == io.EOF {
				err = fmt.Errorf("token non disponible, veuillez vous authentifier")
			}
			b.Say(msg.Channel(), err.Error())
			return
		}
	}

	switch cmd {
	case "list":
		l.doShowsList(b, msg, token)

	case "search":
		l.doShowsSearch(b, msg, token)

	case "ep":
		l.doShowsSearchEp(b, msg, token)

	case "add":
		l.doShowsAdd(b, msg, token)

	case "signup":
		l.doUsersSignUp(b, msg)

	case "signin":
		l.doUsersSignIn(b, msg)
	}
}

type APIResp struct {
	Status  string          `json:"status"`
	Payload json.RawMessage `json:"msg"`
}

func (l *LetMeKnow) doShowsAdd(b *ircbot.IrcBot, msg *ircbot.IrcMsg, token string) error {
	if len(msg.Trailing) < 3 {
		b.Say(msg.Channel(), "missing parameter")
		return nil
	}

	title := url.QueryEscape(strings.Join(msg.Trailing[2:], " "))
	apiURL := fmt.Sprintf("%s%s/%s?token=%s", apiRoot, "shows/add", title, token)
	resp, err := http.Post(apiURL, "text/html", nil)
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

type showsList []struct {
	Title string `json:"title"`
}

func (l *LetMeKnow) doShowsList(b *ircbot.IrcBot, msg *ircbot.IrcMsg, token string) error {
	url := fmt.Sprintf("%s%s?token=%s", apiRoot, "shows/list", token)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("error : ", err)
		return err
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
	return err
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

func (l *LetMeKnow) doShowsSearch(b *ircbot.IrcBot, msg *ircbot.IrcMsg, token string) error {
	if len(msg.Trailing) < 3 {
		b.Say(msg.Channel(), "missing parameter")
		return nil
	}

	title := url.QueryEscape(strings.Join(msg.Trailing[2:], " "))
	apiURL := fmt.Sprintf("%s%s/%s?token=%s", apiRoot, "shows/search", title, token)
	resp, err := http.Get(apiURL)
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

func (l *LetMeKnow) doShowsSearchEp(b *ircbot.IrcBot, msg *ircbot.IrcMsg, token string) error {
	if len(msg.Trailing) < 5 {
		b.Say(msg.Channel(), "missing parameter")
		return nil
	}

	season := msg.Trailing[2]
	number := msg.Trailing[3]
	title := url.QueryEscape(strings.Join(msg.Trailing[4:], " "))
	apiURL := fmt.Sprintf("%s%s/%s/%s/%s/%s?token=%s", apiRoot, "shows/search", title, "episodes", season, number, token)
	resp, err := http.Get(apiURL)
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

func (l *LetMeKnow) doUsersSignUp(b *ircbot.IrcBot, msg *ircbot.IrcMsg) error {
	if len(msg.Trailing) < 5 {
		b.Say(msg.Channel(), "missing parameter")
		return nil
	}

	mail := msg.Trailing[2]
	username := msg.Trailing[3]
	password := msg.Trailing[4]

	url := fmt.Sprintf("%s%s?mail=%s&username=%s&password=%s", apiRoot, "users/sign_up", mail, username, password)
	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		fmt.Println("sign_up error get :", err)
		return err
	}
	defer resp.Body.Close()

	apiResp := &APIResp{}
	if err := json.NewDecoder(resp.Body).Decode(apiResp); err != nil {
		fmt.Println("sign_up error decode :", err)
		return err
	}

	if err := checkApiRespError(apiResp, b, msg); err != nil {
		return err
	}

	if apiResp.Status == "ok" {
		resp := ""
		if err := json.Unmarshal(apiResp.Payload, &resp); err != nil {
			fmt.Println("sign_up decode error :", err)
			return err
		}
		b.Say(msg.Channel(), resp)
	}
	return nil
}

func (l *LetMeKnow) doUsersSignIn(b *ircbot.IrcBot, msg *ircbot.IrcMsg) error {
	if len(msg.Trailing) < 4 {
		b.Say(msg.Channel(), "missing parameter")
		return nil
	}

	username := msg.Trailing[2]
	password := msg.Trailing[3]

	url := fmt.Sprintf("%s%s?username=%s&password=%s", apiRoot, "users/sign_in", username, password)
	fmt.Println("DEBUG :", url)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("sign_in error get :", err)
		return err
	}
	defer resp.Body.Close()

	apiResp := &APIResp{}
	if err := json.NewDecoder(resp.Body).Decode(apiResp); err != nil {
		fmt.Println("sign_in error decode :", err)
		return err
	}

	if err := checkApiRespError(apiResp, b, msg); err != nil {
		return err
	}

	if apiResp.Status == "ok" {
		token := ""
		if err := json.Unmarshal(apiResp.Payload, &token); err != nil {
			fmt.Println("sign_in decode error :", err)
			return err
		}

		sql := "INSERT INTO lmk_tokens(prefix,token,timestamp) VALUES ($prefix,$token,$timestamp)"
		if err := l.dbConn.Exec(sql, msg.Prefix, token, time.Now()); err != nil {
			b.Say(msg.Nick(), err.Error())
		}

		//save token in memory for next use
		l.tokens[msg.Prefix] = token
		b.Say(msg.Channel(), "token send in query")
		b.Say(msg.Nick(), fmt.Sprintf("your token : %s", token))
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

// func isAuthentificated(b *ircbot.IrcBot, msg *ircbot.IrcMsg) bool {
// 	if _, ok := tokens[msg.Prefix]; !ok {
// 		b.Say(msg.Channel(), "vous devez vous authentifier avant de faire cette commande")
// 		return false
// 	}
// 	return true
// }

func (l *LetMeKnow) getToken(prefix string) (string, error) {
	//test if token is in memory
	token, ok := l.tokens[prefix]
	if ok {
		return token, nil
	}

	//if not, try to retreive from database
	sql := "SELECT token FROM lmk_tokens WHERE prefix=$prefix"
	stmt, err := l.dbConn.Query(sql, prefix)
	if err != nil {
		return "", err
	}
	err = stmt.Scan(&token)
	if err == nil {
		l.tokens[prefix] = token
	}

	if err := stmt.Close(); err != nil {
		fmt.Printf("ERROR close statement query : %s\n", err)
		return "", err
	}
	return token, nil
}
