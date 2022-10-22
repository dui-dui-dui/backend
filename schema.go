package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"
)

type TSMap struct {
	StartTS map[string]time.Time
	EndTS   map[string]time.Time
}

type Schema struct {
	TS          int64  `json:"ts"`
	Name        string `json:"name"`
	Description string `json:"description"`
	StartKey    string `json:"start_key"`
	EndKey      string `json:"end_key"`
}

func loadSchema() ([]Schema, TSMap, error) {
	type database struct {
		ID     int64 `json:"id"`
		DBName struct {
			O string `json:"O"`
		} `json:"db_name"`
	}

	type table struct {
		ID     int64 `json:"id"`
		DBName string
		Name   struct {
			O string `json:"O"`
		} `json:"name"`
	}

	tsMap := TSMap{
		StartTS: make(map[string]time.Time),
		EndTS:   make(map[string]time.Time),
	}

	var databases []database
	err := httpGet("http://"+*tidbAddr+"/schema", &databases)
	if err != nil {
		return nil, tsMap, err
	}

	var schemas []Schema
	schemas = append(schemas, Schema{
		Name:        "meta",
		Description: "meta data of tidb cluster",
		EndKey:      tableKey(0),
	})

	// system (mysql)
	var tables []table
	err = httpGet("http://"+*tidbAddr+"/schema/mysql", &tables)
	if err != nil {
		return nil, tsMap, err
	}
	var maxID int64
	for _, table := range tables {
		if table.ID > maxID && table.ID < 10000000 {
			maxID = table.ID
		}
	}
	schemas = append(schemas, Schema{
		Name:        "system",
		Description: "system tables of mysql database",
		StartKey:    tableKey(maxID),
		EndKey:      tableKey(maxID + 2),
	})

	// user tables
	var allTables []table
	for _, db := range databases {
		if strings.HasSuffix(db.DBName.O, "_SCHEMA") || db.DBName.O == "mysql" {
			continue
		}
		var tbls []table
		err = httpGet("http://"+*tidbAddr+"/schema/"+db.DBName.O, &tbls)
		if err != nil {
			return nil, tsMap, err
		}
		for i := range tbls {
			tbls[i].DBName = db.DBName.O
		}
		allTables = append(allTables, tbls...)
	}
	sort.Slice(allTables, func(i, j int) bool {
		return allTables[i].ID < allTables[j].ID
	})
	for _, t := range allTables {
		schemas = append(schemas, Schema{
			Name:        t.Name.O,
			Description: fmt.Sprintf("%s/%s", t.DBName, t.Name.O),
			StartKey:    tableKey(t.ID),
			EndKey:      tableKey(t.ID + 2),
		})
	}

	// future tables
	schemas = append(schemas, Schema{
		Name:        "default",
		Description: "future tables",
		StartKey:    schemas[len(schemas)-1].EndKey,
		EndKey:      "",
	})

	start, _ := time.Parse("2006-01-02 15:04:05 MST", "2022-01-01 00:00:00 CST")
	tsMap.StartTS[""] = start
	tsMap.EndTS[""] = start.AddDate(0, 0, len(schemas)-1)
	for i, s := range schemas {
		ts := start.AddDate(0, 0, i)
		tsMap.StartTS[s.StartKey], tsMap.EndTS[s.EndKey] = ts, ts
		schemas[i].TS = ts.UnixMicro()
	}

	return schemas, tsMap, nil
}

func httpGet(url string, v interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(v)
}
