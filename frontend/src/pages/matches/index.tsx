import { SessionTimeline } from "@/components/timeline"

export default function Page() {
  return (
    <main className="min-h-screen bg-background p-4 md:p-6">
      <div className="mx-auto max-w-7xl">
        <header className="mb-6">
          <h1 className="text-2xl font-bold text-foreground">Timeline</h1>
          <p className="text-sm text-muted-foreground">Match history & stream analytics</p>
        </header>

        <SessionTimeline />
      </div>
    </main>
  )
}

