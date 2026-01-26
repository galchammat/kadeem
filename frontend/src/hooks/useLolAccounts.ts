import { useState } from 'react';
import * as RiotClient from '../../wailsjs/go/riot/RiotClient';
import { models } from '../../wailsjs/go/models';

export interface UseLolAccountsReturn {
  accounts: models.LeagueOfLegendsAccount[];
  loading: boolean;
  error: string | null;
  fetchAccounts: (streamerId?: number) => Promise<void>;
  addAccount: (region: string, gameName: string, tagLine: string, streamerId: number) => Promise<void>;
  updateAccount: (region: string, gameName: string, tagLine: string, puuid: string) => Promise<void>;
  deleteAccount: (puuid: string) => Promise<void>;
}

export function useLolAccounts(): UseLolAccountsReturn {
  const [accounts, setAccounts] = useState<models.LeagueOfLegendsAccount[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchAccounts = async (streamerId?: number) => {
    try {
      setLoading(true);
      setError(null);
      const filter = streamerId
        ? { streamerId, puuid: '', gameName: '', tagLine: '', region: '' }
        : new models.LeagueOfLegendsAccount();
      const result = await RiotClient.ListAccounts(filter);
      setAccounts(result ?? []);
    } catch (err) {
      setError(`Failed to load accounts: ${err}`);
      setAccounts([]);
    } finally {
      setLoading(false);
    }
  };

  const addAccount = async (region: string, gameName: string, tagLine: string, streamerId: number) => {
    await RiotClient.AddAccount(region, gameName, tagLine, streamerId);
  };

  const updateAccount = async (region: string, gameName: string, tagLine: string, puuid: string) => {
    await RiotClient.UpdateAccount(region, gameName, tagLine, puuid);
  };

  const deleteAccount = async (puuid: string) => {
    await RiotClient.DeleteAccount(puuid);
  };

  return {
    accounts,
    loading,
    error,
    fetchAccounts,
    addAccount,
    updateAccount,
    deleteAccount,
  };
}
