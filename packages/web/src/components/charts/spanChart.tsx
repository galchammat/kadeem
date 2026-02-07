"use client"

import {
  CartesianGrid,
  ComposedChart,
  ReferenceArea,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts"

import { ChartConfig, ChartContainer } from "@/components/ui/chart"

export const description = "Activity spans chart (horizontal spans per entity)"

// sample span data (timestamps in seconds) - ~50 entries across entity1..entity5 spanning a week
const chartData = [
  { source: "entity1", start: 1704067200, end: 1704078000 },
  { source: "entity2", start: 1704072000, end: 1704082800 },
  { source: "entity3", start: 1704075600, end: 1704090000 },
  { source: "entity4", start: 1704086400, end: 1704097200 },
  { source: "entity5", start: 1704090000, end: 1704100800 },

  { source: "entity1", start: 1704104400, end: 1704111600 },
  { source: "entity2", start: 1704108000, end: 1704122400 },
  { source: "entity3", start: 1704115200, end: 1704126000 },
  { source: "entity4", start: 1704122400, end: 1704133200 },
  { source: "entity5", start: 1704129600, end: 1704136800 },

  { source: "entity1", start: 1704133200, end: 1704144000 },
  { source: "entity2", start: 1704136800, end: 1704147600 },
  { source: "entity3", start: 1704140400, end: 1704151200 },
  { source: "entity4", start: 1704147600, end: 1704154800 },
  { source: "entity5", start: 1704151200, end: 1704162000 },

  { source: "entity1", start: 1704165600, end: 1704176400 },
  { source: "entity2", start: 1704169200, end: 1704180000 },
  { source: "entity3", start: 1704172800, end: 1704183600 },
  { source: "entity4", start: 1704180000, end: 1704190800 },
  { source: "entity5", start: 1704187200, end: 1704198000 },

  { source: "entity1", start: 1704194400, end: 1704205200 },
  { source: "entity2", start: 1704198000, end: 1704208800 },
  { source: "entity3", start: 1704201600, end: 1704212400 },
  { source: "entity4", start: 1704208800, end: 1704216000 },
  { source: "entity5", start: 1704212400, end: 1704223200 },

  { source: "entity1", start: 1704226800, end: 1704234000 },
  { source: "entity2", start: 1704230400, end: 1704244800 },
  { source: "entity3", start: 1704237600, end: 1704248400 },
  { source: "entity4", start: 1704244800, end: 1704255600 },
  { source: "entity5", start: 1704248400, end: 1704255600 },

  { source: "entity1", start: 1704259200, end: 1704266400 },
  { source: "entity2", start: 1704262800, end: 1704270000 },
  { source: "entity3", start: 1704266400, end: 1704277200 },
  { source: "entity4", start: 1704273600, end: 1704284400 },
  { source: "entity5", start: 1704277200, end: 1704284400 },

  { source: "entity1", start: 1704288000, end: 1704295200 },
  { source: "entity2", start: 1704291600, end: 1704302400 },
  { source: "entity3", start: 1704295200, end: 1704306000 },
  { source: "entity4", start: 1704302400, end: 1704313200 },
  { source: "entity5", start: 1704306000, end: 1704316800 },

  { source: "entity1", start: 1704320400, end: 1704327600 },
  { source: "entity2", start: 1704324000, end: 1704334800 },
  { source: "entity3", start: 1704327600, end: 1704338400 },
  { source: "entity4", start: 1704334800, end: 1704345600 },
  { source: "entity5", start: 1704338400, end: 1704349200 },
]

const colors: { [key: string]: string } = {
  entity1: "var(--chart-1)",
  entity2: "var(--chart-2)",
  entity3: "var(--chart-3)",
  entity4: "var(--chart-4)",
  entity5: "var(--chart-5)",
}

// compute helpers and mappings
const minStart = Math.min(...chartData.map((d) => d.start))
const maxEnd = Math.max(...chartData.map((d) => d.end))

// list entities in desired order
const entities = Array.from(new Set(chartData.map((d) => d.source)))
// map entity -> numeric index (y position)
const entityIndex = Object.fromEntries(entities.map((e, i) => [e, i]))

// convert spans to values used by ReferenceArea (x are seconds)
const spans = chartData.map((d, i) => ({
  ...d,
  id: i,
  yi: entityIndex[d.source],
  x1: d.start,
  x2: d.end,
  color: colors[d.source],
}))

const chartConfig = {
  desktop: {
    label: "Desktop",
    color: "var(--chart-1)",
  },
} satisfies ChartConfig

const fmt = (secs: number) =>
  new Date(secs * 1000).toLocaleString(undefined, {
    month: "short",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  })

export function SpanChart() {
  const yCount = entities.length
  const yPadding = 0.25 // half height of each band
  
  // Calculate dynamic height: much smaller per entity
  const LINE_HEIGHT = 25 // Very small height per entity
  const chartHeight = yCount * LINE_HEIGHT + 20 // No minimum, just calculated height
  
  console.log('Chart height calculated:', chartHeight, 'for', yCount, 'entities')

  return (

          <ChartContainer config={chartConfig} style={{ height: chartHeight, width: '100%' }}>
            <ResponsiveContainer width="100%" height="100%">
            <ComposedChart
              data={entities.map((s, i) => ({ name: s, index: i }))}
              margin={{ top: 5, right: 12, left: 12, bottom: 5 }}
            >
              <CartesianGrid vertical={false} strokeDasharray="3 3" />
              <XAxis
                type="number"
                dataKey="x" /* not used directly */
                domain={[minStart, maxEnd]}
                tickFormatter={(v) =>
                  new Date(Number(v) * 1000).toLocaleTimeString(undefined, {
                    hour: "2-digit",
                    minute: "2-digit",
                  })
                }
              />
              <YAxis
                type="number"
                domain={[-0.5, Math.max(entities.length - 0.5, 0)]}
                ticks={entities.map((_, i) => i)}
                tickFormatter={(v: number) => entities[Number(v)] ?? ""}
                width={80}
                tick={{ fontSize: 10 }}
              />
              <Tooltip
                formatter={(_: any, __: any, payload: any) => {
                  // Tooltip is left generic; custom tooltip can be implemented if needed.
                  return _
                }}
              />

              {/* Render a ReferenceArea per span; y band is centered at index */}
              {spans.map((s) => (
                <ReferenceArea
                  key={s.id}
                  x1={s.x1}
                  x2={s.x2}
                  y1={s.yi - yPadding}
                  y2={s.yi + yPadding}
                  strokeOpacity={0}
                  fill={s.color}
                  fillOpacity={0.75}
                />
              ))}
            </ComposedChart>
          </ResponsiveContainer>
        </ChartContainer>

  )
}