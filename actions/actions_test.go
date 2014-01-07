package actions

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"testing"
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

	if jsonData["low"] != 0.00085961 {
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
