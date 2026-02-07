package datadragon

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/galchammat/kadeem/internal/logging"
)

const (
	versionsURL = "https://ddragon.leagueoflegends.com/api/versions.json"
	cdnBaseURL  = "https://ddragon.leagueoflegends.com/cdn"
)

// Client handles Data Dragon requests with local caching
type DataDragonClient struct {
	ctx        context.Context
	httpClient *http.Client
	version    string
	cacheDir   string
	cacheMu    sync.RWMutex

	// In-memory ID→Name/Path mappings (lazy loaded)
	championIDMap   map[int]string
	championMapOnce sync.Once
	championMapMu   sync.RWMutex

	itemIDMap   map[int]string
	itemMapOnce sync.Once
	itemMapMu   sync.RWMutex

	perkIDMap     map[int]string
	perkTreeIDMap map[int]string
	perkMapOnce   sync.Once
	perkMapMu     sync.RWMutex

	spellIDMap   map[int]string
	spellMapOnce sync.Once
	spellMapMu   sync.RWMutex
}

// NewClient creates a new Data Dragon client
// It fetches the latest version on startup and sets up local caching
func NewDataDragonClient(ctx context.Context, cacheDir string) *DataDragonClient {
	if cacheDir == "" {
		cacheDir = "./bin/datadragon"
	}

	client := &DataDragonClient{
		ctx: ctx,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		cacheDir:      cacheDir,
		championIDMap: make(map[int]string),
		itemIDMap:     make(map[int]string),
		perkIDMap:     make(map[int]string),
		perkTreeIDMap: make(map[int]string),
		spellIDMap:    make(map[int]string),
	}

	// Fetch the latest version on startup
	if err := client.updateVersion(); err != nil {
		panic(fmt.Errorf("failed to fetch Data Dragon version: %w", err))
	}

	logging.Info("Data Dragon client initialized", "version", client.version, "cache", cacheDir)

	return client
}

// updateVersion fetches the latest Data Dragon version
func (c *DataDragonClient) updateVersion() error {
	req, err := http.NewRequestWithContext(c.ctx, "GET", versionsURL, nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch versions: status %d", resp.StatusCode)
	}

	var versions []string
	if err := json.NewDecoder(resp.Body).Decode(&versions); err != nil {
		return err
	}

	if len(versions) == 0 {
		return fmt.Errorf("no versions available")
	}

	newVersion := versions[0] // First element is the latest version

	// If version changed, clear old cache
	if c.version != "" && c.version != newVersion {
		logging.Info("New Data Dragon version detected, clearing old cache", "old", c.version, "new", newVersion)
		_ = c.clearCache()
	}

	c.version = newVersion
	return nil
}

// GetVersion returns the current Data Dragon version
func (c *DataDragonClient) GetVersion() string {
	return c.version
}

// clearCache removes all cached files
func (c *DataDragonClient) clearCache() error {
	c.cacheMu.Lock()
	defer c.cacheMu.Unlock()

	return os.RemoveAll(c.cacheDir)
}

// getCacheDir returns the cache directory for the current version
func (c *DataDragonClient) getCacheDir(subdir string) string {
	return filepath.Join(c.cacheDir, c.version, subdir)
}

// getCachePath returns the full path for a cached file
func (c *DataDragonClient) getCachePath(subdir, filename string) string {
	return filepath.Join(c.getCacheDir(subdir), filename)
}

// fetchAndCache downloads a file from a URL and caches it locally
func (c *DataDragonClient) fetchAndCache(url, subdir, filename string) ([]byte, error) {
	cachePath := c.getCachePath(subdir, filename)

	// Check if file exists in cache
	c.cacheMu.RLock()
	if data, err := os.ReadFile(cachePath); err == nil {
		c.cacheMu.RUnlock()
		return data, nil
	}
	c.cacheMu.RUnlock()

	// Download the file
	req, err := http.NewRequestWithContext(c.ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch %s: status %d", url, resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Cache the file
	c.cacheMu.Lock()
	defer c.cacheMu.Unlock()

	cacheDir := c.getCacheDir(subdir)
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		logging.Warn("Failed to create cache directory", "error", err)
		return data, nil // Return data even if caching fails
	}

	if err := os.WriteFile(cachePath, data, 0644); err != nil {
		logging.Warn("Failed to cache file", "path", cachePath, "error", err)
	}

	return data, nil
}

// loadChampionMap fetches champion.json and builds ID→Name map
func (c *DataDragonClient) loadChampionMap() error {
	championData, err := c.GetChampionData("en_US")
	if err != nil {
		return err
	}

	c.championMapMu.Lock()
	defer c.championMapMu.Unlock()

	for _, champion := range championData.Data {
		id, err := strconv.Atoi(champion.Key)
		if err != nil {
			logging.Warn("Failed to parse champion ID", "key", champion.Key, "champion", champion.ID)
			continue
		}
		c.championIDMap[id] = champion.ID
	}

	return nil
}

// getChampionName returns champion name for ID, loading map if needed
func (c *DataDragonClient) getChampionName(championID int) (string, error) {
	c.championMapOnce.Do(func() {
		if err := c.loadChampionMap(); err != nil {
			logging.Error("Failed to load champion map", "error", err)
		}
	})

	c.championMapMu.RLock()
	defer c.championMapMu.RUnlock()

	name, ok := c.championIDMap[championID]
	if !ok {
		return "", nil // ID not found
	}
	return name, nil
}

// loadItemMap fetches item.json and builds ID→ImageName map
func (c *DataDragonClient) loadItemMap() error {
	itemData, err := c.GetItemData("en_US")
	if err != nil {
		return err
	}

	c.itemMapMu.Lock()
	defer c.itemMapMu.Unlock()

	for idStr, item := range itemData.Data {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			logging.Warn("Failed to parse item ID", "key", idStr)
			continue
		}
		c.itemIDMap[id] = item.Image.Full
	}

	return nil
}

// getItemImageName returns item image filename for ID, loading map if needed
func (c *DataDragonClient) getItemImageName(itemID int) (string, error) {
	c.itemMapOnce.Do(func() {
		if err := c.loadItemMap(); err != nil {
			logging.Error("Failed to load item map", "error", err)
		}
	})

	c.itemMapMu.RLock()
	defer c.itemMapMu.RUnlock()

	imageName, ok := c.itemIDMap[itemID]
	if !ok {
		return "", nil // ID not found
	}
	return imageName, nil
}

// loadPerkMap fetches runesReforged.json and builds ID→IconPath maps
func (c *DataDragonClient) loadPerkMap() error {
	runeData, err := c.GetRuneData("en_US")
	if err != nil {
		return err
	}

	c.perkMapMu.Lock()
	defer c.perkMapMu.Unlock()

	// Map tree IDs to their icons
	for _, tree := range runeData {
		c.perkTreeIDMap[tree.ID] = tree.Icon

		// Map individual perk IDs to their icons
		for _, slot := range tree.Slots {
			for _, rune := range slot.Runes {
				c.perkIDMap[rune.ID] = rune.Icon
			}
		}
	}

	return nil
}

// getPerkIconPath returns perk icon path for ID, loading map if needed
func (c *DataDragonClient) getPerkIconPath(perkID int) (string, error) {
	c.perkMapOnce.Do(func() {
		if err := c.loadPerkMap(); err != nil {
			logging.Error("Failed to load perk map", "error", err)
		}
	})

	c.perkMapMu.RLock()
	defer c.perkMapMu.RUnlock()

	iconPath, ok := c.perkIDMap[perkID]
	if !ok {
		return "", nil // ID not found
	}
	return iconPath, nil
}

// getPerkTreeIconPath returns perk tree icon path for ID, loading map if needed
func (c *DataDragonClient) getPerkTreeIconPath(treeID int) (string, error) {
	c.perkMapOnce.Do(func() {
		if err := c.loadPerkMap(); err != nil {
			logging.Error("Failed to load perk map", "error", err)
		}
	})

	c.perkMapMu.RLock()
	defer c.perkMapMu.RUnlock()

	iconPath, ok := c.perkTreeIDMap[treeID]
	if !ok {
		return "", nil // ID not found
	}
	return iconPath, nil
}

// loadSpellMap fetches summoner.json and builds ID→Name map
func (c *DataDragonClient) loadSpellMap() error {
	spellData, err := c.GetSummonerSpellData("en_US")
	if err != nil {
		return err
	}

	c.spellMapMu.Lock()
	defer c.spellMapMu.Unlock()

	for _, spell := range spellData.Data {
		id, err := strconv.Atoi(spell.Key)
		if err != nil {
			logging.Warn("Failed to parse summoner spell ID", "key", spell.Key, "spell", spell.ID)
			continue
		}
		c.spellIDMap[id] = spell.ID
	}

	return nil
}

// getSummonerSpellName returns summoner spell name for ID, loading map if needed
func (c *DataDragonClient) getSummonerSpellName(spellID int) (string, error) {
	c.spellMapOnce.Do(func() {
		if err := c.loadSpellMap(); err != nil {
			logging.Error("Failed to load summoner spell map", "error", err)
		}
	})

	c.spellMapMu.RLock()
	defer c.spellMapMu.RUnlock()

	name, ok := c.spellIDMap[spellID]
	if !ok {
		return "", nil // ID not found
	}
	return name, nil
}

// GetChampionIcon fetches a champion icon by champion ID
func (c *DataDragonClient) GetChampionIcon(championID int) ([]byte, error) {
	name, err := c.getChampionName(championID)
	if err != nil {
		return nil, err
	}
	if name == "" {
		return nil, nil // ID not found
	}

	url := fmt.Sprintf("%s/%s/img/champion/%s.png", cdnBaseURL, c.version, name)
	return c.fetchAndCache(url, "champions", fmt.Sprintf("%d.png", championID))
}

// GetItemIcon fetches an item icon by item ID
func (c *DataDragonClient) GetItemIcon(itemID int) ([]byte, error) {
	imageName, err := c.getItemImageName(itemID)
	if err != nil {
		return nil, err
	}
	if imageName == "" {
		return nil, nil // ID not found
	}

	url := fmt.Sprintf("%s/%s/img/item/%s", cdnBaseURL, c.version, imageName)
	return c.fetchAndCache(url, "items", fmt.Sprintf("%d.png", itemID))
}

// GetPerkIcon fetches a perk/rune icon by perk ID (for keystones)
// GetPerkIcon fetches a perk icon by perk ID (for keystones like Electrocute)
func (c *DataDragonClient) GetPerkIcon(perkID int) ([]byte, error) {
	iconPath, err := c.getPerkIconPath(perkID)
	if err != nil {
		return nil, err
	}
	if iconPath == "" {
		return nil, nil // ID not found
	}

	// Perk icons don't use version in URL
	url := fmt.Sprintf("%s/img/%s", cdnBaseURL, iconPath)
	filename := fmt.Sprintf("perk_%d.png", perkID)
	return c.fetchAndCache(url, "perks", filename)
}

// GetPerkTreeIcon fetches a perk tree icon by tree ID (for secondary path)
func (c *DataDragonClient) GetPerkTreeIcon(treeID int) ([]byte, error) {
	iconPath, err := c.getPerkTreeIconPath(treeID)
	if err != nil {
		return nil, err
	}
	if iconPath == "" {
		return nil, nil // ID not found
	}

	// Perk tree icons don't use version in URL
	url := fmt.Sprintf("%s/img/%s", cdnBaseURL, iconPath)
	filename := fmt.Sprintf("tree_%d.png", treeID)
	return c.fetchAndCache(url, "perks", filename)
}

// GetSummonerSpellIcon fetches a summoner spell icon by spell ID
func (c *DataDragonClient) GetSummonerSpellIcon(spellID int) ([]byte, error) {
	spellName, err := c.getSummonerSpellName(spellID)
	if err != nil {
		return nil, err
	}
	if spellName == "" {
		return nil, nil // ID not found
	}

	url := fmt.Sprintf("%s/%s/img/spell/%s.png", cdnBaseURL, c.version, spellName)
	return c.fetchAndCache(url, "spells", fmt.Sprintf("%d.png", spellID))
}

// BatchFetchChampionIcons fetches multiple champion icons concurrently
func (c *DataDragonClient) BatchFetchChampionIcons(championIDs []int) (map[int][]byte, error) {
	results := make(map[int][]byte)
	var resultsMu sync.Mutex
	var wg sync.WaitGroup

	for _, id := range championIDs {
		wg.Add(1)
		go func(championID int) {
			defer wg.Done()

			data, err := c.GetChampionIcon(championID)
			if err != nil {
				logging.Warn("Failed to fetch champion icon", "id", championID, "error", err)
			}

			resultsMu.Lock()
			results[championID] = data // Store nil if not found
			resultsMu.Unlock()
		}(id)
	}

	wg.Wait()
	return results, nil
}

// BatchFetchItemIcons fetches multiple item icons concurrently
func (c *DataDragonClient) BatchFetchItemIcons(itemIDs []int) (map[int][]byte, error) {
	results := make(map[int][]byte)
	var resultsMu sync.Mutex
	var wg sync.WaitGroup

	for _, id := range itemIDs {
		wg.Add(1)
		go func(itemID int) {
			defer wg.Done()

			data, err := c.GetItemIcon(itemID)
			if err != nil {
				logging.Warn("Failed to fetch item icon", "id", itemID, "error", err)
			}

			resultsMu.Lock()
			results[itemID] = data
			resultsMu.Unlock()
		}(id)
	}

	wg.Wait()
	return results, nil
}

// BatchFetchPerkIcons fetches multiple perk icons concurrently
func (c *DataDragonClient) BatchFetchPerkIcons(perkIDs []int) (map[int][]byte, error) {
	results := make(map[int][]byte)
	var resultsMu sync.Mutex
	var wg sync.WaitGroup

	for _, id := range perkIDs {
		wg.Add(1)
		go func(perkID int) {
			defer wg.Done()

			data, err := c.GetPerkIcon(perkID)
			if err != nil {
				logging.Warn("Failed to fetch perk icon", "id", perkID, "error", err)
			}

			resultsMu.Lock()
			results[perkID] = data
			resultsMu.Unlock()
		}(id)
	}

	wg.Wait()
	return results, nil
}

// BatchFetchPerkTreeIcons fetches multiple perk tree icons concurrently
func (c *DataDragonClient) BatchFetchPerkTreeIcons(treeIDs []int) (map[int][]byte, error) {
	results := make(map[int][]byte)
	var resultsMu sync.Mutex
	var wg sync.WaitGroup

	for _, id := range treeIDs {
		wg.Add(1)
		go func(treeID int) {
			defer wg.Done()

			data, err := c.GetPerkTreeIcon(treeID)
			if err != nil {
				logging.Warn("Failed to fetch perk tree icon", "id", treeID, "error", err)
			}

			resultsMu.Lock()
			results[treeID] = data
			resultsMu.Unlock()
		}(id)
	}

	wg.Wait()
	return results, nil
}

// BatchFetchSummonerSpellIcons fetches multiple summoner spell icons concurrently
func (c *DataDragonClient) BatchFetchSummonerSpellIcons(spellIDs []int) (map[int][]byte, error) {
	results := make(map[int][]byte)
	var resultsMu sync.Mutex
	var wg sync.WaitGroup

	for _, id := range spellIDs {
		wg.Add(1)
		go func(spellID int) {
			defer wg.Done()

			data, err := c.GetSummonerSpellIcon(spellID)
			if err != nil {
				logging.Warn("Failed to fetch summoner spell icon", "id", spellID, "error", err)
			}

			resultsMu.Lock()
			results[spellID] = data
			resultsMu.Unlock()
		}(id)
	}

	wg.Wait()
	return results, nil
}
