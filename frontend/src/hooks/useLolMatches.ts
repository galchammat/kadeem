import { useState } from 'react';
import { models } from '@wails/go/models';
import { ListMatches } from '@wails/go/riot/RiotClient';

export interface UseLolMatchesReturn {
  matches: models.LeagueOfLegendsMatch[];
  loading: boolean;
  error: string | null;
  fetchMatches: (
    filter: models.LolMatchFilter,
    limit: number,
    offset: number
  ) => Promise<void>;
}

export function useLolMatches(account: models.LeagueOfLegendsAccount): UseLolMatchesReturn {
  const [matches, setMatches] = useState<models.LeagueOfLegendsMatch[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchMatches = async (
    filter: models.LolMatchFilter,
    limit: number,
    offset: number
  ) => {
    try {
      setLoading(true);
      const result = await ListMatches(filter, account, limit, offset);
      setMatches(result ?? []);
      setError(null);
    } catch (err) {
      setError(`Failed to load matches: ${err}`);
      setMatches([]);
    } finally {
      setLoading(false);
    }
  };

  return { matches, loading, error, fetchMatches };
}
