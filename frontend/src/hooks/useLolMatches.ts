import { useState, useEffect } from 'react';
import { models } from '@wails/go/models'
import { ListMatches } from '@wails/go/riot/RiotClient'

export interface UseLolMatchesReturn {
  matches: models.LeagueOfLegendsMatch[];
  fetchMatches: () => Promise<void>;
}

export function useLolMatches(): UseLolMatchesReturn {
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [matches, setMatches] = useState<models.LeagueOfLegendsMatch[]>([]);

  const fetchMatches = async () => {
    try {
      setLoading(true);
      const filter = new models.LeagueOfLegendsMatch();
      const result = await ListMatches(filter);
      // Defensive: ensure we always have an array, even if backend returns null
      setMatches(result ?? []);
      setError(null);
    } catch (err) {
      setError(`Failed to load accounts: ${err}`);
      setMatches([]); // Reset to empty array on error
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchMatches();
  }, []);

  return {
    matches,
    fetchMatches
  }
}
