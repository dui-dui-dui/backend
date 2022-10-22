package main

import (
	"fmt"
	"strings"
	"time"
)

func convertToMarkdown(schemas []Schema, groups []Group) ([]eSchema, string) {
	start, _ := time.Parse("2006-01-02 15:04:05 MST", "2022-01-01 00:00:00 CST")
	var ess []eSchema
	startTS, endTS := make(map[string]time.Time), make(map[string]time.Time)
	startTS[""] = start
	endTS[""] = start.AddDate(0, 0, len(schemas)-1)
	for i, s := range schemas {
		ts := start.AddDate(0, 0, i)
		startTS[s.StartKey], endTS[s.EndKey] = ts, ts
		ess = append(ess, eSchema{
			TS:          ts.UnixMicro(),
			Name:        s.Name,
			Description: s.Description,
			StartKey:    s.StartKey,
			EndKey:      s.EndKey,
		})
	}

	var sb strings.Builder
	for _, g := range groups {
		sb.WriteString("group " + g.ID + "\n\n")
		for _, r := range g.Rules {
			sb.WriteString(fmt.Sprintf("%s - %s: %s\n", startTS[r.StartKeyHex].Format("01/02/2006"), endTS[r.EndKeyHex].Format("01/02/2006"), r.ID))
		}
	}
	return ess, sb.String()
}
