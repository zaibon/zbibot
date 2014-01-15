package actions

import (
	"encoding/json"
	"fmt"
	"github.com/Zaibon/ircbot"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	apiUrl = "http://www.laminerie.eu/index.php?page=api&action=%s&api_key=%s"
	apiKey = "dfc87f06d1f4b93f7b97209396d48647ed0c53daf7ba33eaaa5a0f0fd152bbd0"
)

type block struct {
	AccountId     string  `json:account_id`
	Accounted     int     `json:accounted`
	Amount        float64 `jsonamount`
	BlockHash     string  `blockhash`
	Confirmations uint    `json:confirmations`
	Difficulty    float64 `json:difficulty`
	EstShares     float64 `json:estshares`
	Finder        string  `json:finder`
	Height        uint64  `json:height`
	Id            uint64  `json:id`
	IsAnonymous   int     `json:is_anonymous`
	ShareId       uint64  `json:share_id`
	Shares        float64 `json:shares`
	Time          int64   `json:time`
	WorkerName    string  `json:worker_name`
}

func (b *block) Ratio() float64 {
	return (b.Shares / b.EstShares) * 100.0
}

func LastBlock(b *ircbot.IrcBot, m *ircbot.IrcMsg) {
	url := fmt.Sprintf(apiUrl+"&limit=1", "getblocksfound", apiKey)

	resp, err := http.Get(url)
	if err != nil {
		b.Error <- err
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		b.Error <- err
		return
	}

	var data map[string]struct {
		Data []block
	}
	if err := json.Unmarshal(body, &data); err != nil {
		b.Error <- err
		return
	}
	lastBlock := data["getblocksfound"].Data[0]
	foundSince := time.Now().Sub(time.Unix(lastBlock.Time, 0))

	//%%%% => irc requierd double % too
	output := fmt.Sprintf("Last : #%d | Ratio %.3f%%%% | Confirmation %d | Mined by %s | Found Since %s",
		lastBlock.Height, lastBlock.Ratio(), lastBlock.Confirmations, lastBlock.Finder, foundSince)

	b.Say(m.Channel, output)
}
