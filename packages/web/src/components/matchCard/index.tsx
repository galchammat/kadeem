"use client"

import { useState } from "react"
import { cn } from "@/lib/utils"
import { Film, ChevronDown, ChevronUp } from "lucide-react"
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "@/components/ui/tooltip"

interface Player {
  name: string
  champion: string
}

export interface Match {
  id: number
  queueType: string
  timeAgo: string
  result: "Victory" | "Defeat"
  duration: string
  startedAt: number
  durationSeconds: number
  champion: {
    name: string
    image: string
    level: number
  }
  kda: { kills: number; deaths: number; assists: number }
  kdaRatio: string
  summonerSpells: string[]
  items: string[]
  stats: {
    laning: string
    pKill: string
    cs: string
    rank: string
  }
  placement: string
  performanceTag: string
  teams: {
    blue: Player[]
    red: Player[]
  }
}

export function MatchCard({ match }: { match: Match }) {
  const [expanded, setExpanded] = useState(false)
  const isVictory = match.result === "Victory"

  return (
    <div className="flex flex-col">
      <div
        className={cn(
          "relative flex items-stretch gap-2 rounded-lg border bg-card p-2.5 transition-colors hover:bg-card/80 overflow-hidden",
          isVictory ? "border-l-4 border-l-[var(--victory)]" : "border-l-4 border-l-[var(--defeat)]",
        )}
      >
        {/* Left section: Queue type, time, result, duration */}
        <div className="flex w-[72px] shrink-0 flex-col justify-center gap-0.5">
          <span className="text-[11px] font-semibold text-[var(--info)] leading-tight">{match.queueType}</span>
          <span className="text-[10px] text-foreground/60">{match.timeAgo}</span>
          <span className={cn("text-sm font-bold", isVictory ? "text-[var(--victory)]" : "text-[var(--defeat)]")}>
            {match.result}
          </span>
          <span className="text-[10px] text-foreground/60">{match.duration}</span>
        </div>

        {/* Champion portrait with level */}
        <div className="relative shrink-0 self-center">
          <div className="relative h-14 w-14 overflow-hidden rounded-lg border-2 border-border">
            <img
              src={match.champion.image || "/placeholder.svg"}
              alt={match.champion.name}
              className="h-full w-full object-cover"
            />
          </div>
          <div className="absolute -bottom-1 -right-1 flex h-5 w-5 items-center justify-center rounded-full bg-secondary border border-border text-[10px] font-bold text-foreground">
            {match.champion.level}
          </div>
        </div>

        {/* Summoner spells + Items */}
        <div className="flex shrink-0 flex-col justify-center gap-1 self-center min-w-[78px]">
          <div className="flex gap-0.5">
            {match.summonerSpells.map((spell, i) => (
              <div key={i} className="relative h-5 w-5 overflow-hidden rounded border border-border">
                <img src={spell || "/placeholder.svg"} alt="Spell" className="h-full w-full object-cover" />
              </div>
            ))}
          </div>
          <div className="grid grid-cols-3 gap-0.5 w-fit">
            {match.items.map((item, i) => (
              <div key={i} className="relative h-6 w-6 overflow-hidden rounded border border-border bg-secondary">
                <img src={item || "/placeholder.svg"} alt="Item" className="h-full w-full object-cover" />
              </div>
            ))}
          </div>
        </div>

        {/* KDA section */}
        <div className="flex shrink-0 flex-col items-center justify-center gap-0.5 px-2 self-center min-w-[80px]">
          <span className="text-base font-bold text-foreground whitespace-nowrap">
            {match.kda.kills} / <span className="text-[var(--defeat)]">{match.kda.deaths}</span> / {match.kda.assists}
          </span>
          <span className="text-[10px] font-medium text-foreground/70">{match.kdaRatio} KDA</span>
        </div>

        {/* Stats section */}
        <div className="flex shrink-0 flex-col justify-center gap-0.5 border-l border-border pl-2 text-[11px] self-center min-w-[92px]">
          <div className="flex items-center gap-1.5">
            <span className="text-[var(--gold)] font-semibold">Laning</span>
            <span className="text-foreground/80 font-medium">{match.stats.laning}</span>
          </div>
          <div className="flex items-center gap-1.5">
            <span className="text-foreground/60">P/Kill</span>
            <span className="text-foreground/90 font-medium">{match.stats.pKill}</span>
          </div>
          <div className="flex items-center gap-1.5">
            <span className="text-foreground/60">CS</span>
            <span className="text-foreground/90 font-medium">{match.stats.cs}</span>
          </div>
          <div className="flex items-center gap-1.5">
            <span className="text-[var(--warning)]">⚔</span>
            <span className="text-foreground/90 font-medium">{match.stats.rank}</span>
          </div>
        </div>

        {/* Placement + Performance tag */}
        <div className="flex shrink-0 flex-col items-center justify-center gap-1 px-2 self-center min-w-[60px]">
          <span className="text-lg font-bold text-foreground">{match.placement}</span>
          <span
            className={cn(
              "rounded px-2 py-0.5 text-[10px] font-semibold",
              match.performanceTag === "Carry"
                ? "bg-[var(--victory)]/20 text-[var(--victory)]"
                : match.performanceTag === "Struggle"
                  ? "bg-[var(--defeat)]/20 text-[var(--defeat)]"
                  : "bg-secondary text-foreground/70",
            )}
          >
            {match.performanceTag}
          </span>
        </div>

        <div className="flex items-center gap-2 shrink-0 self-center">
          {/* Blue team - neutral color now */}
          <div className="flex flex-col gap-0.5">
            {match.teams.blue.map((player, i) => (
              <div key={i} className="flex items-center gap-1">
                <div className="relative h-4 w-4 overflow-hidden rounded border border-border">
                  <img
                    src={player.champion || "/placeholder.svg"}
                    alt={player.name}
                    className="h-full w-full object-cover"
                  />
                </div>
                <span className="w-12 truncate text-[10px] font-medium text-foreground/70">{player.name}</span>
              </div>
            ))}
          </div>

          {/* Red team - neutral color now */}
          <div className="flex flex-col gap-0.5">
            {match.teams.red.map((player, i) => (
              <div key={i} className="flex items-center gap-1">
                <div className="relative h-4 w-4 overflow-hidden rounded border border-border">
                  <img
                    src={player.champion || "/placeholder.svg"}
                    alt={player.name}
                    className="h-full w-full object-cover"
                  />
                </div>
                <span className="w-12 truncate text-[10px] font-medium text-foreground/70">{player.name}</span>
              </div>
            ))}
          </div>
        </div>

        <div className="flex flex-col gap-1 shrink-0 self-center ml-auto pl-2">
          <TooltipProvider delayDuration={200}>
            <Tooltip>
              <TooltipTrigger asChild>
                <button
                  className="flex h-7 w-7 items-center justify-center rounded border border-border bg-secondary hover:bg-secondary/80 transition-colors"
                  onClick={() => {
                    /* Navigate to hyper-vod editor */
                  }}
                >
                  <Film className="h-3.5 w-3.5 text-foreground/70" />
                </button>
              </TooltipTrigger>
              <TooltipContent side="left" className="text-xs">
                Edit Hyper-VOD
              </TooltipContent>
            </Tooltip>
          </TooltipProvider>

          <TooltipProvider delayDuration={200}>
            <Tooltip>
              <TooltipTrigger asChild>
                <button
                  className="flex h-7 w-7 items-center justify-center rounded border border-border bg-secondary hover:bg-secondary/80 transition-colors"
                  onClick={() => setExpanded(!expanded)}
                >
                  {expanded ? (
                    <ChevronUp className="h-3.5 w-3.5 text-foreground/70" />
                  ) : (
                    <ChevronDown className="h-3.5 w-3.5 text-foreground/70" />
                  )}
                </button>
              </TooltipTrigger>
              <TooltipContent side="left" className="text-xs">
                {expanded ? "Collapse details" : "Expand details"}
              </TooltipContent>
            </Tooltip>
          </TooltipProvider>
        </div>
      </div>

      {expanded && (
        <div className="rounded-b-lg border border-t-0 border-border bg-card/50 p-4 h-[200px]">
          <p className="text-sm text-foreground/60">Match details coming soon...</p>
        </div>
      )}
    </div>
  )
}

