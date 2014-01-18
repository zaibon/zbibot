package actions

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestBlockUpdate(t *testing.T) {
	url := fmt.Sprintf(apiUrl+`&limit=1`, `getblocksfound`, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		t.Fail()
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fail()
	}
	resp.Body.Close()

	var data map[string]struct {
		Data []block
	}
	if err := json.Unmarshal(body, &data); err != nil {
		t.Fail()
	}
	blockInfo := data["getblocksfound"].Data[0]
	t.Log(blockInfo)
	output := fmt.Sprintf("BLOCK FOUND !!! #%d | %f %%%% | Mined By %s | Amount %f",
		blockInfo.Height, blockInfo.Ratio(), blockInfo.WorkerName, blockInfo.Amount)
	t.Log(output)
	t.Fail()
}
