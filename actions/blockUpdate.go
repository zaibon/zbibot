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
	for coin, api := range apis {

		var previous_height uint64 = 0
		url := fmt.Sprintf(api.Url+`&limit=1`, `getblocksfound`, api.Key)

		go func(b *ircbot.IrcBot, coin string, api apiInfo) {
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

						output := fmt.Sprintf("[%s] BLOCK FOUND !!! #%d | Ratio %f %%%% | Mined By %s | Amount %f",
							coin, blockInfo.Height, blockInfo.Ratio(), blockInfo.WorkerName, blockInfo.Amount)

						b.Say(b.Channel[0], output) //FIXME : we can't be sure the first channel is the good one
					} else {
						fmt.Printf("INFO : [%s],no new blocks\n", coin)
					}
				}
				//check every 30 seconds for new found block
				time.Sleep(30 * time.Second)
			}
		}(b, coin, api)
	}
}
