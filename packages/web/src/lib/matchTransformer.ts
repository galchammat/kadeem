import type {
  LeagueOfLegendsMatch,
  LeagueOfLegendsMatchParticipantSummary,
  PlayerRank,
  ChampionData,
  SummonerSpellData,
  ItemData,
} from "@/types"
import * as api from "@/lib/api"
import type { Match } from "@/components/matchCard"

// Queue ID to name mapping
const QUEUE_NAMES: Record<number, string> = {
  420: "Ranked Solo/Duo",
  440: "Ranked Flex",
  400: "Normal Draft",
  430: "Normal Blind",
  450: "ARAM",
  700: "Clash",
  900: "URF",
  1020: "One For All",
  1300: "Nexus Blitz",
  1400: "Ultimate Spellbook",
  1700: "Arena",
}

// Cached DataDragon data
let cachedVersion: string | null = null
let cachedChampions: ChampionData | null = null
let cachedItems: ItemData | null = null
let cachedSpells: SummonerSpellData | null = null

// ID lookup maps (built lazily from cached data)
let championIdMap: Map<number, string> | null = null
let itemIdMap: Map<number, string> | null = null
let spellIdMap: Map<number, string> | null = null

async function ensureDataDragonData() {
  if (!cachedVersion) {
    cachedVersion = await api.getDataDragonVersion()
  }
  if (!cachedChampions) {
    cachedChampions = await api.getChampionData()
    championIdMap = new Map()
    for (const champ of Object.values(cachedChampions.data)) {
      championIdMap.set(Number(champ.key), champ.id)
    }
  }
  if (!cachedItems) {
    cachedItems = await api.getItemData()
    itemIdMap = new Map()
    for (const [idStr, item] of Object.entries(cachedItems.data)) {
      itemIdMap.set(Number(idStr), item.image.full)
    }
  }
  if (!cachedSpells) {
    cachedSpells = await api.getSummonerSpellData()
    spellIdMap = new Map()
    for (const spell of Object.values(cachedSpells.data)) {
      spellIdMap.set(Number(spell.key), spell.id)
    }
  }
}

function cdnUrl(path: string): string {
  return `https://ddragon.leagueoflegends.com/cdn/${cachedVersion}/${path}`
}

function championIconUrl(championId: number): string {
  const name = championIdMap?.get(championId)
  if (!name) return "/placeholder.svg"
  return cdnUrl(`img/champion/${name}.png`)
}

function itemIconUrl(itemId: number): string {
  if (itemId === 0) return "/placeholder.svg"
  const imageName = itemIdMap?.get(itemId)
  if (!imageName) return "/placeholder.svg"
  return cdnUrl(`img/item/${imageName}`)
}

function spellIconUrl(spellId: number): string {
  const name = spellIdMap?.get(spellId)
  if (!name) return "/placeholder.svg"
  return cdnUrl(`img/spell/${name}.png`)
}

// Format timestamp to relative time
function formatTimeAgo(timestamp: number): string {
  const now = Date.now()
  const diffMs = now - timestamp * 1000 // timestamp is in seconds
  const diffMins = Math.floor(diffMs / 60000)
  const diffHours = Math.floor(diffMs / 3600000)
  const diffDays = Math.floor(diffMs / 86400000)

  if (diffMins < 60) {
    return `${diffMins} minute${diffMins === 1 ? "" : "s"} ago`
  } else if (diffHours < 24) {
    return `${diffHours} hour${diffHours === 1 ? "" : "s"} ago`
  } else {
    return `${diffDays} day${diffDays === 1 ? "" : "s"} ago`
  }
}

// Format duration
function formatDuration(seconds: number): string {
  const mins = Math.floor(seconds / 60)
  const secs = seconds % 60
  return `${mins}m ${secs}s`
}

// Calculate KDA ratio
function calculateKdaRatio(kills: number, deaths: number, assists: number): string {
  if (deaths === 0) return "Perfect"
  const ratio = ((kills + assists) / deaths).toFixed(2)
  return `${ratio}:1`
}

// Format rank string
function formatRank(rank: PlayerRank | null): string {
  if (!rank) return "Unranked"

  const tier = rank.tier.charAt(0) + rank.tier.slice(1).toLowerCase()

  // Master+ tiers don't have divisions
  if (["MASTER", "GRANDMASTER", "CHALLENGER"].includes(rank.tier)) {
    return tier
  }

  return `${tier} ${rank.rank}`
}

// Determine performance tag based on placement and result
function determinePerformanceTag(placement: number, isVictory: boolean): string {
  if (placement <= 3 && isVictory) return "Carry"
  if (placement >= 8 && !isVictory) return "Struggle"
  return "Average"
}

// Calculate participation percentage
function calculateParticipation(
  playerKills: number,
  playerAssists: number,
  teamKills: number
): string {
  if (teamKills === 0) return "0%"
  const participation = ((playerKills + playerAssists) / teamKills) * 100
  return `${Math.round(participation)}%`
}

// Transform backend match to frontend Match interface
export async function transformMatch(
  backendMatch: LeagueOfLegendsMatch,
  trackedPUUID: string,
  accountId: string
): Promise<Match> {
  await ensureDataDragonData()

  // Find the tracked player's participant data
  const playerParticipant = backendMatch.participants.find(
    (p) => p.puuid === trackedPUUID
  )

  if (!playerParticipant) {
    throw new Error(`Tracked PUUID ${trackedPUUID} not found in match participants`)
  }

  // Sort participants by KDA to determine placement
  const sortedParticipants = [...backendMatch.participants].sort((a, b) => {
    const kdaA = a.deaths === 0 ? 999 : (a.kills + a.assists) / a.deaths
    const kdaB = b.deaths === 0 ? 999 : (b.kills + b.assists) / b.deaths
    return kdaB - kdaA
  })
  const placement = sortedParticipants.findIndex((p) => p.puuid === trackedPUUID) + 1

  // Get icons via CDN URLs
  const champIcon = championIconUrl(playerParticipant.championId)
  const summonerSpells = [
    spellIconUrl(playerParticipant.summoner1Id),
    spellIconUrl(playerParticipant.summoner2Id),
  ]

  // Get item icons
  const items = [
    itemIconUrl(playerParticipant.item0),
    itemIconUrl(playerParticipant.item1),
    itemIconUrl(playerParticipant.item2),
    itemIconUrl(playerParticipant.item3),
    itemIconUrl(playerParticipant.item4),
    itemIconUrl(playerParticipant.item5),
  ]

  // Get rank at match time
  const queueId = backendMatch.summary.queueId || 420
  let rank: PlayerRank | null = null
  try {
    rank = await api.getPlayerRankAtTime(
      accountId,
      queueId,
      backendMatch.summary.startedAt || 0
    )
  } catch (error) {
    console.warn(`Failed to get rank for match ${backendMatch.summary.gameId}:`, error)
  }

  // Calculate team kills for participation
  const teamParticipants = backendMatch.participants.filter((p) => p.win === playerParticipant.win)
  const teamKills = teamParticipants.reduce((sum, p) => sum + p.kills, 0)

  // Calculate CS per minute
  const csPerMin =
    backendMatch.summary.duration && backendMatch.summary.duration > 0
      ? (playerParticipant.totalMinionsKilled / (backendMatch.summary.duration / 60)).toFixed(1)
      : "0.0"

  // Build teams data
  const winners = backendMatch.participants.filter((p) => p.win)
  const losers = backendMatch.participants.filter((p) => !p.win)

  const buildTeamData = (participants: LeagueOfLegendsMatchParticipantSummary[]) => {
    return participants.slice(0, 5).map((p) => ({
      name: p.riotIdGameName || "Unknown",
      champion: championIconUrl(p.championId),
    }))
  }

  const blueTeam = playerParticipant.win ? winners : losers
  const redTeam = playerParticipant.win ? losers : winners

  const teams = {
    blue: buildTeamData(blueTeam),
    red: buildTeamData(redTeam),
  }

  return {
    id: backendMatch.summary.gameId,
    startedAt: backendMatch.summary.startedAt ?? 0,
    durationSeconds: backendMatch.summary.duration ?? 0,
    queueType: QUEUE_NAMES[queueId] || "Unknown Queue",
    timeAgo: formatTimeAgo(backendMatch.summary.startedAt || 0),
    result: playerParticipant.win ? "Victory" : "Defeat",
    duration: formatDuration(backendMatch.summary.duration || 0),
    champion: {
      name: "",
      image: champIcon,
      level: playerParticipant.champLevel || 1,
    },
    kda: {
      kills: playerParticipant.kills,
      deaths: playerParticipant.deaths,
      assists: playerParticipant.assists,
    },
    kdaRatio: calculateKdaRatio(
      playerParticipant.kills,
      playerParticipant.deaths,
      playerParticipant.assists
    ),
    summonerSpells,
    items,
    stats: {
      laning: "N/A",
      pKill: calculateParticipation(
        playerParticipant.kills,
        playerParticipant.assists,
        teamKills
      ),
      cs: `${playerParticipant.totalMinionsKilled} (${csPerMin})`,
      rank: formatRank(rank),
    },
    placement: `${placement}${getOrdinalSuffix(placement)}`,
    performanceTag: determinePerformanceTag(placement, playerParticipant.win),
    teams,
  }
}

// Helper to get ordinal suffix (1st, 2nd, 3rd, etc.)
function getOrdinalSuffix(num: number): string {
  const j = num % 10
  const k = num % 100
  if (j === 1 && k !== 11) return "st"
  if (j === 2 && k !== 12) return "nd"
  if (j === 3 && k !== 13) return "rd"
  return "th"
}
