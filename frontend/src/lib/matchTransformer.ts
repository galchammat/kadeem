import { models } from "@wails/go/models"
import { GetPlayerRankAtTime } from "@wails/go/riot/RiotClient"
import * as DataDragonClient from "@wails/go/datadragon/DataDragonClient"
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

// Convert bytes to data URL
function bytesToDataURL(bytes: number[]): string {
  const binary = String.fromCharCode(...bytes)
  return `data:image/png;base64,${btoa(binary)}`
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
function formatRank(rank: models.PlayerRank | null): string {
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
  backendMatch: models.LeagueOfLegendsMatch,
  trackedPUUID: string
): Promise<Match> {
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

  // Get champion icon
  const champIconBytes = await DataDragonClient.GetChampionIcon(playerParticipant.championId)
  const champIcon = bytesToDataURL(Array.from(champIconBytes))

  // Get summoner spell icons
  const spell1Bytes = await DataDragonClient.GetSummonerSpellIcon(
    playerParticipant.summoner1Id
  )
  const spell2Bytes = await DataDragonClient.GetSummonerSpellIcon(
    playerParticipant.summoner2Id
  )
  const summonerSpells = [
    bytesToDataURL(Array.from(spell1Bytes)),
    bytesToDataURL(Array.from(spell2Bytes)),
  ]

  // Get item icons
  const itemIds = [
    playerParticipant.item0,
    playerParticipant.item1,
    playerParticipant.item2,
    playerParticipant.item3,
    playerParticipant.item4,
    playerParticipant.item5,
  ]
  const itemPromises = itemIds.map(async (itemId) => {
    if (itemId === 0) return "/placeholder.svg" // Empty slot
    const bytes = await DataDragonClient.GetItemIcon(itemId)
    return bytesToDataURL(Array.from(bytes))
  })
  const items = await Promise.all(itemPromises)

  // Get rank at match time
  const queueId = backendMatch.summary.queueId || 420 // Default to Ranked Solo/Duo
  let rank: models.PlayerRank | null = null
  try {
    rank = await GetPlayerRankAtTime(
      trackedPUUID,
      queueId,
      backendMatch.summary.startedAt || 0
    )
  } catch (error) {
    console.warn(`Failed to get rank for match ${backendMatch.summary.gameId}:`, error)
  }

  // Calculate team kills for participation
  const playerTeam = playerParticipant.win ? 100 : 200 // Assumes team 100 = blue, team 200 = red
  // Note: We don't have teamId in participant data, so we'll use win status as proxy
  // In reality, you'd want to track teamId properly
  const teamParticipants = backendMatch.participants.filter((p) => p.win === playerParticipant.win)
  const teamKills = teamParticipants.reduce((sum, p) => sum + p.kills, 0)

  // Calculate CS per minute
  const csPerMin =
    backendMatch.summary.duration && backendMatch.summary.duration > 0
      ? (playerParticipant.totalMinionsKilled / (backendMatch.summary.duration / 60)).toFixed(1)
      : "0.0"

  // Build teams data (simplified - we'll show all participants)
  // Team 100 = Blue, Team 200 = Red
  // Since we don't have teamId in participant, we'll split by win status as approximation
  const winners = backendMatch.participants.filter((p) => p.win)
  const losers = backendMatch.participants.filter((p) => !p.win)

  const buildTeamData = async (participants: models.LeagueOfLegendsMatchParticipantSummary[]) => {
    return Promise.all(
      participants.slice(0, 5).map(async (p) => {
        const champBytes = await DataDragonClient.GetChampionIcon(p.championId)
        return {
          name: p.riotIdGameName || "Unknown",
          champion: bytesToDataURL(Array.from(champBytes)),
        }
      })
    )
  }

  // Assuming blue = winners, red = losers (this is a simplification)
  const blueTeam = playerParticipant.win ? winners : losers
  const redTeam = playerParticipant.win ? losers : winners

  const teams = {
    blue: await buildTeamData(blueTeam),
    red: await buildTeamData(redTeam),
  }

  return {
    id: backendMatch.summary.gameId,
    queueType: QUEUE_NAMES[queueId] || "Unknown Queue",
    timeAgo: formatTimeAgo(backendMatch.summary.startedAt || 0),
    result: playerParticipant.win ? "Victory" : "Defeat",
    duration: formatDuration(backendMatch.summary.duration || 0),
    champion: {
      name: "", // We don't have champion name from backend, could add later
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
      laning: "N/A", // Skip laning for MVP
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
