import { models } from "@wails/go/models";
import { SpanChart } from "@/components/charts/spanChart"

function StreamerTimeline({ channels, LeagueOfLegendsAccounts }: { channels: models.Channel[], LeagueOfLegendsAccounts: models.LeagueOfLegendsAccount[] }) {
    return (
        <div className="p-6">
            <SpanChart />
            {/* Timeline content goes here */}
        </div>
    );
}
export default StreamerTimeline;