package riot

import (
	"fmt"
)

var apiRegionMap = map[string]string{
	"NA":   "americas",
	"EUW":  "europe",
	"EUNE": "europe",
	"KR":   "asia",
	"JP":   "asia",
	"BR":   "americas",
	"LAN":  "americas",
	"LAS":  "americas",
	"OCE":  "sea",
	"RU":   "europe",
	"TR":   "europe",
}

func GetAPIRegion(region string) (string, error) {
	apiRegion, exists := apiRegionMap[region]
	if !exists {
		return "", fmt.Errorf("unsupported region: %s", region)
	}
	return apiRegion, nil
}
