package actions

import (
	"encoding/json"
	"os"
	"testing"
)

func TestLink(t *testing.T) {
	f, err := os.Open("infoUrl.json")
	if err != nil {
		t.Error(err.Error())
	}

	data := make(map[string]string)
	dec := json.NewDecoder(f)
	dec.Decode(&data)

	if data["bter"] != "https://bter.com/" {
		t.Log(data)
		t.Fail()
	}
}
