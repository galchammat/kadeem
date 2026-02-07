package datadragon

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetChampionData(t *testing.T) {
	ctx := context.Background()
	tempDir := t.TempDir()

	client := NewDataDragonClient(ctx, tempDir)

	// Fetch champion data
	championData, err := client.GetChampionData("en_US")
	require.NoError(t, err)
	require.NotNil(t, championData)

	// Verify structure
	assert.NotEmpty(t, championData.Version)
	assert.NotEmpty(t, championData.Data)

	// Check for some well-known champions
	assert.Contains(t, championData.Data, "Ahri")
	assert.Contains(t, championData.Data, "Ashe")

	// Verify champion data structure
	ahri := championData.Data["Ahri"]
	assert.Equal(t, "Ahri", ahri.ID)
	assert.NotEmpty(t, ahri.Name)
	assert.NotEmpty(t, ahri.Key) // Numeric ID
	assert.NotEmpty(t, ahri.Title)
	assert.NotEmpty(t, ahri.Image.Full)
}

func TestGetItemData(t *testing.T) {
	ctx := context.Background()
	tempDir := t.TempDir()

	client := NewDataDragonClient(ctx, tempDir)

	// Fetch item data
	itemData, err := client.GetItemData("en_US")
	require.NoError(t, err)
	require.NotNil(t, itemData)

	// Verify structure
	assert.NotEmpty(t, itemData.Version)
	assert.NotEmpty(t, itemData.Data)

	// Check for a well-known item (Infinity Edge - 3031)
	assert.Contains(t, itemData.Data, "3031")

	infinityEdge := itemData.Data["3031"]
	assert.NotEmpty(t, infinityEdge.Name)
	assert.NotEmpty(t, infinityEdge.Image.Full)
	assert.Greater(t, infinityEdge.Gold.Total, 0)
}

func TestGetRuneData(t *testing.T) {
	ctx := context.Background()
	tempDir := t.TempDir()

	client := NewDataDragonClient(ctx, tempDir)

	// Fetch rune data
	runeTrees, err := client.GetRuneData("en_US")
	require.NoError(t, err)
	require.NotEmpty(t, runeTrees)

	// There should be 5 rune trees (Precision, Domination, Sorcery, Resolve, Inspiration)
	assert.GreaterOrEqual(t, len(runeTrees), 5)

	// Verify structure of first tree
	tree := runeTrees[0]
	assert.NotEmpty(t, tree.Name)
	assert.NotEmpty(t, tree.Key)
	assert.NotEmpty(t, tree.Icon)
	assert.NotEmpty(t, tree.Slots)

	// Each tree should have slots with runes
	for _, slot := range tree.Slots {
		assert.NotEmpty(t, slot.Runes)
		for _, rune := range slot.Runes {
			assert.NotEmpty(t, rune.Name)
			assert.NotEmpty(t, rune.Icon)
			assert.Greater(t, rune.ID, 0)
		}
	}
}

func TestGetSummonerSpellData(t *testing.T) {
	ctx := context.Background()
	tempDir := t.TempDir()

	client := NewDataDragonClient(ctx, tempDir)

	// Fetch summoner spell data
	spellData, err := client.GetSummonerSpellData("en_US")
	require.NoError(t, err)
	require.NotNil(t, spellData)

	// Verify structure
	assert.NotEmpty(t, spellData.Version)
	assert.NotEmpty(t, spellData.Data)

	// Check for well-known summoner spell (Flash)
	assert.Contains(t, spellData.Data, "SummonerFlash")

	flash := spellData.Data["SummonerFlash"]
	assert.Equal(t, "4", flash.Key)
	assert.Equal(t, "Flash", flash.Name)
	assert.NotEmpty(t, flash.Image.Full)
}

func TestDataCaching(t *testing.T) {
	ctx := context.Background()
	tempDir := t.TempDir()

	client := NewDataDragonClient(ctx, tempDir)

	// Fetch champion data twice
	data1, err := client.GetChampionData("en_US")
	require.NoError(t, err)

	data2, err := client.GetChampionData("en_US")
	require.NoError(t, err)

	// Both should return the same data
	assert.Equal(t, data1.Version, data2.Version)
	assert.Equal(t, len(data1.Data), len(data2.Data))
}

func TestMultipleLocales(t *testing.T) {
	ctx := context.Background()
	tempDir := t.TempDir()

	client := NewDataDragonClient(ctx, tempDir)

	// Fetch data in different locales
	enData, err := client.GetChampionData("en_US")
	require.NoError(t, err)

	koData, err := client.GetChampionData("ko_KR")
	require.NoError(t, err)

	// Both should have the same champions
	assert.Equal(t, len(enData.Data), len(koData.Data))

	// Names might be different (localized)
	// But IDs and Keys should be the same
	assert.Contains(t, enData.Data, "Ahri")
	assert.Contains(t, koData.Data, "Ahri")

	enAhri := enData.Data["Ahri"]
	koAhri := koData.Data["Ahri"]
	assert.Equal(t, enAhri.Key, koAhri.Key)
	assert.Equal(t, enAhri.ID, koAhri.ID)
}
