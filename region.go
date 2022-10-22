package main

type Region struct {
	StartKey string `json:"start_key"`
	EndKey   string `json:"end_key"`
	Peers    []struct {
		StoreID  int64  `json:"store_id"`
		RoleName string `json:"role_name"`
	} `json:"peers"`
	Leader struct {
		StoreID int64 `json:"store_id"`
	} `json:"leader"`
}

func loadRegions() ([]Region, error) {
	type regionData struct {
		Regions []Region `json:"regions"`
	}
	var rd regionData
	err := httpGet("http://"+*pdAddr+"/pd/api/v1/regions", &rd)
	if err != nil {
		return nil, err
	}
	return rd.Regions, nil
}
