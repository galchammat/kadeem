import { useState, useEffect } from "react"
import { MatchCard, type Match } from "@/components/matchCard"
import { MatchStreamEvents, type StreamEvent } from "@/components/matchStreamEvents"
import Sectionheader from "@/components/sectionHeader"
import { useLolAccounts } from "@/hooks/useLolAccounts"
import useLolMatches from "@/hooks/useLolMatches"
import { transformMatch } from "@/lib/matchTransformer"
import { Alert, AlertDescription } from "@/components/ui/alert"
import { Button } from "@/components/ui/button"
import { Loader2 } from "lucide-react"

interface SessionEntry {
  match: Match
  events: StreamEvent[]
}

interface SessionTimelineProps {
  streamerId?: number
}

export default function SessionTimeline({ streamerId }: SessionTimelineProps) {
  const { accounts, loading: accountsLoading, error: accountsError, fetchAccounts } = useLolAccounts()
  const { matches, loading: matchesLoading, error: matchesError, partialErrors, fetchMatches } = useLolMatches(accounts)
  const [sessionEntries, setSessionEntries] = useState<SessionEntry[]>([])
  const [transforming, setTransforming] = useState(false)
  const [transformError, setTransformError] = useState<string | null>(null)
  const [currentOffset, setCurrentOffset] = useState(0)
  const MATCHES_PER_PAGE = 10

  // Fetch accounts on mount
  useEffect(() => {
    fetchAccounts(streamerId)
  }, [streamerId])

  // Transform matches when they change
  useEffect(() => {
    const transformMatches = async () => {
      if (matches.length === 0) {
        setSessionEntries([])
        return
      }

      setTransforming(true)
      setTransformError(null)

      try {
        // Get list of tracked PUUIDs
        const trackedPUUIDs = accounts.map(acc => acc.puuid)
        
        // Transform each match
        const transformedEntries: SessionEntry[] = []
        for (const match of matches) {
          // Find which tracked PUUID is in this match
          const trackedPUUID = trackedPUUIDs.find(puuid => 
            match.participants.some(p => p.puuid === puuid)
          )

          if (!trackedPUUID) {
            console.warn(`No tracked PUUID found in match ${match.summary.gameId}`)
            continue
          }

          try {
            const transformedMatch = await transformMatch(match, trackedPUUID)
            transformedEntries.push({
              match: transformedMatch,
              events: [], // Stream events stubbed for now
            })
          } catch (err) {
            console.error(`Failed to transform match ${match.summary.gameId}:`, err)
          }
        }

        setSessionEntries(transformedEntries)
      } catch (err) {
        setTransformError(`Failed to transform matches: ${err}`)
      } finally {
        setTransforming(false)
      }
    }

    transformMatches()
  }, [matches, accounts])

  const handleLoadMore = async () => {
    const newOffset = currentOffset + MATCHES_PER_PAGE
    setCurrentOffset(newOffset)
    await fetchMatches({ offset: newOffset, limit: MATCHES_PER_PAGE })
  }

  // Loading state
  if (accountsLoading || (matchesLoading && sessionEntries.length === 0)) {
    return (
      <div className="space-y-3">
        <Sectionheader title="Session Timeline" />
        <div className="flex items-center justify-center py-8">
          <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
          <span className="ml-2 text-muted-foreground">Loading matches...</span>
        </div>
      </div>
    )
  }

  // Error states
  if (accountsError) {
    return (
      <div className="space-y-3">
        <Sectionheader title="Session Timeline" />
        <Alert variant="destructive">
          <AlertDescription>{accountsError}</AlertDescription>
        </Alert>
      </div>
    )
  }

  if (accounts.length === 0) {
    return (
      <div className="space-y-3">
        <Sectionheader title="Session Timeline" />
        <Alert>
          <AlertDescription>
            No League of Legends accounts configured for this streamer. Add accounts to see match history.
          </AlertDescription>
        </Alert>
      </div>
    )
  }

  if (matchesError) {
    return (
      <div className="space-y-3">
        <Sectionheader title="Session Timeline" />
        <Alert variant="destructive">
          <AlertDescription>{matchesError}</AlertDescription>
        </Alert>
      </div>
    )
  }

  // Partial errors (some accounts failed)
  const hasPartialErrors = partialErrors.size > 0
  const partialErrorMessage = hasPartialErrors
    ? `Failed to load matches from ${partialErrors.size} of ${accounts.length} accounts`
    : null

  return (
    <div className="space-y-3">
      <Sectionheader title="Session Timeline" />

      {/* Partial error alert */}
      {hasPartialErrors && (
        <Alert variant="destructive">
          <AlertDescription>{partialErrorMessage}</AlertDescription>
        </Alert>
      )}

      {/* Transform error */}
      {transformError && (
        <Alert variant="destructive">
          <AlertDescription>{transformError}</AlertDescription>
        </Alert>
      )}

      {/* Empty state */}
      {sessionEntries.length === 0 && !transforming && (
        <Alert>
          <AlertDescription>
            No matches found. Matches will appear here once they are synced.
          </AlertDescription>
        </Alert>
      )}

      {/* Match list */}
      {sessionEntries.length > 0 && (
        <div className="flex flex-col gap-3">
          {sessionEntries.map((entry) => (
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
      )}

      {/* Load More button */}
      {sessionEntries.length > 0 && matches.length >= currentOffset + MATCHES_PER_PAGE && (
        <div className="flex justify-center pt-4">
          <Button
            onClick={handleLoadMore}
            disabled={matchesLoading || transforming}
            variant="outline"
          >
            {matchesLoading || transforming ? (
              <>
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                Loading...
              </>
            ) : (
              "Load More"
            )}
          </Button>
        </div>
      )}

      {/* Loading indicator for pagination */}
      {(matchesLoading || transforming) && sessionEntries.length > 0 && (
        <div className="flex items-center justify-center py-4">
          <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
        </div>
      )}
    </div>
  )
}
