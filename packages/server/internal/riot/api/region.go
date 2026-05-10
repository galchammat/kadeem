package api

var apiRegionMap = map[string]string{
	"NA1":  "americas",
	"EUW1": "europe",
	"EUN1": "europe",
	"KR":   "asia",
	"JP1":  "asia",
	"BR1":  "americas",
	"LA1":  "americas",
	"LA2":  "americas",
	"OC1":  "sea",
	"RU":   "europe",
	"TR1":  "europe",
}

func GetAPIRegion(region string) (string, error) {
	apiRegion, exists := apiRegionMap[region]
	if !exists {
		return region, nil // default region
	}
	return apiRegion, nil
}
