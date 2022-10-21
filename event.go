package main

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
			FromDateTime    int64  `json:"fromDateTime"`
			ToDateTime      int64  `json:"toDateTime"`
			OriginalString  string `json:"originalString"`
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
		Min    int64 `json:"min"`
		Max    int64 `json:"max"`
		Latest int64 `json:"latest"`
	} `json:"range"`
	StartExpanded bool   `json:"startExpanded"`
	Style         string `json:"style"`
}

func convertSchemas(ss []Schema) []eSchema {
	start := int64(1666195200000)
	var es []eSchema
	for i, s := range ss {
		es = append(es, eSchema{
			TS:          start + int64(i)*24*3600*1000,
			Name:        s.Name,
			Description: s.Description,
			StartKey:    s.StartKey,
			EndKey:      s.EndKey,
		})
	}
	return es
}

func convertToEvents(schemas []Schema, groups []Group) ([]eSchema, []eEvent, []eGroup) {
	start := int64(1666195200000)
	var ess []eSchema
	startTS, endTS := make(map[string]int64), make(map[string]int64)
	startTS[""] = start
	endTS[""] = start + int64(len(schemas))*24*3600*1000
	for i, s := range schemas {
		ts := start + int64(i)*24*3600*1000
		startTS[s.StartKey], endTS[s.EndKey] = ts, ts
		ess = append(ess, eSchema{
			TS:          ts,
			Name:        s.Name,
			Description: s.Description,
			StartKey:    s.StartKey,
			EndKey:      s.EndKey,
		})
	}
	var ees []eEvent
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
			ee.Ranges.Event.From = startTS[r.StartKeyHex]
			ee.Ranges.Event.To = endTS[r.EndKeyHex]
			ee.Ranges.Event.Type = "event"
			ee.Ranges.Date.FromDateTime = startTS[r.StartKeyHex]
			ee.Ranges.Date.ToDateTime = endTS[r.EndKeyHex]
			ee.Ranges.Date.DateRangeInText.From = startTS[r.StartKeyHex]
			ee.Ranges.Date.DateRangeInText.To = endTS[r.EndKeyHex]
			ee.Ranges.Date.DateRangeInText.Type = "dateRange"
			ees = append(ees, ee)

			if i == 0 || startTS[r.StartKeyHex] < eg.Range.Min {
				eg.Range.Min = startTS[r.StartKeyHex]
			}
			if i == 0 || endTS[r.EndKeyHex] > eg.Range.Max {
				eg.Range.Max = endTS[r.EndKeyHex]
				eg.Range.Latest = endTS[r.EndKeyHex]
			}
		}

		egs = append(egs, eg)
	}
	return ess, ees, egs
}
