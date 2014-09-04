package zbibot

import (
	"bytes"
	"crypto/tls"
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
	apiRoot = "https://lmk.hito.be"
)

// var tokens map[string]string

type LetMeKnow struct {
	dbConn *db.DB
	client *http.Client
	bot    *ircbot.IrcBot

	tokens map[string]string
}

func NewLetMeKnow(bot *ircbot.IrcBot) *LetMeKnow {
	conn, err := bot.DBConnection()
	if err != nil {
		panic(err)
	}

	initDB(conn)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	return &LetMeKnow{
		dbConn: conn,
		tokens: map[string]string{},
		bot:    bot,
		client: client,
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
		b.Say(msg.Channel(), "list | search :title | ep :season :episode :title | add :title | follow :title | unfollow :title | followed ?:username")
		b.Say(msg.Channel(), "users cmd:")
		b.Say(msg.Channel(), "signup :mail :username :password | signin :username :password | lost :email")

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

	case "follow":
		l.doFollowShow(b, msg, token)

	case "unfollow":
		l.doUnfollowShow(b, msg, token)

	case "followed":
		l.doUsersFollowed(b, msg, token)

	case "lost":
		l.doLostPassword(b, msg, token)

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

	title := strings.Join(msg.Trailing[2:], " ")
	apiURL := fmt.Sprintf("%s/%s", apiRoot, "shows")

	body := struct {
		Title string `json:"title"`
	}{
		title,
	}
	apiResp, err := l.post(apiURL, body, token, msg)
	if err != nil {
		return err
	}

	var addMsg string
	if err := json.Unmarshal(apiResp.Payload, &addMsg); err != nil {
		fmt.Println("error decode : ", err)
		return err
	}

	b.Say(msg.Channel(), addMsg)

	return nil
}

type showsList []struct {
	Title string `json:"title"`
}

func (l *LetMeKnow) doShowsList(b *ircbot.IrcBot, msg *ircbot.IrcMsg, token string) error {
	apiURL := fmt.Sprintf("%s/%s", apiRoot, "shows")

	apiResp, err := l.get(apiURL, token, msg)
	if err != nil {
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

	title := strings.Join(msg.Trailing[2:], " ")
	apiURL := fmt.Sprintf("%s/%s/%s", apiRoot, "shows", title)

	apiResp, err := l.get(apiURL, token, msg)
	if err != nil {
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
	title := strings.Join(msg.Trailing[4:], " ")
	apiURL := fmt.Sprintf("%s/%s/%s/%s/%s/%s", apiRoot, "shows", title, "episodes", season, number)

	apiResp, err := l.get(apiURL, token, msg)
	if err != nil {
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

	body := struct {
		Username string `json:"username"`
		Pass     string `json:"password"`
		Email    string `json:"email"`
	}{
		username,
		password,
		mail,
	}

	apiURL := fmt.Sprintf("%s/%s", apiRoot, "users/sign_up")
	apiResp, err := l.post(apiURL, body, "", msg)
	if apiResp != nil {
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

	apiURL := fmt.Sprintf("%s/%s?username=%s&password=%s", apiRoot, "users/sign_in", username, password)
	apiResp, err := l.get(apiURL, "", msg)
	if err != nil {
		return err
	}

	if apiResp.Status == "ok" {
		token := ""
		if err := json.Unmarshal(apiResp.Payload, &token); err != nil {
			fmt.Println("sign_in decode error :", err)
			return err
		}

		//this is ugly...sorry
		sql := "DELETE FROM lmk_tokens WHERE prefix=$prefix"
		if err := l.dbConn.Exec(sql, msg.Prefix); err != nil {
			b.Say(msg.Nick(), err.Error())
		}
		sql = "INSERT INTO lmk_tokens(prefix,token,timestamp) VALUES ($prefix,$token,$timestamp)"
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

func (l *LetMeKnow) doFollowShow(b *ircbot.IrcBot, msg *ircbot.IrcMsg, token string) error {
	if len(msg.Trailing) < 3 {
		b.Say(msg.Channel(), "missing parameter")
		return nil
	}

	title := strings.Join(msg.Trailing[2:], " ")
	body := struct {
		Title string `json:"title"`
	}{
		title,
	}

	apiURL := fmt.Sprintf("%s/%s/%s", apiRoot, "users", "follows")
	apiResp, err := l.post(apiURL, body, token, msg)
	if err != nil {
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

func (l *LetMeKnow) doUnfollowShow(b *ircbot.IrcBot, msg *ircbot.IrcMsg, token string) error {
	if len(msg.Trailing) < 3 {
		b.Say(msg.Channel(), "missing parameter")
		return nil
	}

	title := url.QueryEscape(strings.Join(msg.Trailing[2:], " "))
	apiURL := fmt.Sprintf("%s/%s/%s/%s", apiRoot, "users", "follows", title)

	apiResp, err := l.delete(apiURL, token, msg)
	if err != nil {
		return err
	}

	if apiResp.Status == "ok" {
		resp := ""
		if err := json.Unmarshal(apiResp.Payload, &resp); err != nil {
			fmt.Println("error decode : ", err)
			return err
		}

		b.Say(msg.Channel(), resp)

	}
	return nil
}

type followed []struct {
	Title string `json:"title"`
}

func (l *LetMeKnow) doUsersFollowed(b *ircbot.IrcBot, msg *ircbot.IrcMsg, token string) error {
	apiURL := ""

	if len(msg.Trailing) > 2 {
		username := msg.Trailing[2]
		apiURL = fmt.Sprintf("%s/%s/%s/%s", apiRoot, "users", "follows", username)
	} else {
		apiURL = fmt.Sprintf("%s/%s/%s", apiRoot, "users", "follows")
	}

	apiResp, err := l.get(apiURL, token, msg)
	if err != nil {
		return err
	}

	if apiResp.Status == "ok" {
		shows := showsList{}
		if err := json.Unmarshal(apiResp.Payload, &shows); err != nil {
			fmt.Println("followed, decode error :", err)
			return err
		}
		for _, title := range shows {
			b.Say(msg.Channel(), title.Title)
		}
	}
	return err
}

func (l *LetMeKnow) doLostPassword(b *ircbot.IrcBot, msg *ircbot.IrcMsg, token string) error {
	if len(msg.Trailing) < 3 {
		b.Say(msg.Channel(), "missing parameter")
		return nil
	}

	email := url.QueryEscape(strings.Join(msg.Trailing[2:], " "))
	apiURL := fmt.Sprintf("%s/%s/%s?email=%s", apiRoot, "users", "forgot_password", email)
	apiResp, err := l.get(apiURL, token, msg)
	if err != nil {
		return err
	}

	if apiResp.Status == "ok" {
		var resp string
		if err := json.Unmarshal(apiResp.Payload, &resp); err != nil {
			fmt.Println("error decode : ", err)
			return err
		}
		b.Say(msg.Channel(), resp)
	}
	return nil
}

//API helpers
func (l *LetMeKnow) get(apiURL, token string, msg *ircbot.IrcMsg) (*APIResp, error) {
	return callAPI(l.client, apiURL, nil, "GET", token, l.bot, msg)
}

func (l *LetMeKnow) post(apiURL string, body interface{}, token string, msg *ircbot.IrcMsg) (*APIResp, error) {
	w := &bytes.Buffer{}
	if err := json.NewEncoder(w).Encode(body); err != nil {
		return nil, err
	}
	return callAPI(l.client, apiURL, w, "POST", token, l.bot, msg)
}

func (l *LetMeKnow) delete(apiURL, token string, msg *ircbot.IrcMsg) (*APIResp, error) {
	return callAPI(l.client, apiURL, nil, "DELETE", token, l.bot, msg)
}

func callAPI(client *http.Client, apiURL string, body io.Reader, method, token string, b *ircbot.IrcBot, m *ircbot.IrcMsg) (*APIResp, error) {
	u, err := url.Parse(apiURL)
	if err != nil {
		fmt.Println("error parse url : ", err)
		return nil, err
	}
	if token != "" {
		v := u.Query()
		v.Add("token", token)
		u.RawQuery = v.Encode()
	}

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		fmt.Println("error create request : ", err)
		return nil, err
	}
	fmt.Printf("req : %+v\n", req)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error do request : ", err)
		return nil, err
	}
	defer resp.Body.Close()

	apiResp := &APIResp{}
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		fmt.Println("error decoding api response: ", err)
		return nil, err
	}

	if err := checkApiRespError(apiResp, b, m); err != nil {
		fmt.Println("api respond with error :", err)
		return nil, err
	}

	return apiResp, nil
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
