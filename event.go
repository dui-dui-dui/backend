package main

import (
	"fmt"
	"time"
)

type eSchema struct {
	TS          int64  `json:"ts"`
	Size        int    `json:"size"`
	Left        int    `json:"left"`
	Name        string `json:"name"`
	Description string `json:"description"`
	StartKey    string `json:"start_key"`
	EndKey      string `json:"end_key"`
}

type eEvent struct {
	EventString string `json:"eventString"`
	Ranges      struct {
		Event struct {
			From    int64  `json:"from"`
			To      int64  `json:"to"`
			Type    string `json:"type"`
			Content any    `json:"content"`
		} `json:"event"`
		Date struct {
			FromDateTime    time.Time `json:"fromDateTime"`
			ToDateTime      time.Time `json:"toDateTime"`
			OriginalString  string    `json:"originalString"`
			DateRangeInText struct {
				From    int64  `json:"from"`
				To      int64  `json:"to"`
				Type    string `json:"type"`
				Content any    `json:"content"`
			} `json:"dateRangeInText"`
		} `json:"date"`
	} `json:"ranges"`
	Event struct {
		EventDescription string   `json:"eventDescription"`
		Tags             []string `json:"tags"`
		Supplemental     []string `json:"supplemental"`
		GooglePhotosLink string   `json:"googlePhotosLink"`
		Locations        []string `json:"locations"`
		ID               string   `json:"id"`
		Percent          int      `json:"percent"`
	} `json:"event"`
}

type eGroup struct {
	Tags  []string `json:"tags"`
	Title string   `json:"title"`
	Range struct {
		Min    time.Time `json:"min"`
		Max    time.Time `json:"max"`
		Latest time.Time `json:"latest"`
	} `json:"range"`
	StartExpanded bool     `json:"startExpanded"`
	Style         string   `json:"style"`
	Events        []eEvent `json:"events"`
}

func convertToEvents(schemas []Schema, groups []Group) ([]eSchema, []eGroup) {
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
	var egs []eGroup

	for _, g := range groups {
		var eg eGroup
		eg.StartExpanded = true
		eg.Style = "group"
		eg.Title = g.ID
		eg.Range.Min = startTS[""]
		eg.Range.Max = endTS[""]
		eg.Range.Latest = endTS[""]

		for i, r := range g.Rules {
			var ee eEvent
			ee.EventString = r.ID
			ee.Ranges.Event.Type = "event"
			ee.Ranges.Date.FromDateTime = startTS[r.StartKeyHex]
			ee.Ranges.Date.ToDateTime = endTS[r.EndKeyHex]
			ee.Ranges.Date.DateRangeInText.Type = "dateRange"
			ee.Event.EventDescription = display(r.Count, r.Role)
			eg.Events = append(eg.Events, ee)

			if i == 0 || startTS[r.StartKeyHex].Before(eg.Range.Min) {
				eg.Range.Min = startTS[r.StartKeyHex]
			}
			if i == 0 || endTS[r.EndKeyHex].After(eg.Range.Max) {
				eg.Range.Max = endTS[r.EndKeyHex]
				eg.Range.Latest = endTS[r.EndKeyHex]
			}
		}

		egs = append(egs, eg)
	}
	return ess, egs
}

func display(count int, name string) string {
	if count > 1 {
		return fmt.Sprintf("%d %ss", count, name)
	}
	return fmt.Sprintf("%d %s", count, name)
}
