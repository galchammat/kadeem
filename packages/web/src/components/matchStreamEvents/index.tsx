import {
  Flame,
  Gift,
  TrendingUp,
  TrendingDown,
  Star,
  MessageSquare,
  DollarSign,
  Zap,
  AlertTriangle,
  Trophy,
  Clapperboard,
} from "lucide-react"
import { cn } from "@/lib/utils"
import { ScrollArea } from "@/components/ui/scroll-area"
import type { EventType, StreamEvent } from "@/types"

const eventConfig: Record<EventType, { icon: typeof Flame; color: string; bg: string }> = {
  hype_train: { icon: Flame, color: "text-orange-400", bg: "bg-orange-500/15" },
  gifted_subs: { icon: Gift, color: "text-pink-400", bg: "bg-pink-500/15" },
  peak_viewers: { icon: TrendingUp, color: "text-emerald-400", bg: "bg-emerald-500/15" },
  low_viewers: { icon: TrendingDown, color: "text-amber-400", bg: "bg-amber-500/15" },
  raid: { icon: Zap, color: "text-blue-400", bg: "bg-blue-500/15" },
  donation: { icon: DollarSign, color: "text-emerald-400", bg: "bg-emerald-500/15" },
  new_subscriber: { icon: Star, color: "text-yellow-400", bg: "bg-yellow-500/15" },
  chat_milestone: { icon: MessageSquare, color: "text-blue-400", bg: "bg-blue-500/15" },
  follower_goal: { icon: Trophy, color: "text-yellow-400", bg: "bg-yellow-500/15" },
  ban_wave: { icon: AlertTriangle, color: "text-red-400", bg: "bg-red-500/15" },
  clip: { icon: Clapperboard, color: "text-purple-400", bg: "bg-purple-500/15" },
}

interface MatchStreamEventsProps {
  events: StreamEvent[]
}

export function MatchStreamEvents({ events }: MatchStreamEventsProps) {
  return (
    <div className="rounded-lg border border-border bg-card/50 overflow-hidden h-full min-h-[110px]">
      <ScrollArea className="h-full">
        <div className="flex flex-col gap-1 p-2">
          {events.map((event) => {
            const config = eventConfig[event.type]
            const Icon = config.icon

            return (
              <div
                key={event.id}
                className={cn("flex items-center gap-2 rounded-md px-2 py-1.5 transition-colors", config.bg)}
              >
                <div className={cn("shrink-0", config.color)}>
                  <Icon className="h-3.5 w-3.5" />
                </div>

                <div className="flex flex-1 items-center justify-between min-w-0 gap-2">
                  <div className="min-w-0 flex-1">
                    <span className={cn("text-[11px] font-semibold block truncate", config.color)}>{event.title}</span>
                    <span className="text-[9px] text-foreground/50 block truncate">{event.description}</span>
                  </div>
                  {event.value && (
                    <span className="text-[11px] font-bold text-foreground/90 shrink-0">{event.value}</span>
                  )}
                </div>
              </div>
            )
          })}

          {events.length === 0 && (
            <div className="text-[10px] text-foreground/40 text-center py-4">No notable events</div>
          )}
        </div>
      </ScrollArea>
    </div>
  )
}

