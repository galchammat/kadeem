export type DialogMode = null | "add" | "edit"

// Riot / League of Legends
export interface LolAccount {
  puuid: string
  streamerId?: number
  tagLine: string
  gameName: string
  region?: string
  syncedAt?: number | null
}

export interface LolMatch {
  summary: LolMatchSummary
  participants: LolMatchParticipantSummary[]
  replay?: string | null
}

export interface LolMatchSummary {
  gameId: number
  startedAt?: number | null
  duration?: number | null
  queueId?: number | null
  replaySynced?: boolean | null
}

export interface LolMatchParticipantSummary {
  gameId: number
  championId: number
  champLevel: number
  kills: number
  deaths: number
  assists: number
  totalMinionsKilled: number
  doubleKills: number
  tripleKills: number
  quadraKills: number
  pentaKills: number
  item0: number
  item1: number
  item2: number
  item3: number
  item4: number
  item5: number
  item6: number
  summoner1Id: number
  summoner2Id: number
  lane: string
  participantId: number
  puuid: string
  riotIdGameName: string
  riotIdTagline: string
  totalDamageDealtToChampions: number
  totalDamageTaken: number
  win: boolean
}

export interface PlayerRank {
  puuid: string
  timestamp: number
  tier: string
  rank: string
  leaguePoints: number
  wins: number
  losses: number
  queueId: number
}

// Streamers / Channels / Broadcasts
export interface StreamerView {
  id: number
  name: string
  channels: Channel[]
  lastLive?: number | null
  avatarUrl?: string | null
}

export interface Channel {
  id: string
  streamerId: number
  platform: string
  channelName: string
  avatarUrl: string
  syncedAt?: number | null
}

export interface Broadcast {
  id: number
  channel_id: string
  title: string
  url: string
  thumbnail_url: string
  viewable: string
  created_at: number
  published_at: number
  duration: number
}

export type EventType =
  | "hype_train"
  | "gifted_subs"
  | "peak_viewers"
  | "low_viewers"
  | "raid"
  | "donation"
  | "new_subscriber"
  | "chat_milestone"
  | "follower_goal"
  | "ban_wave"
  | "clip"

export interface StreamEvent {
  id: number
  type: EventType
  title: string
  description: string
  timestamp: number
  value?: string | number
}

// DataDragon
export interface ChampionData {
  type: string
  format: string
  version: string
  data: Record<string, Champion>
}

export interface Champion {
  version: string
  id: string
  key: string
  name: string
  title: string
  image: DDImage
}

export interface DDImage {
  full: string
  sprite: string
  group: string
  x: number
  y: number
  w: number
  h: number
}

export interface ItemData {
  type: string
  version: string
  data: Record<string, Item>
}

export interface Item {
  name: string
  description: string
  image: DDImage
}

export interface SummonerSpellData {
  type: string
  version: string
  data: Record<string, SummonerSpell>
}

export interface SummonerSpell {
  id: string
  key: string
  name: string
  description: string
  image: DDImage
}
