package datadragon

import (
	"fmt"
	"strconv"
)

// GetChampionIDByName returns the champion ID for a given name
// Uses the cached champion map, loading it if necessary
func (c *DataDragonClient) GetChampionIDByName(name string) (int, error) {
	c.championMapOnce.Do(func() {
		_ = c.loadChampionMap()
	})

	c.championMapMu.RLock()
	defer c.championMapMu.RUnlock()

	// Search through map for matching name
	for id, championName := range c.championIDMap {
		if championName == name {
			return id, nil
		}
	}

	return 0, fmt.Errorf("champion '%s' not found", name)
}

// GetItemIDByName returns the item ID for a given name.
// If multiple items have the same name (e.g., arena variants), returns the lowest ID (base item).
func (c *DataDragonClient) GetItemIDByName(name string) (int, error) {
	itemData, err := c.GetItemData("en_US")
	if err != nil {
		return 0, err
	}

	var foundID int
	for idStr, item := range itemData.Data {
		if item.Name == name {
			id, err := strconv.Atoi(idStr)
			if err != nil {
				return 0, err
			}
			// Prefer lower ID (base item over variants)
			if foundID == 0 || id < foundID {
				foundID = id
			}
		}
	}

	if foundID > 0 {
		return foundID, nil
	}

	return 0, fmt.Errorf("item '%s' not found", name)
}

// GetSummonerSpellIDByName returns the summoner spell ID for a given name.
// If multiple spells have the same name (e.g., arena variants), returns the lowest ID (base spell).
func (c *DataDragonClient) GetSummonerSpellIDByName(name string) (int, error) {
	c.spellMapOnce.Do(func() {
		_ = c.loadSpellMap()
	})

	c.spellMapMu.RLock()
	defer c.spellMapMu.RUnlock()

	// Get spell data to match names
	spellData, err := c.GetSummonerSpellData("en_US")
	if err != nil {
		return 0, err
	}

	var foundID int
	// Search through map for matching name
	for id, spellName := range c.spellIDMap {
		// Get the actual display name from the spell data
		for _, spell := range spellData.Data {
			if spell.ID == spellName && spell.Name == name {
				// Prefer lower ID (base spell over variants)
				if foundID == 0 || id < foundID {
					foundID = id
				}
			}
		}
	}

	if foundID > 0 {
		return foundID, nil
	}

	return 0, fmt.Errorf("summoner spell '%s' not found", name)
}

// GetPerkIDByName returns the perk/rune ID for a given name
func (c *DataDragonClient) GetPerkIDByName(name string) (int, error) {
	runeData, err := c.GetRuneData("en_US")
	if err != nil {
		return 0, err
	}

	for _, tree := range runeData {
		for _, slot := range tree.Slots {
			for _, rune := range slot.Runes {
				if rune.Name == name {
					return rune.ID, nil
				}
			}
		}
	}

	return 0, fmt.Errorf("perk '%s' not found", name)
}

// GetPerkTreeIDByName returns the perk tree ID for a given name
func (c *DataDragonClient) GetPerkTreeIDByName(name string) (int, error) {
	runeData, err := c.GetRuneData("en_US")
	if err != nil {
		return 0, err
	}

	for _, tree := range runeData {
		if tree.Name == name {
			return tree.ID, nil
		}
	}

	return 0, fmt.Errorf("perk tree '%s' not found", name)
}

// GetChampionIDsByNames returns a map of champion names to IDs
func (c *DataDragonClient) GetChampionIDsByNames(names []string) (map[string]int, error) {
	results := make(map[string]int)

	for _, name := range names {
		id, err := c.GetChampionIDByName(name)
		if err == nil {
			results[name] = id
		}
	}

	return results, nil
}
