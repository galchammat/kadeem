package datadragon

import (
	"encoding/json"
	"fmt"
)

// Image represents image metadata from Data Dragon
type Image struct {
	Full   string `json:"full"`
	Sprite string `json:"sprite"`
	Group  string `json:"group"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
	W      int    `json:"w"`
	H      int    `json:"h"`
}

// Champion represents basic champion data from Data Dragon
type Champion struct {
	Version string `json:"version"`
	ID      string `json:"id"`
	Key     string `json:"key"`
	Name    string `json:"name"`
	Title   string `json:"title"`
	Image   Image  `json:"image"`
}

// ChampionData represents the champion.json data structure
type ChampionData struct {
	Type    string              `json:"type"`
	Format  string              `json:"format"`
	Version string              `json:"version"`
	Data    map[string]Champion `json:"data"`
}

// GetChampionData fetches the champion.json data file
// This includes all champion info including IDs and names
func (c *DataDragonClient) GetChampionData(locale string) (*ChampionData, error) {
	if locale == "" {
		locale = "en_US"
	}

	url := fmt.Sprintf("%s/%s/data/%s/champion.json", cdnBaseURL, c.version, locale)

	// Use fetchAndCache to get the data
	data, err := c.fetchAndCache(url, "data", fmt.Sprintf("champion_%s.json", locale))
	if err != nil {
		return nil, err
	}

	var championData ChampionData
	if err := json.Unmarshal(data, &championData); err != nil {
		return nil, fmt.Errorf("failed to parse champion data: %w", err)
	}

	return &championData, nil
}

// Gold represents item gold data from Data Dragon
type Gold struct {
	Base        int  `json:"base"`
	Purchasable bool `json:"purchasable"`
	Total       int  `json:"total"`
	Sell        int  `json:"sell"`
}

// Item represents basic item data from Data Dragon
type Item struct {
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Colloq      string             `json:"colloq"`
	Plaintext   string             `json:"plaintext"`
	Into        []string           `json:"into"`
	Image       Image              `json:"image"`
	Gold        Gold               `json:"gold"`
	Tags        []string           `json:"tags"`
	Maps        map[string]bool    `json:"maps"`
	Stats       map[string]float64 `json:"stats"`
}

// ItemData represents the item.json data structure
type ItemData struct {
	Type    string          `json:"type"`
	Version string          `json:"version"`
	Data    map[string]Item `json:"data"`
}

// GetItemData fetches the item.json data file
func (c *DataDragonClient) GetItemData(locale string) (*ItemData, error) {
	if locale == "" {
		locale = "en_US"
	}

	url := fmt.Sprintf("%s/%s/data/%s/item.json", cdnBaseURL, c.version, locale)

	data, err := c.fetchAndCache(url, "data", fmt.Sprintf("item_%s.json", locale))
	if err != nil {
		return nil, err
	}

	var itemData ItemData
	if err := json.Unmarshal(data, &itemData); err != nil {
		return nil, fmt.Errorf("failed to parse item data: %w", err)
	}

	return &itemData, nil
}

// Rune represents rune data from Data Dragon
type Rune struct {
	ID        int    `json:"id"`
	Key       string `json:"key"`
	Icon      string `json:"icon"`
	Name      string `json:"name"`
	ShortDesc string `json:"shortDesc"`
	LongDesc  string `json:"longDesc"`
}

// RuneSlot represents a slot containing runes in a rune tree
type RuneSlot struct {
	Runes []Rune `json:"runes"`
}

// RuneTree represents a rune tree (path)
type RuneTree struct {
	ID    int        `json:"id"`
	Key   string     `json:"key"`
	Icon  string     `json:"icon"`
	Name  string     `json:"name"`
	Slots []RuneSlot `json:"slots"`
}

// GetRuneData fetches the runesReforged.json data file
func (c *DataDragonClient) GetRuneData(locale string) ([]RuneTree, error) {
	if locale == "" {
		locale = "en_US"
	}

	url := fmt.Sprintf("%s/%s/data/%s/runesReforged.json", cdnBaseURL, c.version, locale)

	data, err := c.fetchAndCache(url, "data", fmt.Sprintf("runes_%s.json", locale))
	if err != nil {
		return nil, err
	}

	var runeTrees []RuneTree
	if err := json.Unmarshal(data, &runeTrees); err != nil {
		return nil, fmt.Errorf("failed to parse rune data: %w", err)
	}

	return runeTrees, nil
}

// SummonerSpell represents summoner spell data from Data Dragon
type SummonerSpell struct {
	ID          string `json:"id"`
	Key         string `json:"key"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       Image  `json:"image"`
}

// SummonerSpellData represents the summoner.json data structure
type SummonerSpellData struct {
	Type    string                   `json:"type"`
	Version string                   `json:"version"`
	Data    map[string]SummonerSpell `json:"data"`
}

// GetSummonerSpellData fetches the summoner.json data file
func (c *DataDragonClient) GetSummonerSpellData(locale string) (*SummonerSpellData, error) {
	if locale == "" {
		locale = "en_US"
	}

	url := fmt.Sprintf("%s/%s/data/%s/summoner.json", cdnBaseURL, c.version, locale)

	data, err := c.fetchAndCache(url, "data", fmt.Sprintf("summoner_%s.json", locale))
	if err != nil {
		return nil, err
	}

	var spellData SummonerSpellData
	if err := json.Unmarshal(data, &spellData); err != nil {
		return nil, fmt.Errorf("failed to parse summoner spell data: %w", err)
	}

	return &spellData, nil
}
