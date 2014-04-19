package actions

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestOverAllStats(t *testing.T) {
	urlStatus := fmt.Sprintf(apiUrl, "getblockstats", apiKey)
	resp, err := http.Get(urlStatus)
	if err != nil {
		t.Fail()
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fail()
	}

	var data map[string]struct {
		Data blockStats
	}
	if err := json.Unmarshal(body, &data); err != nil {
		t.Fail()
	}
	stats := data["getblockstats"].Data

	output1 := fmt.Sprintf("Last Hour  | Found : %4s | Valid : %4s | Orphan %s | Efficiency %f %%%%",
		stats.OneHourhTotal, stats.OneHourhValid, stats.OneHourhOrphan, stats.OneHourEfficency())

	output2 := fmt.Sprintf("Last 24H   | Found : %4s | Valid : %4s | Orphan %s | Efficiency %f %%%%",
		stats.YesterdayTotal, stats.YesterdayValid, stats.YesterdayOrphan, stats.YestardayEfficency())

	output3 := fmt.Sprintf("Last Week  | Found : %4s | Valid : %4s | Orphan %s | Efficiency %f %%%%",
		stats.WeekTotal, stats.WeekValid, stats.WeekOrphan, stats.WeekEfficency())

	output4 := fmt.Sprintf("Last Year  | Found : %4s | Valid : %4s | Orphan %s | Efficiency %f %%%%",
		stats.YearTotal, stats.YearValid, stats.YearOrphan, stats.YearEfficency())

	output5 := fmt.Sprintf("Last Year  | Found : %4d | Valid : %4s | Orphan %s | Efficiency %f %%%%",
		stats.Total, stats.TotalValid, stats.TotalOrphan, stats.TotalEfficency())

	t.Log(output1)
	t.Log(output2)
	t.Log(output3)
	t.Log(output4)
	t.Log(output5)
}
