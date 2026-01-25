import { MatchCard, type Match } from "@/components/matchCard"
import { MatchStreamEvents, type StreamEvent } from "@/components/matchStreamEvents"

import Sectionheader from '@/components/sectionHeader';
interface SessionEntry {
  match: Match
  events: StreamEvent[]
}

const sessionData: SessionEntry[] = [
  {

    match: {
      id: 1,
      queueType: "Ranked Solo/Duo",
      timeAgo: "18 hours ago",
      result: "Defeat",
      duration: "23m 35s",
      champion: {
        name: "Darius",
        image: "/darius-league-champion-portrait.jpg",
        level: 11,
      },
      kda: { kills: 1, deaths: 7, assists: 2 },
      kdaRatio: "0.43:1",
      summonerSpells: ["/flash-league-spell-icon.jpg", "/teleport-league-spell-icon.jpg"],
      items: [
        "/trinity-force-league-item.jpg",
        "/plated-steelcaps-league-item.jpg",
        "/deaths-dance-league-item.jpg",
        "/steraks-gage-league-item.jpg",
        "/guardian-angel-league-item.jpg",
        "/control-ward-league-item.jpg",
      ],
      stats: {
        laning: "37 : 63",
        pKill: "50%",
        cs: "170 (7.2)",
        rank: "Grandmaster",
      },
      placement: "9th",
      performanceTag: "Struggle",
      teams: {
        blue: [
          { name: "IMPACT", champion: "/gnar-league-champion.jpg" },
          { name: "rummiku...", champion: "/viego-league-champion.jpg" },
          { name: "kisno", champion: "/ahri-league-champion.jpg" },
          { name: "Morgan...", champion: "/jinx-league-champion.jpg" },
          { name: "Secret of...", champion: "/thresh-league-champion.jpg" },
        ],
        red: [
          { name: "Aivina", champion: "/darius-league-champion.jpg" },
          { name: "Dantes", champion: "/lee-sin-champion.jpg" },
          { name: "FREE PAL...", champion: "/syndra-champion.jpg" },
          { name: "EtooHai", champion: "/kaisa-champion.jpg" },
          { name: "Virtuosa", champion: "/nautilus-champion.jpg" },
        ],
      },
    },
    events: [
      {
        id: 1,
        type: "hype_train",
        title: "Hype Train Level 5!",
        description: "Community reached level 5 hype train",
        timestamp: "During match",
        value: "500%",
      },
      {
        id: 2,
        type: "peak_viewers",
        title: "Peak Viewers",
        description: "Highest concurrent viewers this session",
        timestamp: "18h ago",
        value: "12,847",
      },
      {
        id: 3,
        type: "gifted_subs",
        title: "Gift Sub Bomb!",
        description: "Anonymous gifted 50 subs",
        timestamp: "Mid-game",
        value: "50",
      },
      {
        id: 4,
        type: "chat_milestone",
        title: "Chat Exploded",
        description: "1,000+ messages per minute",
        timestamp: "Post-death",
        value: "1,234 msg/min",
      },
    ],
  },
  {
    match: {
      id: 2,
      queueType: "Ranked Solo/Duo",
      timeAgo: "19 hours ago",
      result: "Victory",
      duration: "31m 12s",
      champion: {
        name: "Garen",
        image: "/garen-league-champion.jpg",
        level: 16,
      },
      kda: { kills: 8, deaths: 3, assists: 12 },
      kdaRatio: "6.67:1",
      summonerSpells: ["/flash-league-spell-icon.jpg", "/ignite-spell.jpg"],
      items: [
        "/stridebreaker-item.jpg",
        "/dead-mans-plate-item.jpg",
        "/force-of-nature-item.jpg",
        "/mortal-reminder-item.jpg",
        "/warmogs-item.jpg",
        "/control-ward-league-item.jpg",
      ],
      stats: {
        laning: "58 : 42",
        pKill: "71%",
        cs: "245 (7.9)",
        rank: "Grandmaster",
      },
      placement: "2nd",
      performanceTag: "Carry",
      teams: {
        blue: [
          { name: "Player1", champion: "/ornn-champion.jpg" },
          { name: "Player2", champion: "/graves-champion.jpg" },
          { name: "Player3", champion: "/orianna-champion.jpg" },
          { name: "Player4", champion: "/ezreal-champion.jpg" },
          { name: "Player5", champion: "/nautilus-champion.jpg" },
        ],
        red: [
          { name: "Enemy1", champion: "/darius-league-champion.jpg" },
          { name: "Enemy2", champion: "/lee-sin-champion.jpg" },
          { name: "Enemy3", champion: "/syndra-champion.jpg" },
          { name: "Enemy4", champion: "/kaisa-champion.jpg" },
          { name: "Enemy5", champion: "/thresh-league-champion.jpg" },
        ],
      },
    },
    events: [
      {
        id: 5,
        type: "raid",
        title: "Incoming Raid!",
        description: "xQc raided with 8,500 viewers",
        timestamp: "Start of match",
        value: "8,500",
      },
      {
        id: 6,
        type: "donation",
        title: "Big Donation",
        description: "ChampViewer donated $500",
        timestamp: "19h ago",
        value: "$500",
      },
    ],
  },
  {
    match: {
      id: 3,
      queueType: "Ranked Solo/Duo",
      timeAgo: "20 hours ago",
      result: "Defeat",
      duration: "27m 45s",
      champion: {
        name: "Sett",
        image: "/sett-champion.jpg",
        level: 14,
      },
      kda: { kills: 4, deaths: 5, assists: 6 },
      kdaRatio: "2.00:1",
      summonerSpells: ["/flash-league-spell-icon.jpg", "/teleport-league-spell-icon.jpg"],
      items: [
        "/stridebreaker-item.jpg",
        "/plated-steelcaps-league-item.jpg",
        "/deaths-dance-league-item.jpg",
        "/steraks-gage-league-item.jpg",
        "/guardian-angel-league-item.jpg",
        "/control-ward-league-item.jpg",
      ],
      stats: {
        laning: "45 : 55",
        pKill: "42%",
        cs: "198 (7.1)",
        rank: "Grandmaster",
      },
      placement: "6th",
      performanceTag: "Average",
      teams: {
        blue: [
          { name: "TopMain", champion: "/gnar-league-champion.jpg" },
          { name: "JgDiff", champion: "/viego-league-champion.jpg" },
          { name: "MidGap", champion: "/ahri-league-champion.jpg" },
          { name: "ADCarry", champion: "/jinx-league-champion.jpg" },
          { name: "Support", champion: "/thresh-league-champion.jpg" },
        ],
        red: [
          { name: "OpTop", champion: "/darius-league-champion.jpg" },
          { name: "OpJg", champion: "/lee-sin-champion.jpg" },
          { name: "OpMid", champion: "/syndra-champion.jpg" },
          { name: "OpAdc", champion: "/kaisa-champion.jpg" },
          { name: "OpSup", champion: "/nautilus-champion.jpg" },
        ],
      },
    },
    events: [
      {
        id: 7,
        type: "low_viewers",
        title: "Viewer Dip",
        description: "Unusually low viewers detected",
        timestamp: "Mid-game",
        value: "2,341",
      },
      {
        id: 8,
        type: "new_subscriber",
        title: "New Sub Streak",
        description: "15 new subscribers in 5 minutes",
        timestamp: "20h ago",
        value: "15",
      },
      {
        id: 9,
        type: "follower_goal",
        title: "Follower Goal Hit!",
        description: "Reached 1M followers milestone",
        timestamp: "End of match",
        value: "1,000,000",
      },
      {
        id: 10,
        type: "ban_wave",
        title: "Mod Action",
        description: "23 users timed out for spam",
        timestamp: "Post-game",
        value: "23",
      },
      {
        id: 11,
        type: "gifted_subs",
        title: "Sub Train!",
        description: "Community gifted 25 subs",
        timestamp: "During fights",
        value: "25",
      },
      {
        id: 12,
        type: "chat_milestone",
        title: "Poggers Spam",
        description: "500+ emotes in 30 seconds",
        timestamp: "Team fight",
        value: "523 emotes",
      },
    ],
  },
  {
    match: {
      id: 4,
      queueType: "Ranked Solo/Duo",
      timeAgo: "21 hours ago",
      result: "Victory",
      duration: "25m 18s",
      champion: {
        name: "Gnar",
        image: "/gnar-league-champion.jpg",
        level: 15,
      },
      kda: { kills: 5, deaths: 2, assists: 14 },
      kdaRatio: "9.50:1",
      summonerSpells: ["/flash-league-spell-icon.jpg", "/teleport-league-spell-icon.jpg"],
      items: [
        "/stridebreaker-item.jpg",
        "/plated-steelcaps-league-item.jpg",
        "/force-of-nature-item.jpg",
        "/steraks-gage-league-item.jpg",
        "/warmogs-item.jpg",
        "/control-ward-league-item.jpg",
      ],
      stats: {
        laning: "62 : 38",
        pKill: "79%",
        cs: "210 (8.3)",
        rank: "Grandmaster",
      },
      placement: "1st",
      performanceTag: "Carry",
      teams: {
        blue: [
          { name: "Ally1", champion: "/gnar-league-champion.jpg" },
          { name: "Ally2", champion: "/viego-league-champion.jpg" },
          { name: "Ally3", champion: "/ahri-league-champion.jpg" },
          { name: "Ally4", champion: "/jinx-league-champion.jpg" },
          { name: "Ally5", champion: "/thresh-league-champion.jpg" },
        ],
        red: [
          { name: "Foe1", champion: "/darius-league-champion.jpg" },
          { name: "Foe2", champion: "/lee-sin-champion.jpg" },
          { name: "Foe3", champion: "/syndra-champion.jpg" },
          { name: "Foe4", champion: "/kaisa-champion.jpg" },
          { name: "Foe5", champion: "/nautilus-champion.jpg" },
        ],
      },
    },
    events: [
      {
        id: 13,
        type: "hype_train",
        title: "Hype Train Level 3!",
        description: "Community hyped up",
        timestamp: "Early game",
        value: "300%",
      },
    ],
  },
  {
    match: {
      id: 5,
      queueType: "Ranked Solo/Duo",
      timeAgo: "22 hours ago",
      result: "Defeat",
      duration: "35m 42s",
      champion: {
        name: "Viego",
        image: "/viego-league-champion.jpg",
        level: 17,
      },
      kda: { kills: 6, deaths: 8, assists: 5 },
      kdaRatio: "1.38:1",
      summonerSpells: ["/flash-league-spell-icon.jpg", "/ignite-spell.jpg"],
      items: [
        "/trinity-force-league-item.jpg",
        "/plated-steelcaps-league-item.jpg",
        "/deaths-dance-league-item.jpg",
        "/steraks-gage-league-item.jpg",
        "/guardian-angel-league-item.jpg",
        "/control-ward-league-item.jpg",
      ],
      stats: {
        laning: "48 : 52",
        pKill: "55%",
        cs: "267 (7.5)",
        rank: "Grandmaster",
      },
      placement: "7th",
      performanceTag: "Struggle",
      teams: {
        blue: [
          { name: "Blue1", champion: "/gnar-league-champion.jpg" },
          { name: "Blue2", champion: "/viego-league-champion.jpg" },
          { name: "Blue3", champion: "/ahri-league-champion.jpg" },
          { name: "Blue4", champion: "/ezreal-champion.jpg" },
          { name: "Blue5", champion: "/thresh-league-champion.jpg" },
        ],
        red: [
          { name: "Red1", champion: "/darius-league-champion.jpg" },
          { name: "Red2", champion: "/graves-champion.jpg" },
          { name: "Red3", champion: "/orianna-champion.jpg" },
          { name: "Red4", champion: "/kaisa-champion.jpg" },
          { name: "Red5", champion: "/nautilus-champion.jpg" },
        ],
      },
    },
    events: [
      {
        id: 14,
        type: "donation",
        title: "Donation",
        description: "SupportFan donated $50",
        timestamp: "Mid game",
        value: "$50",
      },
      {
        id: 15,
        type: "low_viewers",
        title: "Late Night Dip",
        description: "Viewers dropped off late",
        timestamp: "Late game",
        value: "3,102",
      },
    ],
  },
  {
    match: {
      id: 6,
      queueType: "Ranked Solo/Duo",
      timeAgo: "23 hours ago",
      result: "Victory",
      duration: "29m 05s",
      champion: {
        name: "Ahri",
        image: "/ahri-league-champion.jpg",
        level: 16,
      },
      kda: { kills: 12, deaths: 1, assists: 8 },
      kdaRatio: "20.00:1",
      summonerSpells: ["/flash-league-spell-icon.jpg", "/ignite-spell.jpg"],
      items: [
        "/stridebreaker-item.jpg",
        "/dead-mans-plate-item.jpg",
        "/force-of-nature-item.jpg",
        "/mortal-reminder-item.jpg",
        "/warmogs-item.jpg",
        "/control-ward-league-item.jpg",
      ],
      stats: {
        laning: "71 : 29",
        pKill: "87%",
        cs: "232 (8.0)",
        rank: "Masters",
      },
      placement: "1st",
      performanceTag: "Carry",
      teams: {
        blue: [
          { name: "Team1", champion: "/gnar-league-champion.jpg" },
          { name: "Team2", champion: "/viego-league-champion.jpg" },
          { name: "Team3", champion: "/ahri-league-champion.jpg" },
          { name: "Team4", champion: "/jinx-league-champion.jpg" },
          { name: "Team5", champion: "/nautilus-champion.jpg" },
        ],
        red: [
          { name: "Opp1", champion: "/darius-league-champion.jpg" },
          { name: "Opp2", champion: "/lee-sin-champion.jpg" },
          { name: "Opp3", champion: "/syndra-champion.jpg" },
          { name: "Opp4", champion: "/kaisa-champion.jpg" },
          { name: "Opp5", champion: "/thresh-league-champion.jpg" },
        ],
      },
    },
    events: [
      {
        id: 16,
        type: "peak_viewers",
        title: "Peak Hit!",
        description: "Session peak viewers",
        timestamp: "Team fight",
        value: "14,203",
      },
      {
        id: 17,
        type: "gifted_subs",
        title: "Gift Bomb",
        description: "ProFan gifted 100 subs",
        timestamp: "Post-win",
        value: "100",
      },
      {
        id: 18,
        type: "raid",
        title: "Outgoing Raid",
        description: "Raided FriendStreamer",
        timestamp: "End of session",
        value: "9,500",
      },
    ],
  },
]

export function SessionTimeline() {
  return (
    <div className="flex flex-col gap-1">
      <div className="grid grid-cols-[1fr_280px] gap-4 mb-2">
        <h2 className="text-base font-semibold text-foreground/90">Match History</h2>
        <h2 className="text-base font-semibold text-foreground/90">Stream Events</h2>
      </div>

      <div className="flex flex-col gap-3">
        {sessionData.map((entry) => (
          <div key={entry.match.id} className="grid grid-cols-[1fr_280px] gap-4 items-start">
            {/* Match card */}
            <div className="min-w-0">
              <MatchCard match={entry.match} />
            </div>

            {/* Stream events - aligned to match card */}
            <div className="h-full">
              <MatchStreamEvents events={entry.events} />
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}

