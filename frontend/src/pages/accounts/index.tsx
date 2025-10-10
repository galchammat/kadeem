import { useState, useEffect } from 'react';
import * as RiotClient from '../../../wailsjs/go/riot/RiotClient';
import { models } from '../../../wailsjs/go/models';

export function AccountsPage() {
  const [accounts, setAccounts] = useState<models.LeagueOfLegendsAccount[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchAccounts = async () => {
      try {
        setLoading(true);
        const filter = new models.LeagueOfLegendsAccount();
        const result = await RiotClient.ListAccounts(filter);
        setAccounts(result);
        setError(null);
      } catch (err) {
        setError(`Failed to load accounts: ${err}`);
      } finally {
        setLoading(false);
      }
    };

    fetchAccounts();
  }, []);

  if (loading) {
    return (
      <div className="p-6">
        <h1 className="text-3xl font-bold mb-6">Accounts</h1>
        <p className="text-muted-foreground">Loading accounts...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="p-6">
        <h1 className="text-3xl font-bold mb-6">Accounts</h1>
        <p className="text-destructive">{error}</p>
      </div>
    );
  }

  return (
    <div className="p-6">
      <h1 className="text-3xl font-bold mb-6">Accounts</h1>
      
      {accounts.length === 0 ? (
        <p className="text-muted-foreground">No accounts found. Add an account to get started.</p>
      ) : (
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
          {accounts.map((account) => (
            <div
              key={account.puuid}
              className="border rounded-lg p-4 hover:shadow-md transition-shadow"
            >
              <div className="flex items-start justify-between">
                <div className="space-y-1">
                  <h3 className="font-semibold text-lg">
                    {account.gameName}
                    <span className="text-muted-foreground">#{account.tagLine}</span>
                  </h3>
                  <p className="text-sm text-muted-foreground">Region: {account.region}</p>
                  {account.streamer && (
                    <p className="text-sm text-muted-foreground">Streamer: {account.streamer}</p>
                  )}
                  <p className="text-xs text-muted-foreground font-mono mt-2 break-all">
                    {account.puuid}
                  </p>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

export default AccountsPage;
