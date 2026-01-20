package datadragon_test

import (
	"context"
	"fmt"
	"os"

	"github.com/galchammat/kadeem/internal/riot/datadragon"
)

// ExampleNewClient demonstrates how to create a new Data Dragon client
func ExampleNewDataDragonClient() {
	ctx := context.Background()

	// Create client with default cache directory (./bin/datadragon)
	client := datadragon.NewDataDragonClient(ctx, "")

	fmt.Println("Client version:", client.GetVersion())
}

// ExampleDataDragonClient_GetChampionIcon demonstrates fetching a single champion icon by ID
func ExampleDataDragonClient_GetChampionIcon() {
	ctx := context.Background()
	client := datadragon.NewDataDragonClient(ctx, "")

	// Fetch Ahri's icon (ID: 103)
	iconData, err := client.GetChampionIcon(103)
	if err != nil {
		panic(err)
	}

	// Save to file
	err = os.WriteFile("ahri_icon.png", iconData, 0644)
	if err != nil {
		panic(err)
	}

	fmt.Println("Ahri icon saved to ahri_icon.png")
}

// ExampleDataDragonClient_GetChampionIDByName demonstrates looking up ID by name
func ExampleDataDragonClient_GetChampionIDByName() {
	ctx := context.Background()
	client := datadragon.NewDataDragonClient(ctx, "")

	// If you only have a name, look up ID first
	championID, err := client.GetChampionIDByName("Ahri")
	if err != nil {
		panic(err)
	}

	// Then fetch the icon
	iconData, err := client.GetChampionIcon(championID)
	if err != nil {
		panic(err)
	}

	// Save to file
	err = os.WriteFile("ahri_icon.png", iconData, 0644)
	if err != nil {
		panic(err)
	}

	fmt.Println("Ahri icon saved to ahri_icon.png")
}

// ExampleDataDragonClient_GetItemIcon demonstrates fetching a single item icon
func ExampleDataDragonClient_GetItemIcon() {
	ctx := context.Background()
	client := datadragon.NewDataDragonClient(ctx, "")

	// Fetch item icon (Infinity Edge - item ID 3031)
	iconData, err := client.GetItemIcon(3031)
	if err != nil {
		panic(err)
	}

	// Save to file
	err = os.WriteFile("item_3031.png", iconData, 0644)
	if err != nil {
		panic(err)
	}

	fmt.Println("Item icon saved to item_3031.png")
}

// ExampleDataDragonClient_BatchFetchChampionIcons demonstrates fetching multiple icons at once
func ExampleDataDragonClient_BatchFetchChampionIcons() {
	ctx := context.Background()
	client := datadragon.NewDataDragonClient(ctx, "")

	// Fetch 25 champion icons concurrently by ID
	championIDs := []int{
		103, 84, 12, 32, 34, // Ahri, Akali, Alistar, Amumu, Anivia
		1, 22, 268, 432, 53, // Annie, Ashe, Azir, Bard, Blitzcrank
		63, 201, 51, 69, 31, // Brand, Braum, Caitlyn, Cassiopeia, Cho'Gath
		42, 122, 131, 119, 36, // Corki, Darius, Diana, Draven, Dr. Mundo
		245, 60, 28, 81, 9, // Ekko, Elise, Evelynn, Ezreal, Fiddlesticks
	}

	results, err := client.BatchFetchChampionIcons(championIDs)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Successfully fetched %d champion icons\n", len(results))

	// Save each icon
	for championID, iconData := range results {
		if iconData != nil {
			filename := fmt.Sprintf("champion_%d.png", championID)
			err = os.WriteFile(filename, iconData, 0644)
			if err != nil {
				fmt.Printf("Failed to save %s: %v\n", filename, err)
			}
		}
	}
}

// ExampleDataDragonClient_GetPerkIcon demonstrates fetching a keystone perk icon
func ExampleDataDragonClient_GetPerkIcon() {
	ctx := context.Background()
	client := datadragon.NewDataDragonClient(ctx, "")

	// Fetch a perk icon (Electrocute - ID: 8112)
	iconData, err := client.GetPerkIcon(8112)
	if err != nil {
		panic(err)
	}

	// Save to file
	err = os.WriteFile("electrocute.png", iconData, 0644)
	if err != nil {
		panic(err)
	}

	fmt.Println("Perk icon saved to electrocute.png")
}

// ExampleDataDragonClient_GetPerkTreeIcon demonstrates fetching a perk tree icon (for secondary path)
func ExampleDataDragonClient_GetPerkTreeIcon() {
	ctx := context.Background()
	client := datadragon.NewDataDragonClient(ctx, "")

	// Fetch Sorcery tree icon (ID: 8200)
	iconData, err := client.GetPerkTreeIcon(8200)
	if err != nil {
		panic(err)
	}

	// Save to file
	err = os.WriteFile("sorcery_tree.png", iconData, 0644)
	if err != nil {
		panic(err)
	}

	fmt.Println("Perk tree icon saved to sorcery_tree.png")
}

// ExampleDataDragonClient_GetSummonerSpellIcon demonstrates fetching a summoner spell icon
func ExampleDataDragonClient_GetSummonerSpellIcon() {
	ctx := context.Background()
	client := datadragon.NewDataDragonClient(ctx, "")

	// Fetch Flash icon (ID: 4)
	iconData, err := client.GetSummonerSpellIcon(4)
	if err != nil {
		panic(err)
	}

	// Save to file
	err = os.WriteFile("flash.png", iconData, 0644)
	if err != nil {
		panic(err)
	}

	fmt.Println("Flash icon saved to flash.png")
}

// ExampleDataDragonClient_BatchFetchItemIcons demonstrates fetching multiple item icons
func ExampleDataDragonClient_BatchFetchItemIcons() {
	ctx := context.Background()
	client := datadragon.NewDataDragonClient(ctx, "")

	// Fetch multiple item icons by ID
	itemIDs := []int{
		3031, // Infinity Edge
		3153, // Blade of the Ruined King
		3089, // Rabadon's Deathcap
		3078, // Trinity Force
		3742, // Dead Man's Plate
	}

	results, err := client.BatchFetchItemIcons(itemIDs)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Successfully fetched %d item icons\n", len(results))
}
