import { useState } from 'react';

// Placeholder match data structure - will be replaced with actual API data
interface Match {
  id: string;
  champion: string;
  result: 'win' | 'loss';
  kda: { kills: number; deaths: number; assists: number };
  gameMode: string;
  gameDuration: string;
  timestamp: string;
}

// Mock data for demonstration
const mockMatches: Match[] = [
  {
    id: '1',
    champion: 'Ahri',
    result: 'win',
    kda: { kills: 12, deaths: 3, assists: 8 },
    gameMode: 'Ranked Solo',
    gameDuration: '28:45',
    timestamp: '2 hours ago',
  },
  {
    id: '2',
    champion: 'Zed',
    result: 'loss',
    kda: { kills: 5, deaths: 7, assists: 4 },
    gameMode: 'Ranked Solo',
    gameDuration: '32:15',
    timestamp: '5 hours ago',
  },
  {
    id: '3',
    champion: 'Lux',
    result: 'win',
    kda: { kills: 8, deaths: 2, assists: 15 },
    gameMode: 'Ranked Solo',
    gameDuration: '25:30',
    timestamp: '1 day ago',
  },
];

export function MatchesPage() {
  const [matches] = useState<Match[]>(mockMatches);

  const getKDA = (match: Match) => {
    const { kills, deaths, assists } = match.kda;
    const kda = deaths === 0 ? 'Perfect' : ((kills + assists) / deaths).toFixed(2);
    return kda;
  };

  return (
    <div className="p-6">
      <h1 className="text-3xl font-bold mb-6">Match History</h1>
      
      <div className="space-y-2">
        {matches.map((match) => (
          <div
            key={match.id}
            className={`border rounded-lg p-4 flex items-center justify-between hover:shadow-md transition-shadow ${
              match.result === 'win' 
                ? 'bg-blue-50 dark:bg-blue-950/20 border-blue-200 dark:border-blue-900' 
                : 'bg-red-50 dark:bg-red-950/20 border-red-200 dark:border-red-900'
            }`}
          >
            {/* Result indicator */}
            <div className="flex items-center gap-4 flex-1">
              <div className={`px-3 py-1 rounded font-semibold text-sm ${
                match.result === 'win' 
                  ? 'bg-blue-600 text-white' 
                  : 'bg-red-600 text-white'
              }`}>
                {match.result === 'win' ? 'Victory' : 'Defeat'}
              </div>

              {/* Champion name */}
              <div className="min-w-[100px]">
                <p className="font-semibold">{match.champion}</p>
                <p className="text-xs text-muted-foreground">{match.gameMode}</p>
              </div>

              {/* KDA */}
              <div className="flex items-center gap-2">
                <div className="text-center">
                  <p className="font-semibold">
                    {match.kda.kills} / {match.kda.deaths} / {match.kda.assists}
                  </p>
                  <p className="text-xs text-muted-foreground">
                    {getKDA(match)} KDA
                  </p>
                </div>
              </div>
            </div>

            {/* Game info */}
            <div className="flex items-center gap-6 text-sm text-muted-foreground">
              <div className="text-right">
                <p>{match.gameDuration}</p>
                <p className="text-xs">{match.timestamp}</p>
              </div>
            </div>
          </div>
        ))}
      </div>

      {matches.length === 0 && (
        <p className="text-muted-foreground text-center py-8">
          No matches found. Play some games to see your match history!
        </p>
      )}
    </div>
  );
}

export default MatchesPage;
