package main

type StoreLabel struct {
	Key    string   `json:"key"`
	Values []string `json:"values"`
}

func loadStoreLabels() ([]StoreLabel, error) {
	return []StoreLabel{
		{
			Key:    "engine",
			Values: []string{"tikv", "tiflash"},
		},
		{
			Key:    "zone",
			Values: []string{"zone-1", "zone-2", "zone-3"},
		},
		{
			Key:    "disk",
			Values: []string{"ssd", "hdd"},
		},
	}, nil
}
