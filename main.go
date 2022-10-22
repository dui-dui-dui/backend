package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
)

var (
	tidbAddr = flag.String("tidb", "127.0.0.1:10080", "tidb address")
	pdAddr   = flag.String("pd", "127.0.0.1:2379", "pd address:port")
	addr     = flag.String("addr", ":8080", "server address:port")
	dev      = flag.Bool("dev", false, "development mode")
)

var SCHEMAS []Schema
var STORES []Store
var STORELABELS []StoreLabel
var TSMAP TSMap

var (
	//go:embed testdata/schemas.json
	schemasJSON []byte
	//go:embed testdata/stores.json
	storesJSON []byte
	//go:embed testdata/storelabels.json
	storelabelsJSON []byte
	//go:embed testdata/tsmap.json
	tsmapJSON []byte
	//go:embed testdata/groups.json
	groupsJSON []byte
)

func main() {
	flag.Parse()

	if *dev {
		json.Unmarshal(schemasJSON, &SCHEMAS)
		json.Unmarshal(storesJSON, &STORES)
		json.Unmarshal(storelabelsJSON, &STORELABELS)
		json.Unmarshal(tsmapJSON, &TSMAP)
	} else {
		var err error
		SCHEMAS, TSMAP, err = loadSchema()
		if err != nil {
			log.Fatal("load schema", err)
		}
		STORES, STORELABELS, err = loadStores()
		if err != nil {
			log.Fatal("load stores", err)
		}
	}

	http.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		log.Println("handle /config")
		data, err := handleConfig()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(data))
	})
	http.HandleFunc("/region", func(w http.ResponseWriter, r *http.Request) {
		log.Println("handle /region")
		data, err := handleRegions()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(data))
	})
	http.HandleFunc("/save", func(w http.ResponseWriter, r *http.Request) {
		log.Println("handle /save")
		data, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		type saveData struct {
			Markdown string  `json:"markdown"`
			Groups   []Group `json:"groups"`
		}
		var sd saveData
		err = json.Unmarshal(data, &sd)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = handleSave(sd.Markdown, sd.Groups)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte("ok"))
	})

	log.Println("start server", *addr, "dev mode:", *dev)
	http.ListenAndServe(*addr, nil)
}

func handleConfig() (string, error) {
	groups, err := loadRuleConfig()
	if err != nil {
		log.Fatal("load rule config", err)
	}

	return marshal(map[string]any{
		"markdown": convertGroupMarkdown(groups),
		"schemas":  SCHEMAS,
		"groups":   groups,
		"labels":   STORELABELS,
	}), nil
}

func handleRegions() (string, error) {
	if *dev {
		return marshal(map[string]any{
			"markdown": mockRegionMarkdown(),
		}), nil
	}

	regions, err := loadRegions()
	if err != nil {
		return "", err
	}
	sort.Slice(regions, func(i, j int) bool {
		return regions[i].StartKey < regions[j].StartKey
	})
	return marshal(map[string]any{
		"markdown": convertRegionMarkdown(regions),
	}), nil
}

func handleSave(markdown string, groups []Group) error {
	newGroups, err := parseMarkdown(markdown)
	if err != nil {
		return err
	}
	for i := range newGroups {
		for _, g0 := range groups {
			if newGroups[i].ID == g0.ID {
				newGroups[i].Override = g0.Override
			}
			for j := range newGroups[i].Rules {
				for _, r0 := range g0.Rules {
					if newGroups[i].Rules[j].ID == r0.ID {
						newGroups[i].Rules[j].Override = r0.Override
						newGroups[i].Rules[j].Role = r0.Role
						newGroups[i].Rules[j].Count = r0.Count
						newGroups[i].Rules[j].LabelConstraints = r0.LabelConstraints
						newGroups[i].Rules[j].LocationLabels = r0.LocationLabels
						newGroups[i].Rules[j].IsolationLevel = r0.IsolationLevel
					}
				}
			}
		}
	}
	if *dev {
		log.Println("groups.json is updated")
		ioutil.WriteFile("groups.json", []byte(marshal(newGroups)), 0644)
		return nil
	}
	// TODO: save to pd
	return nil
}

func marshal(v interface{}) string {
	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	enc.SetIndent("", "  ")
	enc.Encode(v)
	return b.String()
}
