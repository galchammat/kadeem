import type { Channel, LolAccount } from "@/types"
import { SpanChart } from "@/components/charts/spanChart"

function StreamerTimeline({
  channels,
  lolAccounts,
}: {
  channels: Channel[]
  lolAccounts: LolAccount[]
}) {
  return (
    <div className="p-6">
      <SpanChart />
      {/* Timeline content goes here */}
    </div>
  )
}
export default StreamerTimeline
