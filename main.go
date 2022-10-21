package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
)

var (
	tidbAddr = flag.String("tidb", "127.0.0.1:10080", "tidb address")
	pdAddr   = flag.String("pd", "127.0.0.1:2379", "pd address:port")
	addr     = flag.String("addr", ":8080", "server address:port")
)

func main() {
	flag.Parse()

	schemas, err := loadSchema()
	if err != nil {
		log.Fatal(err)
	}
	groups, err := loadRuleConfig()
	if err != nil {
		log.Fatal(err)
	}
	labels, err := loadStoreLabels()
	if err != nil {
		log.Fatal(err)
	}

	type config struct {
		Schemas     []Schema     `json:"schemas"`
		RuleConfig  []Group      `json:"rule_config"`
		StoreLables []StoreLabel `json:"store_labels"`
	}

	for _, g := range groups {
		for i := range g.Rules {
			for _, s := range schemas {
				if g.Rules[i].StartKeyHex == s.StartKey {
					g.Rules[i].StartSchema = s.Name
				}
				if g.Rules[i].EndKeyHex == s.EndKey {
					g.Rules[i].EndSchema = s.Name
				}
			}
		}
	}

	cfg := config{
		Schemas:     schemas,
		RuleConfig:  groups,
		StoreLables: labels,
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(cfg)
}
