import { useState } from 'react';
import { models } from '@wails/go/models';
import { ListMatches } from '@wails/go/riot/RiotClient';

type fetchMatchesProps = {
  filter: models.LolMatchFilter,
  limit: number,
  offset: number
}

export interface UseLolMatchesReturn {
  matches: models.LeagueOfLegendsMatch[];
  loading: boolean;
  error: string | null;
  fetchMatches: (props: fetchMatchesProps) => Promise<void>;
}

export default function useLolMatches(account: models.LeagueOfLegendsAccount): UseLolMatchesReturn {
  const [matches, setMatches] = useState<models.LeagueOfLegendsMatch[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchMatches = async (props: fetchMatchesProps) => {
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
