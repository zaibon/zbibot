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

	//launch a go routine for each pool we watch
	for _, api := range apis {

		var previous_height uint64 = 0
		url := fmt.Sprintf(api.Url+`&limit=1`, `getblocksfound`, api.Key)

		go func(b *ircbot.IrcBot) {
			for {

				resp, err := http.Get(url)
				if err != nil {
					fmt.Println("ERROR : ", err.Error())
					continue
				}

				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					fmt.Println("ERROR : ", err.Error())
					continue
				}
				resp.Body.Close()

				var data map[string]struct {
					Data []block
				}
				if err := json.Unmarshal(body, &data); err != nil {
					fmt.Println("ERROR : ", err.Error())
					continue
				}

				if len(data[`getblocksfound`].Data) > 0 {

					blockInfo := data["getblocksfound"].Data[0]

					//first round
					if previous_height == 0 {
						previous_height = blockInfo.Height
					} else if previous_height != blockInfo.Height {
						//new block found !
						previous_height = blockInfo.Height
						//FIXME : we can't be sure the first channel is the good one
						b.Say(b.Channel[0], "Hey ! devinez quoi ??")
						time.Sleep(500 * time.Millisecond)

						output := fmt.Sprintf("BLOCK FOUND !!! #%d | Ratio %f %%%% | Mined By %s | Amount %f",
							blockInfo.Height, blockInfo.Ratio(), blockInfo.WorkerName, blockInfo.Amount)
						b.Say(b.Channel[0], output)
					} else {
						fmt.Println("INFO : no new blocks")
					}
				}
				//check every 30 seconds for new found block
				time.Sleep(30 * time.Second)
			}
		}(b)
	}
}
