package gocsv

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
)

func Test_Gocsv(t *testing.T) {

	resp, err := http.Get("https://api.github.com/search/code?q=addClass+in:file+language:js+repo:jquery/jquery")
	if err != nil {
		t.Error(err.Error())
	}
	defer resp.Body.Close()
	request, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err.Error())
	}

	var data struct {
		TotalCount       int           `json:"total_count"`
		IncompleteResult bool          `json:"incomplete_results"`
		Items            []interface{} `json:"items"`
	}
	err = json.Unmarshal(request, &data)
	if err != nil {
		t.Error(err.Error())
	}

	err = NewCsvFile("test.csv", "", []string{}, data.Items)
	if err != nil {
		t.Error(err.Error())
	}
	t.Log("success")
}
