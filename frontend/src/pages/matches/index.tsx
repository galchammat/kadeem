import { useState } from 'react';
import { Button } from '@/components/ui/button';
import { RefreshCwIcon, ExternalLinkIcon } from 'lucide-react';

// Enhanced match data structure - will be replaced with actual API data
interface Match {
  id: string;
  champion: string;
  championThumbnail?: string;
  result: 'win' | 'loss';
  kda: { kills: number; deaths: number; assists: number };
  killParticipation?: number;
  summoner1?: string;
  summoner2?: string;
  primaryRune?: string;
  secondaryRuneTree?: string;
  gameMode: string;
  gameDuration: string;
  timestamp: string;
  cs?: number;
  items?: string[];
  trinket?: string;
  participants?: { champion: string; summonerName: string; team: 'blue' | 'red' }[];
  vodUrl?: string;
}

// Mock data for demonstration with enhanced fields
const mockMatches: Match[] = [
  {
    id: '1',
    champion: 'Ahri',
    result: 'win',
    kda: { kills: 12, deaths: 3, assists: 8 },
    killParticipation: 67,
    summoner1: 'Flash',
    summoner2: 'Ignite',
    primaryRune: 'Electrocute',
    secondaryRuneTree: 'Precision',
    gameMode: 'Ranked Solo',
    gameDuration: '28:45',
    timestamp: '2 hours ago',
    cs: 245,
    items: ['Luden\'s', 'Sorc Shoes', 'Shadowflame', 'Rabadon\'s', 'Void Staff', 'Zhonya\'s'],
    trinket: 'Oracle Lens',
    vodUrl: 'https://example.com/vod/1',
    participants: [
      { champion: 'Ahri', summonerName: 'Player1', team: 'blue' },
      { champion: 'Lee Sin', summonerName: 'Player2', team: 'blue' },
      { champion: 'Darius', summonerName: 'Player3', team: 'blue' },
      { champion: 'Jinx', summonerName: 'Player4', team: 'blue' },
      { champion: 'Thresh', summonerName: 'Player5', team: 'blue' },
      { champion: 'Zed', summonerName: 'Enemy1', team: 'red' },
      { champion: 'Elise', summonerName: 'Enemy2', team: 'red' },
      { champion: 'Ornn', summonerName: 'Enemy3', team: 'red' },
      { champion: 'Caitlyn', summonerName: 'Enemy4', team: 'red' },
      { champion: 'Lux', summonerName: 'Enemy5', team: 'red' },
    ],
  },
  {
    id: '2',
    champion: 'Zed',
    result: 'loss',
    kda: { kills: 5, deaths: 7, assists: 4 },
    killParticipation: 45,
    summoner1: 'Flash',
    summoner2: 'Teleport',
    primaryRune: 'Conqueror',
    secondaryRuneTree: 'Domination',
    gameMode: 'Ranked Solo',
    gameDuration: '32:15',
    timestamp: '5 hours ago',
    cs: 198,
    items: ['Eclipse', 'Ionians', 'Serylda\'s', 'Black Cleaver', 'Edge of Night'],
    trinket: 'Farsight',
    vodUrl: 'https://example.com/vod/2',
    participants: [
      { champion: 'Garen', summonerName: 'Player1', team: 'blue' },
      { champion: 'Amumu', summonerName: 'Player2', team: 'blue' },
      { champion: 'Zed', summonerName: 'Player3', team: 'blue' },
      { champion: 'Ezreal', summonerName: 'Player4', team: 'blue' },
      { champion: 'Leona', summonerName: 'Player5', team: 'blue' },
      { champion: 'Fiora', summonerName: 'Enemy1', team: 'red' },
      { champion: 'Graves', summonerName: 'Enemy2', team: 'red' },
      { champion: 'Syndra', summonerName: 'Enemy3', team: 'red' },
      { champion: 'Kai\'Sa', summonerName: 'Enemy4', team: 'red' },
      { champion: 'Nautilus', summonerName: 'Enemy5', team: 'red' },
    ],
  },
  {
    id: '3',
    champion: 'Lux',
    result: 'win',
    kda: { kills: 8, deaths: 2, assists: 15 },
    killParticipation: 77,
    summoner1: 'Flash',
    summoner2: 'Barrier',
    primaryRune: 'Arcane Comet',
    secondaryRuneTree: 'Inspiration',
    gameMode: 'Ranked Solo',
    gameDuration: '25:30',
    timestamp: '1 day ago',
    cs: 210,
    items: ['Luden\'s', 'Sorc Shoes', 'Horizon Focus', 'Shadowflame', 'Mejai\'s', 'Rabadon\'s'],
    trinket: 'Farsight',
    vodUrl: 'https://example.com/vod/3',
    participants: [
      { champion: 'Sett', summonerName: 'Player1', team: 'blue' },
      { champion: 'Hecarim', summonerName: 'Player2', team: 'blue' },
      { champion: 'Lux', summonerName: 'Player3', team: 'blue' },
      { champion: 'Jhin', summonerName: 'Player4', team: 'blue' },
      { champion: 'Senna', summonerName: 'Player5', team: 'blue' },
      { champion: 'Malphite', summonerName: 'Enemy1', team: 'red' },
      { champion: 'Vi', summonerName: 'Enemy2', team: 'red' },
      { champion: 'Yasuo', summonerName: 'Enemy3', team: 'red' },
      { champion: 'Ashe', summonerName: 'Enemy4', team: 'red' },
      { champion: 'Braum', summonerName: 'Enemy5', team: 'red' },
    ],
  },
];

export function MatchesPage() {
  const [matches, setMatches] = useState<Match[]>(mockMatches);
  const [refreshing, setRefreshing] = useState(false);

  const getKDA = (match: Match) => {
    const { kills, deaths, assists } = match.kda;
    const kda = deaths === 0 ? 'Perfect' : ((kills + assists) / deaths).toFixed(2);
    return kda;
  };

  const getCSPerMin = (cs: number, duration: string) => {
    const [minutes, seconds] = duration.split(':').map(Number);
    const totalMinutes = minutes + seconds / 60;
    return (cs / totalMinutes).toFixed(1);
  };

  const handleRefresh = async () => {
    setRefreshing(true);
    // TODO: Call backend API to fetch new matches
    // For now, just simulate a delay
    await new Promise(resolve => setTimeout(resolve, 1000));
    setRefreshing(false);
    // In a real implementation, you would check if new matches were found
    // and only update if there are new matches
  };

  return (
    <div className="p-6">
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-3xl font-bold">Match History</h1>
        <Button onClick={handleRefresh} disabled={refreshing}>
          <RefreshCwIcon className={`h-4 w-4 mr-2 ${refreshing ? 'animate-spin' : ''}`} />
          {refreshing ? 'Refreshing...' : 'Refresh'}
        </Button>
      </div>
      
      <div className="space-y-3">
        {matches.map((match) => (
          <div
            key={match.id}
            className={`border rounded-lg p-4 hover:shadow-md transition-shadow ${
              match.result === 'win' 
                ? 'bg-blue-50 dark:bg-blue-950/20 border-blue-200 dark:border-blue-900' 
                : 'bg-red-50 dark:bg-red-950/20 border-red-200 dark:border-red-900'
            }`}
          >
            <div className="flex items-start gap-4">
              {/* Result indicator */}
              <div className="flex flex-col items-center gap-2">
                <div className={`px-3 py-1 rounded font-semibold text-sm ${
                  match.result === 'win' 
                    ? 'bg-blue-600 text-white' 
                    : 'bg-red-600 text-white'
                }`}>
                  {match.result === 'win' ? 'Victory' : 'Defeat'}
                </div>
                <div className="text-center text-xs text-muted-foreground">
                  <p>{match.gameDuration}</p>
                  <p>{match.timestamp}</p>
                </div>
              </div>

              {/* Champion and summoners */}
              <div className="flex flex-col items-center gap-1">
                <div className="w-12 h-12 bg-muted rounded-full flex items-center justify-center font-semibold text-xs">
                  {match.champion.slice(0, 3).toUpperCase()}
                </div>
                {match.summoner1 && match.summoner2 && (
                  <div className="flex gap-1">
                    <div className="w-5 h-5 bg-muted rounded text-[8px] flex items-center justify-center">
                      {match.summoner1.slice(0, 1)}
                    </div>
                    <div className="w-5 h-5 bg-muted rounded text-[8px] flex items-center justify-center">
                      {match.summoner2.slice(0, 1)}
                    </div>
                  </div>
                )}
              </div>

              {/* Runes */}
              {match.primaryRune && match.secondaryRuneTree && (
                <div className="flex flex-col gap-1">
                  <div className="w-6 h-6 bg-muted rounded-full flex items-center justify-center text-xs font-semibold">
                    {match.primaryRune.slice(0, 1)}
                  </div>
                  <div className="w-6 h-6 bg-muted/50 rounded-full flex items-center justify-center text-xs">
                    {match.secondaryRuneTree.slice(0, 1)}
                  </div>
                </div>
              )}

              {/* Champion name and mode */}
              <div className="min-w-[120px]">
                <p className="font-semibold text-lg">{match.champion}</p>
                <p className="text-xs text-muted-foreground">{match.gameMode}</p>
              </div>

              {/* KDA and stats */}
              <div className="flex flex-col gap-1 min-w-[140px]">
                <div className="flex items-center gap-2">
                  <p className="font-semibold">
                    {match.kda.kills} / {match.kda.deaths} / {match.kda.assists}
                  </p>
                </div>
                <p className="text-sm text-muted-foreground">
                  {getKDA(match)} KDA
                </p>
                {match.killParticipation && (
                  <p className="text-xs text-muted-foreground">
                    P/Kill {match.killParticipation}%
                  </p>
                )}
              </div>

              {/* CS */}
              {match.cs && (
                <div className="flex flex-col gap-1 min-w-[80px]">
                  <p className="text-sm font-semibold">{match.cs} CS</p>
                  <p className="text-xs text-muted-foreground">
                    {getCSPerMin(match.cs, match.gameDuration)} / min
                  </p>
                </div>
              )}

              {/* Items */}
              {match.items && (
                <div className="flex gap-1 flex-wrap max-w-[200px]">
                  {match.items.map((item, idx) => (
                    <div
                      key={idx}
                      className="w-8 h-8 bg-muted rounded border flex items-center justify-center text-[8px] text-center leading-tight p-0.5"
                      title={item}
                    >
                      {item.slice(0, 3)}
                    </div>
                  ))}
                  {match.trinket && (
                    <div
                      className="w-8 h-8 bg-muted/70 rounded-full border flex items-center justify-center text-[8px]"
                      title={match.trinket}
                    >
                      {match.trinket.slice(0, 1)}
                    </div>
                  )}
                </div>
              )}

              {/* Participants */}
              {match.participants && (
                <div className="flex gap-2 ml-auto">
                  <div className="flex flex-col gap-0.5">
                    {match.participants.filter(p => p.team === 'blue').map((p, idx) => (
                      <div key={idx} className="flex items-center gap-1">
                        <div className="w-4 h-4 bg-blue-600/20 rounded-full flex items-center justify-center text-[8px]">
                          {p.champion.slice(0, 1)}
                        </div>
                        <span className="text-[10px] text-muted-foreground truncate max-w-[60px]">
                          {p.summonerName}
                        </span>
                      </div>
                    ))}
                  </div>
                  <div className="flex flex-col gap-0.5">
                    {match.participants.filter(p => p.team === 'red').map((p, idx) => (
                      <div key={idx} className="flex items-center gap-1">
                        <div className="w-4 h-4 bg-red-600/20 rounded-full flex items-center justify-center text-[8px]">
                          {p.champion.slice(0, 1)}
                        </div>
                        <span className="text-[10px] text-muted-foreground truncate max-w-[60px]">
                          {p.summonerName}
                        </span>
                      </div>
                    ))}
                  </div>
                </div>
              )}

              {/* VOD Link */}
              {match.vodUrl && (
                <div className="ml-auto">
                  <a
                    href={match.vodUrl}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="flex items-center gap-1 text-sm text-primary hover:underline"
                  >
                    <ExternalLinkIcon className="h-4 w-4" />
                    <span>VOD</span>
                  </a>
                </div>
              )}
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
