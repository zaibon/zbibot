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

type blockStats struct {
	OneHourhAmout     float64 `json:"1HourAmount"`
	OneHourhDiff      float64 `json:"1HourDifficulty"`
	OneHourhEstShares uint    `json:"1HourEstimatedShares"`
	OneHourhOrphan    string  `json:"1HourOrphan"`
	OneHourhShares    string  `json:"1HourShares"`
	OneHourhTotal     string  `json:"1HourTotal"`
	OneHourhValid     string  `json:"1HourValid"`

	YesterdayAmout     float64 `json:"24HourAmount"`
	YesterdayDiff      float64 `json:"24HourDifficulty"`
	YesterdayEstShares uint    `json:"24HourEstimatedShares"`
	YesterdayOrphan    string  `json:"24HourOrphan"`
	YesterdayShares    string  `json:"24HourShares"`
	YesterdayTotal     string  `json:"24HourTotal"`
	YesterdayValid     string  `json:"24HourValid"`

	WeekAmout     float64 `json:"7DaysAmount"`
	WeekDiff      float64 `json:"7DaysDifficulty"`
	WeekEstShares uint    `json:"7DaysEstimatedShares"`
	WeekOrphan    string  `json:"7DaysOrphan"`
	WeekShares    string  `json:"7DaysShares"`
	WeekTotal     string  `json:"7DaysTotal"`
	WeekValid     string  `json:"7DaysValid"`

	FourWeeksAmout      float64 `json:"4WeeksAmount"`
	FourWeekshDiff      float64 `json:"4WeeksDifficulty"`
	FourWeekshEstShares uint    `json:"4WeeksEstimatedShares"`
	FourWeekshOrphan    string  `json:"4WeeksOrphan"`
	FourWeekshShares    string  `json:"4WeeksShares"`
	FourWeekshTotal     string  `json:"4WeeksTotal"`
	FourWeekshValid     string  `json:"4WeeksValid"`

	YearAmout     float64 `json:"12MonthAmount"`
	YearDiff      float64 `json:"12MonthDifficulty"`
	YearEstShares uint    `json:"12MonthEstimatedShares"`
	YearOrphan    string  `json:"12MonthOrphan"`
	YearShares    string  `json:"12MonthShares"`
	YearTotal     string  `json:"12MonthTotal"`
	YearValid     string  `json:"12MonthValid"`

	Total          uint    `json:"Total"`
	TotalAmount    float64 `json:"TotalAmount"`
	TotalDiff      float64 `json:"TotalDifficulty"`
	TotalEstSahres uint    `json:"TotalEstimatedShares"`
	TotalOrphan    string  `json:"TotalOrphan"`
	TotalShares    string  `json:"TotalShares"`
	TotalValid     string  `json:"TotalValid"`
}

func (b *blockStats) OneHourEfficency() float32 {
	var oneHourEfficency float32
	if b.OneHourhShares == "0" {
		return 0
	} else {
		oneHourShares, err := strconv.ParseFloat(b.OneHourhShares, 32)
		if err != nil {
			return 0
		}
		oneHourEfficency = (float32(oneHourShares) / float32(b.OneHourhEstShares)) * 100
		return oneHourEfficency
	}
}

func (b *blockStats) YestardayEfficency() float32 {
	var yesterdayEfficency float32
	if b.YesterdayShares == "0" {
		return 0
	} else {
		yesterdayShares, err := strconv.ParseFloat(b.YesterdayShares, 32)
		if err != nil {
			return 0
		}
		yesterdayEfficency = (float32(yesterdayShares) / float32(b.YesterdayEstShares)) * 100
		return yesterdayEfficency
	}
}

func (b *blockStats) WeekEfficency() float32 {
	var WeekEfficency float32
	if b.WeekShares == "0" {
		return 0
	} else {
		WeekShares, err := strconv.ParseFloat(b.WeekShares, 32)
		if err != nil {
			return 0
		}
		WeekEfficency = (float32(WeekShares) / float32(b.WeekEstShares)) * 100
		return WeekEfficency
	}
}

func (b *blockStats) FourWeekEfficency() float32 {
	var fourWeekEfficency float32
	if b.FourWeekshShares == "0" {
		return 0
	} else {
		fourWeekShares, err := strconv.ParseFloat(b.FourWeekshShares, 32)
		if err != nil {
			return 0
		}
		fourWeekEfficency = (float32(fourWeekShares) / float32(b.FourWeekshEstShares)) * 100
		return fourWeekEfficency
	}
}

func (b *blockStats) YearEfficency() float32 {
	var yearEfficency float32
	if b.YearShares == "0" {
		return 0
	} else {
		yearShares, err := strconv.ParseFloat(b.YearShares, 32)
		if err != nil {
			return 0
		}
		yearEfficency = (float32(yearShares) / float32(b.YearEstShares)) * 100
		return yearEfficency
	}
}

func (b *blockStats) TotalEfficency() float32 {
	var totalEfficency float32
	if b.TotalShares == "0" {
		return 0
	} else {
		totalShares, err := strconv.ParseFloat(b.TotalShares, 32)
		if err != nil {
			return 0
		}
		totalEfficency = (float32(totalShares) / float32(b.TotalEstSahres)) * 100
		return totalEfficency
	}
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
	output := fmt.Sprintf(`Pool Hashrate: %.3f MHash/s | Net Hashrate : %f MHash/s | Pool Efficiency: %.2f%%%% | Current Difficulty: %f | Round %.3f%%%% | Workers: %d`,
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

func OverallStats(b *ircbot.IrcBot, m *ircbot.IrcMsg) {
	urlStatus := fmt.Sprintf(apiUrl, "getblockstats", apiKey)
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
		Data blockStats
	}
	if err := json.Unmarshal(body, &data); err != nil {
		b.Error <- err
		return
	}
	stats := data["getblockstats"].Data

	output1 := fmt.Sprintf("Last Hour  | Found : %4s | Valid : %4s | Orphan %s | Efficiency %f %%%%",
		stats.OneHourhTotal, stats.OneHourhValid, stats.OneHourhOrphan, stats.OneHourEfficency())

	output2 := fmt.Sprintf("Last 24H   | Found : %4s | Valid : %4s | Orphan %s | Efficiency %f %%%%",
		stats.YesterdayTotal, stats.YesterdayValid, stats.YesterdayOrphan, stats.YestardayEfficency())

	output3 := fmt.Sprintf("Last Week  | Found : %4s | Valid : %4s | Orphan %s | Efficiency %f %%%%",
		stats.WeekTotal, stats.WeekValid, stats.WeekOrphan, stats.WeekEfficency())

	output4 := fmt.Sprintf("Last Year  | Found : %4s | Valid : %4s | Orphan %s | Efficiency %f %%%%",
		stats.YearTotal, stats.YearValid, stats.YearOrphan, stats.YearEfficency())

	output5 := fmt.Sprintf("Last Year  | Found : %4d | Valid : %4s | Orphan %s | Efficiency %f %%%%",
		stats.Total, stats.TotalValid, stats.TotalOrphan, stats.TotalEfficency())

	b.Say(m.Channel, output1)
	b.Say(m.Channel, output2)
	b.Say(m.Channel, output3)
	b.Say(m.Channel, output4)
	b.Say(m.Channel, output5)
}
