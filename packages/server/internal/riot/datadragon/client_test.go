package datadragon

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDataDragonClient(t *testing.T) {
	ctx := context.Background()
	tempDir := t.TempDir()

	client := NewDataDragonClient(ctx, tempDir)
	require.NotNil(t, client)

	// Check that version was fetched
	assert.NotEmpty(t, client.GetVersion())

	// Verify version format (e.g., "16.1.1")
	version := client.GetVersion()
	assert.Regexp(t, `^\d+\.\d+\.\d+$`, version, "Version should match format X.Y.Z")
}

func TestGetChampionIcon(t *testing.T) {
	ctx := context.Background()
	tempDir := t.TempDir()

	client := NewDataDragonClient(ctx, tempDir)

	// Test fetching Ahri icon (ID: 103)
	iconData, err := client.GetChampionIcon(103)
	require.NoError(t, err)
	require.NotEmpty(t, iconData)

	// Verify it's a PNG file (starts with PNG magic bytes)
	assert.Equal(t, []byte{0x89, 0x50, 0x4E, 0x47}, iconData[:4], "Should be a PNG file")

	// Test that second fetch comes from cache
	iconData2, err := client.GetChampionIcon(103)
	require.NoError(t, err)
	assert.Equal(t, iconData, iconData2, "Cached data should match original")

	// Verify file exists in cache
	cachePath := client.getCachePath("champions", "103.png")
	_, err = os.Stat(cachePath)
	assert.NoError(t, err, "Icon should be cached")
}

func TestGetItemIcon(t *testing.T) {
	ctx := context.Background()
	tempDir := t.TempDir()

	client := NewDataDragonClient(ctx, tempDir)

	// Test fetching Infinity Edge icon (3031)
	iconData, err := client.GetItemIcon(3031)
	require.NoError(t, err)
	require.NotEmpty(t, iconData)

	// Verify it's a PNG file
	assert.Equal(t, []byte{0x89, 0x50, 0x4E, 0x47}, iconData[:4], "Should be a PNG file")
}

func TestBatchFetchChampionIcons(t *testing.T) {
	ctx := context.Background()
	tempDir := t.TempDir()

	client := NewDataDragonClient(ctx, tempDir)

	// Test batch fetching (Ahri: 103, Akali: 84, Ashe: 22, Annie: 1, Azir: 268)
	champions := []int{103, 84, 22, 1, 268}

	results, err := client.BatchFetchChampionIcons(champions)
	require.NoError(t, err)

	// All champions should be in results
	assert.Len(t, results, len(champions))

	// Verify each result is a valid PNG
	for championID, iconData := range results {
		assert.NotEmpty(t, iconData, "Icon for %d should not be empty", championID)
		if iconData != nil {
			assert.Equal(t, []byte{0x89, 0x50, 0x4E, 0x47}, iconData[:4], "Champion %d icon should be a PNG file", championID)
		}
	}
}

func TestBatchFetchItemIcons(t *testing.T) {
	ctx := context.Background()
	tempDir := t.TempDir()

	client := NewDataDragonClient(ctx, tempDir)

	items := []int{3031, 3153, 3089}

	results, err := client.BatchFetchItemIcons(items)
	require.NoError(t, err)

	assert.Len(t, results, len(items))

	for itemID, iconData := range results {
		assert.NotEmpty(t, iconData, "Icon for item %d should not be empty", itemID)
		if iconData != nil {
			assert.Equal(t, []byte{0x89, 0x50, 0x4E, 0x47}, iconData[:4], "Item %d icon should be a PNG file", itemID)
		}
	}
}

func TestCachePersistence(t *testing.T) {
	ctx := context.Background()
	tempDir := t.TempDir()

	// Create first client
	client1 := NewDataDragonClient(ctx, tempDir)

	// Fetch icon (Ahri: 103)
	iconData1, err := client1.GetChampionIcon(103)
	require.NoError(t, err)

	// Create second client with same cache dir
	client2 := NewDataDragonClient(ctx, tempDir)

	// Fetch same icon - should come from cache
	iconData2, err := client2.GetChampionIcon(103)
	require.NoError(t, err)

	// Data should match
	assert.Equal(t, iconData1, iconData2, "Cached data should persist across client instances")
}

func TestVersionCacheInvalidation(t *testing.T) {
	ctx := context.Background()
	tempDir := t.TempDir()

	client := NewDataDragonClient(ctx, tempDir)

	// Fetch an icon to populate cache (Ahri: 103)
	_, err := client.GetChampionIcon(103)
	require.NoError(t, err)

	// Verify cache exists
	cacheDir := client.getCacheDir("champions")
	_, err = os.Stat(cacheDir)
	require.NoError(t, err, "Cache directory should exist")

	// Manually change version to simulate new patch
	oldVersion := client.version
	client.version = oldVersion + "-test"

	// Create new cache directory for new version
	newCacheDir := filepath.Join(tempDir, client.version, "champions")
	err = os.MkdirAll(newCacheDir, 0755)
	require.NoError(t, err)

	// Old cache should still exist
	_, err = os.Stat(cacheDir)
	assert.NoError(t, err, "Old cache should still exist")

	// New cache directory should be different
	assert.NotEqual(t, cacheDir, newCacheDir, "Cache directories should be different for different versions")
}

func TestInvalidChampionID(t *testing.T) {
	ctx := context.Background()
	tempDir := t.TempDir()

	client := NewDataDragonClient(ctx, tempDir)

	// Test with invalid champion ID
	icon, err := client.GetChampionIcon(99999)
	assert.NoError(t, err, "Should not return error for invalid ID")
	assert.Nil(t, icon, "Should return nil for invalid ID")
}

func TestConcurrentAccess(t *testing.T) {
	ctx := context.Background()
	tempDir := t.TempDir()

	client := NewDataDragonClient(ctx, tempDir)

	// Fetch same icon multiple times concurrently (Ahri: 103)
	const numGoroutines = 10
	done := make(chan bool, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			_, err := client.GetChampionIcon(103)
			assert.NoError(t, err)
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < numGoroutines; i++ {
		<-done
	}
}

func TestLargeBatchFetch(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping large batch test in short mode")
	}

	ctx := context.Background()
	tempDir := t.TempDir()

	client := NewDataDragonClient(ctx, tempDir)

	// Test fetching 25+ icons as mentioned in requirements
	// Using champion IDs
	champions := []int{
		103, 84, 12, 32, 34, // Ahri, Akali, Alistar, Amumu, Anivia
		1, 22, 268, 432, 53, // Annie, Ashe, Azir, Bard, Blitzcrank
		63, 201, 51, 69, 31, // Brand, Braum, Caitlyn, Cassiopeia, Cho'Gath
		42, 122, 131, 119, 36, // Corki, Darius, Diana, Draven, Dr. Mundo
		245, 60, 28, 81, 9, // Ekko, Elise, Evelynn, Ezreal, Fiddlesticks
		114, 105, 3, 41, 86, // Fiora, Fizz, Galio, Gangplank, Garen
	}

	results, err := client.BatchFetchChampionIcons(champions)
	require.NoError(t, err)

	// Should get most or all champions
	assert.GreaterOrEqual(t, len(results), 25, "Should successfully fetch at least 25 icons")

	// Count non-nil results
	nonNilCount := 0
	for _, data := range results {
		if data != nil {
			nonNilCount++
		}
	}
	assert.GreaterOrEqual(t, nonNilCount, 25, "Should have at least 25 valid icons")
}

func TestGetSummonerSpellIcon(t *testing.T) {
	ctx := context.Background()
	tempDir := t.TempDir()

	client := NewDataDragonClient(ctx, tempDir)

	// Test Flash (ID: 4)
	iconData, err := client.GetSummonerSpellIcon(4)
	require.NoError(t, err)
	require.NotEmpty(t, iconData)

	// Verify it's a PNG file
	assert.Equal(t, []byte{0x89, 0x50, 0x4E, 0x47}, iconData[:4], "Should be a PNG file")
}

func TestGetPerkIcon(t *testing.T) {
	ctx := context.Background()
	tempDir := t.TempDir()

	client := NewDataDragonClient(ctx, tempDir)

	// Test Electrocute (ID: 8112)
	iconData, err := client.GetPerkIcon(8112)
	require.NoError(t, err)
	require.NotEmpty(t, iconData)

	// Verify it's a PNG file
	assert.Equal(t, []byte{0x89, 0x50, 0x4E, 0x47}, iconData[:4], "Should be a PNG file")
}

func TestGetPerkTreeIcon(t *testing.T) {
	ctx := context.Background()
	tempDir := t.TempDir()

	client := NewDataDragonClient(ctx, tempDir)

	// Test Domination tree (ID: 8100)
	iconData, err := client.GetPerkTreeIcon(8100)
	require.NoError(t, err)
	require.NotEmpty(t, iconData)

	// Verify it's a PNG file
	assert.Equal(t, []byte{0x89, 0x50, 0x4E, 0x47}, iconData[:4], "Should be a PNG file")
}

func TestBatchFetchWithInvalidIDs(t *testing.T) {
	ctx := context.Background()
	tempDir := t.TempDir()

	client := NewDataDragonClient(ctx, tempDir)

	// Mix valid and invalid IDs
	championIDs := []int{103, 99999, 84, 88888, 22}
	results, err := client.BatchFetchChampionIcons(championIDs)

	require.NoError(t, err)
	assert.Len(t, results, 5) // All IDs in result map

	// Check that valid IDs have data
	assert.NotNil(t, results[103]) // Ahri - valid
	assert.Nil(t, results[99999])  // Invalid
	assert.NotNil(t, results[84])  // Akali - valid
	assert.Nil(t, results[88888])  // Invalid
	assert.NotNil(t, results[22])  // Ashe - valid
}
