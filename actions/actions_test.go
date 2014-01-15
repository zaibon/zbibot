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
	url := fmt.Sprintf(apiUrl, "getblocksfound", apiKey)
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
	t.Log(string(body))

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
}

func TestStatus(t *testing.T) {

	urlStatus := fmt.Sprintf(apiUrl, "getpoolstatus", apiKey)
	resp, err := http.Get(urlStatus)
	if err != nil {
		t.Error(err)
	}

	urlPublic := fmt.Sprintf(apiUrl, "public", apiKey)
	respPub, err := http.Get(urlPublic)
	if err != nil {
		t.Error(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}

	defer respPub.Body.Close()
	bodyPub, err := ioutil.ReadAll(respPub.Body)
	if err != nil {
		t.Error(err)
	}

	var data map[string]struct {
		Data poolStatus
	}
	if err := json.Unmarshal(body, &data); err != nil {
		t.Error(err)
	}

	var dataPub poolPublicInfo
	if err := json.Unmarshal(bodyPub, &dataPub); err != nil {
		t.Error(err)
	}

	ratio := (float64(dataPub.ShareCurRound) / data["getpoolstatus"].Data.EstShare) * 100
	hashRate := float32(data["getpoolstatus"].Data.HashRate) / 1000
	output := fmt.Sprintf("Pool Hashrate: %.3f khash/s | Pool Efficiency: %.2f%%%% | Current Difficulty: %f | Round %.3f%%%% | Workers: %d",
		hashRate, data["getpoolstatus"].Data.Efficency, data["getpoolstatus"].Data.NetDiff, ratio, dataPub.WorkersNbr)
	// Pool Hashrate: 18,374 khash | Pool Efficiency: 99.57% | Current difficulty: 55.299 | Round Estimate: 56626 | Current Round: 139991 | Round: 247.22% | Workers: 29

	t.Log(output)
}
