package main

import (
	"fmt"
	"math/rand"
	"regexp"
	"sort"
	"strings"
	"time"
)

func convertGroupMarkdown(groups []Group) string {
	var sb strings.Builder
	for _, g := range groups {
		sb.WriteString("group " + g.ID + "\n\n")
		for _, r := range g.Rules {
			sb.WriteString(fmt.Sprintf("%s - %s: %s\n", TSMAP.StartTS[r.StartKeyHex].Format("1/2/2006"), TSMAP.EndTS[r.EndKeyHex].Format("1/2/2006"), r.ID))
		}
	}
	return sb.String()
}

func convertRegionMarkdown(regions []Region) string {
	peers := make(map[int64][]string)
	for _, s := range STORES {
		peers[s.ID] = make([]string, 0)
	}
	for _, s := range SCHEMAS {
		for _, r := range regions {
			if s.StartKey >= r.StartKey && (r.EndKey == "" || s.EndKey <= r.EndKey) {
				leader := r.Leader.StoreID
				for _, p := range r.Peers {
					role := "Follower"
					if p.StoreID == leader {
						role = "Leader"
					}
					if p.RoleName == "Learner" {
						role = "Learner"
					}
					start := TSMAP.StartTS[s.StartKey].Format("1/2/2006") // + " 08:00:00"
					end := TSMAP.StartTS[s.StartKey].Format("1/2/2006")   // + " 16:00:00"
					peers[p.StoreID] = append(peers[p.StoreID], fmt.Sprintf("%s - %s: %s", start, end, role))
				}
				break
			}
		}
	}
	var storeIDs []int64
	for id := range peers {
		storeIDs = append(storeIDs, id)
	}
	sort.Slice(storeIDs, func(i, j int) bool {
		return storeIDs[i] < storeIDs[j]
	})

	var sb strings.Builder
	for _, id := range storeIDs {
		desc := ""
		for _, s := range STORES {
			if id == s.ID {
				for _, l := range s.Labels {
					desc += l.Value + " "
				}
			}
		}
		sb.WriteString(fmt.Sprintf("section Store%d %s\n\n", id, desc))
		for _, r := range peers[id] {
			sb.WriteString(r + "\n")
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func mockRegionMarkdown() string {
	var storeIDs []int64
	for _, s := range STORES {
		storeIDs = append(storeIDs, s.ID)
	}
	peers := make([][]string, len(storeIDs))
	for _, s := range SCHEMAS {
		var idxs []int
		for i := range storeIDs {
			idxs = append(idxs, i)
		}
		rand.Shuffle(len(idxs), func(i, j int) { idxs[i], idxs[j] = idxs[j], idxs[i] })
		i := 0
		peers[idxs[i]] = append(peers[idxs[i]], fmt.Sprintf("%s: Leader", TSMAP.StartTS[s.StartKey].Format("1/2/2006")))
		i++
		for j := 0; j < rand.Int()%4; j++ {
			peers[idxs[i]] = append(peers[idxs[i]], fmt.Sprintf("%s: Follower", TSMAP.StartTS[s.StartKey].Format("1/2/2006")))
			i++
		}
		for j := 0; j < rand.Int()%2; j++ {
			peers[idxs[i]] = append(peers[idxs[i]], fmt.Sprintf("%s: Learner", TSMAP.StartTS[s.StartKey].Format("1/2/2006")))
			i++
		}
	}

	var sb strings.Builder
	for i := range storeIDs {
		sb.WriteString(fmt.Sprintf("section Store #%d\n\n", storeIDs[i]))
		for _, r := range peers[i] {
			sb.WriteString(r + "\n")
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func parseMarkdown(data string) ([]Group, error) {
	var groups []Group
	var group Group
	lines := strings.Split(data, "\n")
	var index int
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		if strings.HasPrefix(line, "group") {
			if group.ID != "" {
				groups = append(groups, group)
			}
			line = strings.TrimPrefix(line, "group")
			line = strings.TrimSpace(line)
			group = Group{ID: line, Index: index}
			index++
			continue
		}
		var start, end, id string
		regx := regexp.MustCompile(`(\d+/\d+/\d+)\s+-\s+(\d+/\d+/\d+).*:(.*)`)
		matches := regx.FindStringSubmatch(line)
		if len(matches) == 4 {
			start, end, id = matches[1], matches[2], matches[3]
		} else {
			regx := regexp.MustCompile(`(\d+/\d+/\d+).*:(.*)`)
			matches := regx.FindStringSubmatch(line)
			if len(matches) == 3 {
				start, end, id = matches[1], matches[1], matches[2]
			} else {
				return nil, fmt.Errorf("invalid event format: %s", line)
			}
		}
		rule := Rule{
			GroupID:     group.ID,
			ID:          strings.TrimSpace(id),
			Index:       index,
			StartKeyHex: findStartKey(start),
			EndKeyHex:   findEndKey(end),
		}
		index++
		group.Rules = append(group.Rules, rule)
	}
	if group.ID != "" {
		groups = append(groups, group)
	}
	return groups, nil
}

func findStartKey(timeStr string) string {
	t, _ := time.ParseInLocation("1/2/2006", timeStr, time.Local)
	for k, v := range TSMAP.StartTS {
		if v.Unix() == t.Unix() {
			return k
		}
	}
	return ""
}

func findEndKey(timeStr string) string {
	t, _ := time.ParseInLocation("1/2/2006", timeStr, time.Local)
	for k, v := range TSMAP.EndTS {
		if v.Unix() == t.Unix() {
			return k
		}
	}
	return ""
}
