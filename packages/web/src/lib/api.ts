import type {
  LolAccount,
  LolMatch,
  StreamerView,
  PlayerRank,
  Broadcast,
  ChampionData,
  ItemData,
  SummonerSpellData,
  Channel,
  StreamEvent,
} from "@/types"
import { supabase } from "@/lib/supabase"

const API_BASE = `${import.meta.env.VITE_API_URL || ""}/api/v0`

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const { data: { session } } = await supabase.auth.getSession()
  
  const headers: HeadersInit = {
    "Content-Type": "application/json",
    ...init?.headers,
  }

  if (session?.access_token) {
    (headers as Record<string, string>)["Authorization"] = `Bearer ${session.access_token}`
  }

  const res = await fetch(`${API_BASE}${path}`, {
    ...init,
    headers,
  })
  if (!res.ok) {
    const body = await res.json().catch(() => ({}))
    throw new Error(body.error || `Request failed: ${res.status}`)
  }
  if (res.status === 204) return undefined as T
  return res.json()
}

// Riot Accounts
export async function listAccounts(): Promise<LolAccount[]> {
  const data = await request<{ accounts: LolAccount[]; count: number }>("/riot/accounts")
  return data.accounts ?? []
}

export async function addAccount(region: string, gameName: string, tagLine: string, streamerId: number): Promise<void> {
  await request("/riot/accounts", {
    method: "POST",
    body: JSON.stringify({ region, game_name: gameName, tag_line: tagLine, streamer_id: streamerId }),
  })
}

export async function updateAccount(accountId: string, region: string, gameName: string, tagLine: string): Promise<void> {
  await request(`/riot/accounts/${accountId}`, {
    method: "PUT",
    body: JSON.stringify({ region, game_name: gameName, tag_line: tagLine }),
  })
}

export async function deleteAccount(accountId: string): Promise<void> {
  await request(`/riot/accounts/${accountId}`, { method: "DELETE" })
}

// Riot Matches
export async function listMatches(puuid: string, limit: number, offset: number): Promise<LolMatch[]> {
  const params = new URLSearchParams({ puuid, limit: String(limit), offset: String(offset) })
  const data = await request<{ matches: LolMatch[]; count: number }>(`/riot/matches?${params}`)
  return data.matches ?? []
}

// Riot Ranks
export async function getPlayerRankAtTime(accountId: string, queueID: number, timestamp: number): Promise<PlayerRank | null> {
  const params = new URLSearchParams({ queueID: String(queueID), timestamp: String(timestamp) })
  return request<PlayerRank | null>(`/riot/accounts/${accountId}/rank-at-time?${params}`)
}

// Streamers
export async function listStreamers(): Promise<StreamerView[]> {
  const data = await request<{ streamers: StreamerView[]; count: number }>("/streamers")
  return data.streamers ?? []
}

export async function addStreamer(name: string): Promise<void> {
  await request("/streamers", {
    method: "POST",
    body: JSON.stringify({ name }),
  })
}

export async function deleteStreamer(name: string): Promise<void> {
  await request(`/streamers/${encodeURIComponent(name)}`, { method: "DELETE" })
}

// Channels
export async function addChannel(channel: Pick<Channel, "streamerId" | "channelName" | "platform">): Promise<void> {
  await request("/channels", {
    method: "POST",
    body: JSON.stringify({
      streamer_id: channel.streamerId,
      channel_name: channel.channelName,
      platform: channel.platform,
    }),
  })
}

export async function deleteChannel(channelId: string): Promise<void> {
  await request(`/channels/${encodeURIComponent(channelId)}`, { method: "DELETE" })
}

// Broadcasts
export async function listBroadcasts(channelId: string, limit: number, offset: number): Promise<Broadcast[]> {
  const params = new URLSearchParams({ channelID: channelId, limit: String(limit), offset: String(offset) })
  const data = await request<{ broadcasts: Broadcast[]; count: number }>(`/broadcasts?${params}`)
  return data.broadcasts ?? []
}

// Stream Events
export async function listStreamerEvents(streamerId: number, from: number, to: number, limit: number, offset: number): Promise<StreamEvent[]> {
  const params = new URLSearchParams({ from: String(from), to: String(to), limit: String(limit), offset: String(offset) })
  const data = await request<{ events: StreamEvent[]; count: number }>(`/streamers/${streamerId}/events?${params}`)
  return data.events ?? []
}

// DataDragon
export async function getDataDragonVersion(): Promise<string> {
  const data = await request<{ version: string }>("/datadragon/version")
  return data.version
}

export async function getChampionData(): Promise<ChampionData> {
  return request<ChampionData>("/datadragon/champions")
}

export async function getItemData(): Promise<ItemData> {
  return request<ItemData>("/datadragon/items")
}

export async function getSummonerSpellData(): Promise<SummonerSpellData> {
  return request<SummonerSpellData>("/datadragon/summoner-spells")
}
