package actions

import (
	"encoding/json"
	"fmt"
	"github.com/Zaibon/ircbot"
	"io/ioutil"
	"net/http"
	"time"
)

func FindBlock(b *ircbot.IrcBot) {
	var previous_height uint64 = 0
	url := fmt.Sprintf(apiUrl+`&limit=1`, `getblocksfound`, apiKey)

	go func(b *ircbot.IrcBot) {
		for {

			resp, err := http.Get(url)
			if err != nil {
				b.Error <- err
				return
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				b.Error <- err
				return
			}
			resp.Body.Close()

			var data map[string]struct {
				Data []block
			}
			if err := json.Unmarshal(body, &data); err != nil {
				b.Error <- err
				return
			}

			blockInfo := data["getblocksfound"].Data[0]
			if previous_height == 0 {
				previous_height = blockInfo.Height
			} else if previous_height != blockInfo.Height {
				//new block found !

				//we can't be sure the first channel is the good one. need to be fix
				b.Say(b.Channel[0], "Hey ! devinez quoi ??")
				time.Sleep(500 * time.Millisecond)

				output := fmt.Sprintf("BLOCK FOUND !!! #%d | %d shares | Mined By %s | Amount %f",
					blockInfo.Height, blockInfo.Shares, blockInfo.Amount)
				b.Say(b.Channel[0], output)
			} else {
				fmt.Println("INFO : no new blocks")
			}

			//check every 30 seconds for new found block
			time.Sleep(30 * time.Second)
		}
	}(b)
}
