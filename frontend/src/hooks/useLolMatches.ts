import { useState, useEffect } from 'react';
import { models } from '@wails/go/models';
import { ListMatches } from '@wails/go/riot/RiotClient';

type FetchMatchesProps = {
  filter?: models.LolMatchFilter;
  limit?: number;
  offset?: number;
  accounts?: models.LeagueOfLegendsAccount[];
};

export interface UseLolMatchesReturn {
  matches: models.LeagueOfLegendsMatch[];
  loading: boolean;
  error: string | null;
  partialErrors: Map<string, string>;
  fetchMatches: (props?: FetchMatchesProps) => Promise<void>;
}

export default function useLolMatches(
  initialAccounts: models.LeagueOfLegendsAccount[]
): UseLolMatchesReturn {
  const [matches, setMatches] = useState<models.LeagueOfLegendsMatch[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [partialErrors, setPartialErrors] = useState<Map<string, string>>(new Map());

  const fetchMatches = async (props?: FetchMatchesProps) => {
    const { filter = {}, limit = 20, offset = 0, accounts: overrideAccounts } = props || {};
    const fetchAccounts = overrideAccounts || initialAccounts;

    if (fetchAccounts.length === 0) {
      setMatches([]);
      return;
    }

    try {
      setLoading(true);
      const batchLimit = limit * 2; // Fetch 2x for deduplication buffer
      
      // Fetch from all accounts in parallel
      const results = await Promise.allSettled(
        fetchAccounts.map(acc => 
          ListMatches(filter, acc, batchLimit, offset)
        )
      );

      // Track which accounts failed
      const errors = new Map<string, string>();
      const allMatches: models.LeagueOfLegendsMatch[] = [];

      results.forEach((result, idx) => {
        const account = fetchAccounts[idx];
        if (result.status === 'fulfilled') {
          allMatches.push(...(result.value || []));
        } else {
          errors.set(account.puuid, `${result.reason}`);
        }
      });

      // If ALL accounts failed, set error
      if (errors.size === fetchAccounts.length) {
        setError('Failed to fetch matches from all accounts');
        setMatches([]);
        setPartialErrors(errors);
        return;
      }

      // Deduplicate by match ID
      const matchMap = new Map<number, models.LeagueOfLegendsMatch>();
      for (const match of allMatches) {
        if (!matchMap.has(match.summary.gameId)) {
          matchMap.set(match.summary.gameId, match);
        }
      }

      // Sort by timestamp descending
      const dedupedMatches = Array.from(matchMap.values()).sort((a, b) => {
        const timeA = a.summary.startedAt || 0;
        const timeB = b.summary.startedAt || 0;
        return timeB - timeA;
      });

      // Append to existing (for pagination) or replace
      if (offset > 0) {
        setMatches(prev => [...prev, ...dedupedMatches.slice(0, limit)]);
      } else {
        setMatches(dedupedMatches.slice(0, limit));
      }
      
      setPartialErrors(errors);
      setError(null);
    } catch (err) {
      setError(`Failed to load matches: ${err}`);
      setMatches([]);
    } finally {
      setLoading(false);
    }
  };

  // Auto-fetch on mount
  useEffect(() => {
    if (initialAccounts.length > 0) {
      fetchMatches({ limit: 20, offset: 0 });
    }
  }, [initialAccounts.map(a => a.puuid).join(',')]);

  return { matches, loading, error, partialErrors, fetchMatches };
}
