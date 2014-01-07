package actions

import (
	"encoding/json"
	"fmt"
	"github.com/Zaibon/ircbot"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type cryptocoinchartsJson struct {
	Id             string `json:id`
	Price          string `json:price`
	PriceBefore24h string `json:price_before_24h`
	Volumefirst    string `json:volume_first`
	VolumeSecond   string `json:volume_second`
	VolumeBtc      string `json:volume_btc`
	BestMarket     string `json:best_market`
	LatestTrade    string `json:latest_trade`
}

func ExchRate(b *ircbot.IrcBot, m *ircbot.IrcMsg) {
	resp, err := http.PostForm("http://www.cryptocoincharts.info/v2/api/tradingPairs", url.Values{"pairs": {"mec_btc,btc_usb,btc_eur"}})
	if err != nil {
		b.Say(m.Channel, err.Error())
		b.Error <- err
		return
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		b.Say(m.Channel, err.Error())
		b.Error <- err
		return
	}

	var data []cryptocoinchartsJson
	if err := json.Unmarshal(bytes, &data); err != nil {
		b.Say(m.Channel, err.Error())
		b.Error <- err
		return
	}

	var btcEur float64
	var mecBtc float64

	for _, v := range data {
		b.Say(m.Channel, fmt.Sprintf("change       : %s", v.Id))
		b.Say(m.Channel, fmt.Sprintf("price        : %s", v.Price))
		// b.Say(m.Channel, fmt.Sprintf("price -24h   : %s", v.PriceBefore24h))
		b.Say(m.Channel, "------------------")

		switch v.Id {
		case "btc/eur":
			btcEur, _ = strconv.ParseFloat(v.Price, 10)
		case "mec/btc":
			mecBtc, _ = strconv.ParseFloat(v.Price, 10)
		}
	}

	b.Say(m.Channel, fmt.Sprintf("change       : %s", "mec/eur"))
	b.Say(m.Channel, fmt.Sprintf("price        : %f", btcEur*mecBtc))
	b.Say(m.Channel, "------------------")
}
