package actions

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Zaibon/ircbot"
)

func Mintpal(bot *ircbot.IrcBot) {
	threshold := -50.0

	for {
		time.Sleep(10 * time.Minute)

		resp, err := http.Get("https://api.mintpal.com/v1/market/summary")
		if err != nil {
			log.Fatalln(err)
		}
		defer resp.Body.Close()

		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		data := []map[string]string{}
		if err = json.Unmarshal(b, &data); err != nil {
			log.Fatalln(err)
		}

		for _, v := range data {
			change, err := strconv.ParseFloat(v["change"], 32)
			if err != nil {
				log.Println(v["code"] + " : " + err.Error())
			}

			if change <= threshold {
				output := fmt.Sprintf("gros mouvement sur le %s !! %f%% last price %s\n", v["code"], change, v["last_price"])
				bot.Say("#laminerie", output)
			}
		}
	}
}
