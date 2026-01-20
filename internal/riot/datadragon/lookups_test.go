package datadragon

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetChampionIDByName(t *testing.T) {
	ctx := context.Background()
	tempDir := t.TempDir()

	client := NewDataDragonClient(ctx, tempDir)

	// Test exact match
	id, err := client.GetChampionIDByName("Ahri")
	require.NoError(t, err)
	assert.Equal(t, 103, id)

	// Test not found
	_, err = client.GetChampionIDByName("InvalidChampion")
	assert.Error(t, err)
}

func TestGetSummonerSpellIDByName(t *testing.T) {
	ctx := context.Background()
	tempDir := t.TempDir()

	client := NewDataDragonClient(ctx, tempDir)

	id, err := client.GetSummonerSpellIDByName("Flash")
	require.NoError(t, err)
	assert.Equal(t, 4, id)
}

func TestGetPerkIDByName(t *testing.T) {
	ctx := context.Background()
	tempDir := t.TempDir()

	client := NewDataDragonClient(ctx, tempDir)

	id, err := client.GetPerkIDByName("Electrocute")
	require.NoError(t, err)
	assert.Equal(t, 8112, id)
}

func TestGetPerkTreeIDByName(t *testing.T) {
	ctx := context.Background()
	tempDir := t.TempDir()

	client := NewDataDragonClient(ctx, tempDir)

	id, err := client.GetPerkTreeIDByName("Domination")
	require.NoError(t, err)
	assert.Equal(t, 8100, id)
}

func TestGetChampionIDsByNames(t *testing.T) {
	ctx := context.Background()
	tempDir := t.TempDir()

	client := NewDataDragonClient(ctx, tempDir)

	names := []string{"Ahri", "Ashe", "InvalidChamp"}
	results, err := client.GetChampionIDsByNames(names)

	require.NoError(t, err)
	assert.Equal(t, 103, results["Ahri"])
	assert.Equal(t, 22, results["Ashe"])
	assert.NotContains(t, results, "InvalidChamp")
}

func TestGetItemIDByName(t *testing.T) {
	ctx := context.Background()
	tempDir := t.TempDir()

	client := NewDataDragonClient(ctx, tempDir)

	// Test a well-known item
	id, err := client.GetItemIDByName("Infinity Edge")
	require.NoError(t, err)
	assert.Equal(t, 3031, id)
}
