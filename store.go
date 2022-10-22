package main

import (
	"sort"
)

type StoreLabel struct {
	Key    string   `json:"key"`
	Values []string `json:"values"`
}

type Store struct {
	ID      int64  `json:"id"`
	Address string `json:"address"`
	Labels  []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"labels"`
}

func loadStores() ([]Store, []StoreLabel, error) {
	type storesData struct {
		Stores []struct {
			Store Store `json:"store"`
		} `json:"stores"`
	}
	var sd storesData
	err := httpGet("http://"+*pdAddr+"/pd/api/v1/stores", &sd)
	if err != nil {
		return nil, nil, err
	}

	var stores []Store
	for _, store := range sd.Stores {
		stores = append(stores, store.Store)
	}
	sort.Slice(stores, func(i, j int) bool {
		return stores[i].ID < stores[j].ID
	})

	labels := make(map[string]map[string]struct{})
	for _, store := range sd.Stores {
		for _, label := range store.Store.Labels {
			if _, ok := labels[label.Key]; !ok {
				labels[label.Key] = make(map[string]struct{})
			}
			labels[label.Key][label.Value] = struct{}{}
		}
	}

	var storeLabels []StoreLabel
	for key, vs := range labels {
		var values []string
		for k := range vs {
			values = append(values, k)
		}
		sort.Strings(values)
		storeLabels = append(storeLabels, StoreLabel{
			Key:    key,
			Values: values,
		})
		sort.Slice(storeLabels, func(i, j int) bool {
			return storeLabels[i].Key < storeLabels[j].Key
		})
	}

	return stores, storeLabels, nil
}
