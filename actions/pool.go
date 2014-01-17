package actions

//these actions are heavily inspired by https://github.com/WKNiGHT-/mpos-bot

import (
	`encoding/json`
	`fmt`
	`github.com/Zaibon/ircbot`
	`io/ioutil`
	`net/http`
	"strconv"
	`time`
)

const (
	apiUrl = `http://www.laminerie.eu/index.php?page=api&action=%s&api_key=%s`
	apiKey = `dfc87f06d1f4b93f7b97209396d48647ed0c53daf7ba33eaaa5a0f0fd152bbd0`
)

type block struct {
	AccountId     uint    `json:"account_id"`
	Accounted     int     `json:"accounted"`
	Amount        float64 `json:"jsonamount"`
	BlockHash     string  `json:"blockhash"`
	Confirmations uint    `json:"confirmations"`
	Difficulty    float64 `json:"difficulty"`
	EstShares     float64 `json:"estshares"`
	Finder        string  `json:"finder"`
	Height        uint64  `json:"height"`
	Id            uint64  `json:"id"`
	IsAnonymous   int     `json:"is_anonymous"`
	ShareId       uint64  `json:"share_id"`
	Shares        float64 `json:"shares"`
	Time          int64   `json:"time"`
	WorkerName    string  `json:"worker_name"`
}

func (b *block) Ratio() float64 {
	return (b.Shares / b.EstShares) * 100.0
}

type poolStatus struct {
	CurrentNetWorkBlock uint32  `json:"currentnetworkblock"`
	Efficency           float32 `json:"efficiency"`
	EstShare            float64 `json:"estshares"`
	EstTime             float64 `json:"esttime"`
	HashRate            uint32  `json:"hashrate"`
	LastBlock           uint    `json:"lastblock"`
	NetHashRate         uint32  `json:"nethashrate"`
	NetDiff             float32 `json:"networkdiff"`
	NextNetWorkBlock    uint32  `json:"nextnetworkblock"`
	PoolName            string  `json:"pool_name"`
	TimeSinceLast       uint    `json:"timesincelast"`
	WorkersNbr          uint    `json:"workers"`
}

type poolPublicInfo struct {
	HashRate      int    `json:"hashrate"`
	LastBlock     int    `json:"last_block"`
	NetHashRate   int    `json:"network_hashrate"`
	PoolName      string `json:"pool_name"`
	ShareCurRound int    `json:"shares_this_round"`
	WorkersNbr    int    `json:"workers"`
}

type user struct {
	HashRate  int    `json:"hashrate"`
	ShareRate string `sjson:"shareharerate"`
	UserName  string `json:"username"`
	Share     shares `json:"shares"`
}

type shares struct {
	DonatePercent float64 `json:"donate_percent"`
	Id            uint64  `json:"id"`
	Invalid       uint    `json:"invalid"`
	IsAnonymous   int     `json:"is_anonymous"`
	Username      string  `json:"username"`
	Valid         uint64  `json:"valid"`
}

func LastBlock(b *ircbot.IrcBot, m *ircbot.IrcMsg) {
	var (
		nbrLastBlock int = 1
		errConv      error
	)
	if len(m.Args) >= 2 {
		nbrLastBlock, errConv = strconv.Atoi(m.Args[1])
		if errConv != nil {
			nbrLastBlock = 1
		}
	}

	url := fmt.Sprintf(apiUrl, `getblocksfound`, apiKey)

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

	now := time.Now()
	for i := 0; i < nbrLastBlock; i++ {
		lastBlock := data[`getblocksfound`].Data[i]
		foundSince := now.Sub(time.Unix(lastBlock.Time, 0))

		//%%%% => irc requierd double % too
		output := fmt.Sprintf(`Last : #%d | Ratio %.3f%%%% | Confirmation %3d | Mined by %10s| Found Since %s`,
			lastBlock.Height, lastBlock.Ratio(), lastBlock.Confirmations, lastBlock.Finder, foundSince)

		b.Say(m.Channel, output)
	}
}

func Status(b *ircbot.IrcBot, m *ircbot.IrcMsg) {
	urlStatus := fmt.Sprintf(apiUrl, `getpoolstatus`, apiKey)
	resp, err := http.Get(urlStatus)
	if err != nil {
		b.Error <- err
		return
	}

	urlPublic := fmt.Sprintf(apiUrl, `public`, apiKey)
	respPub, err := http.Get(urlPublic)
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

	defer respPub.Body.Close()
	bodyPub, err := ioutil.ReadAll(respPub.Body)
	if err != nil {
		b.Error <- err
		return
	}

	var data map[string]struct {
		Data poolStatus
	}
	if err := json.Unmarshal(body, &data); err != nil {
		b.Error <- err
		return
	}

	var dataPub poolPublicInfo
	if err := json.Unmarshal(bodyPub, &dataPub); err != nil {
		b.Error <- err
		return
	}

	ratio := (float64(dataPub.ShareCurRound) / data[`getpoolstatus`].Data.EstShare) * 100
	poolHashRate := float32(data[`getpoolstatus`].Data.HashRate) / 1000
	netHashRate := float32(data[`getpoolstatus`].Data.NetHashRate) / 1000000
	output := fmt.Sprintf(`Pool Hashrate: %.3f KHash/s | Net Hashrate : %f MHash/s | Pool Efficiency: %.2f%%%% | Current Difficulty: %f | Round %.3f%%%% | Workers: %d`,
		poolHashRate, netHashRate, data[`getpoolstatus`].Data.Efficency, data[`getpoolstatus`].Data.NetDiff, ratio, dataPub.WorkersNbr)

	b.Say(m.Channel, output)
}

func User(b *ircbot.IrcBot, m *ircbot.IrcMsg) {
	var username string
	if len(m.Args) < 2 {
		username = m.Nick
	} else {
		username = m.Args[1]
	}

	urlStatus := fmt.Sprintf(apiUrl+"&id=%s", "getuserstatus", apiKey, username)
	resp, err := http.Get(urlStatus)
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
		Data user
	}
	if err := json.Unmarshal(body, &data); err != nil {
		b.Error <- err
		return
	}
	user := data["getuserstatus"].Data
	output := fmt.Sprintf("Username: %s | HashRate: %d Kh/s | ShareRate %s | Share Valid : %d | Share Invalid: %d",
		user.UserName, user.HashRate, user.ShareRate, user.Share.Valid, user.Share.Invalid)

	b.Say(m.Channel, output)
}
