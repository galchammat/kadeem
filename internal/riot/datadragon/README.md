# Data Dragon Client

A Go package for accessing Riot's Data Dragon CDN to fetch League of Legends static assets (champion icons, item icons, perk/rune icons, summoner spell icons).

## Features

- **Integer ID API**: Uses int IDs as the primary interface (champion IDs, item IDs, perk IDs, spell IDs)
- **Automatic Version Management**: Fetches the latest Data Dragon version on startup
- **Local Caching**: Caches assets locally to minimize network requests (stored in `./bin/datadragon`)
- **Cache Invalidation**: Automatically clears cache when a new patch is detected
- **Concurrent Fetching**: Batch fetch multiple icons simultaneously for optimal performance
- **Thread-Safe**: Safe for concurrent access across goroutines
- **Lookup Helpers**: Optional name→ID conversion when you only have names

## Installation

```go
import "github.com/galchammat/kadeem/internal/riot/datadragon"
```

## Usage

### Initialize Client

```go
ctx := context.Background()

// Create client with default cache directory (./bin/datadragon)
client, err := datadragon.NewClient(ctx, "")
if err != nil {
    log.Fatal(err)
}

fmt.Println("Data Dragon version:", client.GetVersion())
```

### Fetch Champion Icons

```go
// Single champion (by champion ID)
iconData, err := client.GetChampionIcon(103) // Ahri
if err != nil {
    log.Fatal(err)
}

// Invalid IDs return nil (no error)
iconData := client.GetChampionIcon(99999) // returns nil

// Save to file
os.WriteFile("ahri_icon.png", iconData, 0644)
```

**Common Champion IDs:**
- 103 = Ahri
- 22 = Ashe
- 1 = Annie
- 51 = Caitlyn
- 67 = Vayne
- 498 = Xayah

### Fetch Item Icons

```go
// Single item (by item ID)
iconData, err := client.GetItemIcon(3031) // Infinity Edge
if err != nil {
    log.Fatal(err)
}

// Invalid IDs return nil (no error)
iconData := client.GetItemIcon(99999) // returns nil
```

**Common Item IDs:**
- 3031 = Infinity Edge
- 3153 = Blade of the Ruined King
- 3157 = Zhonya's Hourglass
- 3068 = Sunfire Aegis

### Fetch Perk/Rune Icons

```go
// Fetch keystone rune icon (e.g., Electrocute)
iconData, err := client.GetPerkIcon(8112) // Electrocute
if err != nil {
    log.Fatal(err)
}

// Fetch secondary rune tree icon (e.g., Sorcery)
iconData, err := client.GetPerkTreeIcon(8200) // Sorcery
if err != nil {
    log.Fatal(err)
}
```

**Perk IDs (Keystones):**
- 8112 = Electrocute (Domination)
- 8128 = Dark Harvest (Domination)
- 8214 = Summon Aery (Sorcery)
- 8229 = Arcane Comet (Sorcery)
- 8005 = Press the Attack (Precision)

**Perk Tree IDs:**
- 8000 = Precision
- 8100 = Domination
- 8200 = Sorcery
- 8300 = Resolve
- 8400 = Inspiration

### Fetch Summoner Spell Icons

```go
// Fetch summoner spell icon
iconData, err := client.GetSummonerSpellIcon(4) // Flash
if err != nil {
    log.Fatal(err)
}
```

**Common Summoner Spell IDs:**
- 4 = Flash
- 14 = Ignite
- 12 = Teleport
- 7 = Heal
- 3 = Exhaust
- 21 = Barrier
- 11 = Smite

### Batch Fetching

Efficiently fetch multiple icons at once using concurrent requests:

```go
championIDs := []int{103, 22, 1, 51, 67} // Ahri, Ashe, Annie, Caitlyn, Vayne

results, err := client.BatchFetchChampionIcons(championIDs)
if err != nil {
    log.Fatal(err)
}

for champID, iconData := range results {
    filename := fmt.Sprintf("%d.png", champID)
    os.WriteFile(filename, iconData, 0644)
}
```

Batch methods return `map[int][]byte`:
- `BatchFetchChampionIcons([]int) (map[int][]byte, error)`
- `BatchFetchItemIcons([]int) (map[int][]byte, error)`
- `BatchFetchPerkIcons([]int) (map[int][]byte, error)`
- `BatchFetchPerkTreeIcons([]int) (map[int][]byte, error)`
- `BatchFetchSummonerSpellIcons([]int) (map[int][]byte, error)`

### Name Lookup Helpers

If you only have names instead of IDs, use the lookup helpers:

```go
// Champion name → ID
championID, err := client.GetChampionIDByName("Ahri")
// Returns: 103

// Item name → ID
itemID, err := client.GetItemIDByName("Infinity Edge")
// Returns: 3031

// Summoner spell name → ID
spellID, err := client.GetSummonerSpellIDByName("Flash")
// Returns: 4

// Perk name → ID
perkID, err := client.GetPerkIDByName("Electrocute")
// Returns: 8112

// Perk tree name → ID
treeID, err := client.GetPerkTreeIDByName("Sorcery")
// Returns: 8200

// Batch champion lookup
names := []string{"Ahri", "Ashe", "Annie"}
results, err := client.GetChampionIDsByNames(names)
// Returns: map[string]int{"Ahri": 103, "Ashe": 22, "Annie": 1}
```

**Note:** For items and summoner spells with duplicate names (e.g., arena variants), lookup helpers return the **lowest ID** (base item/spell).

## Caching

Assets are cached locally in a version-specific directory structure:

```
./bin/datadragon/
└── 16.1.1/
    ├── champions/
    │   ├── 103.png
    │   └── 22.png
    ├── items/
    │   ├── 3031.png
    │   └── 3153.png
    ├── perks/
    │   ├── perk_8112.png
    │   └── tree_8200.png
    └── spells/
        ├── 4.png
        └── 14.png
```

When a new patch is detected:
1. The client fetches the new version number
2. Old cache is cleared automatically
3. New assets are cached in the new version directory

## Data Dragon URLs

The package uses Riot's official Data Dragon CDN:

- **Versions**: `https://ddragon.leagueoflegends.com/api/versions.json`
- **Champions**: `https://ddragon.leagueoflegends.com/cdn/{version}/img/champion/{name}.png`
- **Items**: `https://ddragon.leagueoflegends.com/cdn/{version}/img/item/{id}.png`
- **Perks**: `https://ddragon.leagueoflegends.com/cdn/img/perk-images/Styles/{path}` (no version)
- **Summoner Spells**: `https://ddragon.leagueoflegends.com/cdn/{version}/img/spell/{name}.png`

**Note:** Perk icons (both individual perks and trees) don't use a version number in their URLs.

## Champion IDs

Champion IDs are numeric and correspond to their internal Riot API IDs. You can find all champion IDs via:
- The Riot API champion endpoint
- The lookup helper: `client.GetChampionIDByName("ChampionName")`
- [champion.json](http://ddragon.leagueoflegends.com/cdn/16.1.1/data/en_US/champion.json)

## Error Handling

- **Invalid IDs**: Methods return `nil` (no error) for invalid IDs
- **Network errors**: Return errors for failed HTTP requests or data parsing issues
- **Lookup failures**: Name→ID helpers return errors when name not found

```go
// Invalid ID returns nil without error
icon := client.GetChampionIcon(99999) // icon == nil

// Network error returns error
icon, err := client.GetChampionIcon(103)
if err != nil {
    log.Fatal(err) // Handle network/API error
}
```

## Testing

```bash
# Run all tests
go test ./internal/riot/datadragon/

# Run tests with coverage
go test -cover ./internal/riot/datadragon/

# Run tests in short mode (skip long-running tests)
go test -short ./internal/riot/datadragon/

# Run specific test
go test -v ./internal/riot/datadragon/ -run TestBatchFetch
```

## References

- [Data Dragon Documentation](https://riot-api-libraries.readthedocs.io/en/latest/ddragon.html)
- [Riot Developer Portal](https://developer.riotgames.com/)
- [Community Dragon](http://raw.communitydragon.org/) - Additional assets not in Data Dragon
