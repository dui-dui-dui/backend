package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	var m map[string]any
	httpGet("http://127.0.0.1:8080/config", &m)
	ioutil.WriteFile("data/config_out.json", []byte(marshal(m)), 0644)
	ioutil.WriteFile("data/config_out.md", []byte(m["markdown"].(string)), 0644)
}

func TestLoadRegion(t *testing.T) {
	var m map[string]any
	httpGet("http://127.0.0.1:8080/region", &m)
	ioutil.WriteFile("data/region_out.md", []byte(m["markdown"].(string)), 0644)
}

func TestSaveConfig(t *testing.T) {
	var m map[string]any
	data, _ := ioutil.ReadFile("data/config_out.json")
	err := json.Unmarshal(data, &m)
	if err != nil {
		t.Fail()
	}
	md, _ := ioutil.ReadFile("data/config_out.md")
	m["markdown"] = string(md)
	data = []byte(marshal(m))
	res, err := http.Post("http://127.0.0.1:8080/save", "application/json", bytes.NewReader(data))
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	defer res.Body.Close()
	data, _ = ioutil.ReadAll(res.Body)
	fmt.Println(string(data))
}
