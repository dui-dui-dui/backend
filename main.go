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
		log.Fatal("load schema", err)
	}
	groups, err := loadRuleConfig()
	if err != nil {
		log.Fatal("load rule config", err)
	}
	labels, err := loadStoreLabels()
	if err != nil {
		log.Fatal("load store labels", err)
	}

	es, eg := convertToEvents(schemas, groups)

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(map[string]any{
		"schemas": es,
		"groups":  eg,
		"labels":  labels,
	})
}
