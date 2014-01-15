package actions

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"
)

func TestLink(t *testing.T) {
	f, err := os.Open("../infoUrl.json")
	if err != nil {
		t.Error(err.Error())
	}

	data := make(map[string][]string)
	dec := json.NewDecoder(f)
	dec.Decode(&data)

	if data["bter"][0] != "https://bter.com/" {
		t.Log(data)
		t.Fail()
	}
}

func TestTicker(t *testing.T) {
	resp, err := http.Get("https://bter.com/api/1/ticker/mec_btc")
	if err != nil {
		t.Error(err)
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}

	jsonData := make(map[string]interface{})
	err = json.Unmarshal(bytes, &jsonData)
	if err != nil {
		t.Error(err)
	}

	// r := jsonData["result"]
	// result, _ := r.(string)
	// if result != "true" {
	if jsonData["result"] != "true" {
		t.Fail()
	}
}

func TestExchRate(t *testing.T) {
	resp, err := http.PostForm("http://www.cryptocoincharts.info/v2/api/tradingPairs", url.Values{"pairs": {"mec_btc,btc_usb,btc_eur"}})
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}

	var data []cryptocoinchartsJson
	if err := json.Unmarshal(bytes, &data); err != nil {
		t.Error(err)
	}

	if data[0].Id != "btc/eur" {
		t.Log(data)
		t.Fail()
	}
}

func TestLastBlock(t *testing.T) {
	const (
		// baseUrl = "http://laminerie.eu?index.php"
		apiUrl = "http://www.laminerie.eu/index.php?page=api&action=%s&api_key=%s"
		apiKey = "dfc87f06d1f4b93f7b97209396d48647ed0c53daf7ba33eaaa5a0f0fd152bbd0"
	)
	url := fmt.Sprintf(apiUrl+"&limit=1", "getblocksfound", apiKey)
	t.Log(url)
	resp, err := http.Get(url)
	if err != nil {
		t.Error(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}

	var data map[string]struct {
		Data []block
	}
	if err := json.Unmarshal(body, &data); err != nil {
		t.Error(err)
	}

	lastBlock := data["getblocksfound"].Data[0]
	foundSince := time.Now().Sub(time.Unix(lastBlock.Time, 0))

	s := fmt.Sprintf("Last : #%d | Ratio %.3f%% | Confirmation %d | Mined by %s | Found Since %s",
		lastBlock.Height, lastBlock.Ratio(), lastBlock.Confirmations, lastBlock.Finder, foundSince)
	t.Log(s)
	t.Fail()
}
