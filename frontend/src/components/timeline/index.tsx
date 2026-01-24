import { MatchCard, type Match } from "@/components/matchCard"
import { MatchStreamEvents, type StreamEvent } from "@/components/matchStreamEvents"
import useLolMatches from '@/hooks/useLolMatches';

import Sectionheader from '@/components/sectionHeader';
interface SessionEntry {
  match: Match
  events: StreamEvent[]
}

export function SessionTimeline() {
  const { matches, loading, error, fetchMatches } = useLolMatches(account);

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

