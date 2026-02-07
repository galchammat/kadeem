import type { Channel, LeagueOfLegendsAccount } from "@/types"
import { SpanChart } from "@/components/charts/spanChart"

function StreamerTimeline({
  channels,
  LeagueOfLegendsAccounts,
}: {
  channels: Channel[]
  LeagueOfLegendsAccounts: LeagueOfLegendsAccount[]
}) {
  return (
    <div className="p-6">
      <SpanChart />
      {/* Timeline content goes here */}
    </div>
  )
}
export default StreamerTimeline
