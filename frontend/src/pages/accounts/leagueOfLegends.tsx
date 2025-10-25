import { useState, useEffect } from 'react';
import * as RiotClient from '@wails/go/riot/RiotClient';
import { models } from '@wails/go/models';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog';
import { PencilIcon, TrashIcon, PlusIcon, AlertCircleIcon, CheckCircle2Icon } from 'lucide-react';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { toast } from 'sonner';
import SectionHeader from '@/components/sectionHeader';
import LeagueOfLegendsIcon from '@/components/icons/leagueOfLegends';
import { type DialogMode } from '@/types';
import { SkeletonCard } from '@/components/skeletonCard';

type propTypes = {
  streamerId: number;
} 

export default function LeagueOfLegendsAccounts({ streamerId }: propTypes) {
  const [accounts, setAccounts] = useState<models.LeagueOfLegendsAccount[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [dialogMode, setDialogMode] = useState<DialogMode>(null);
  const [formData, setFormData] = useState<models.LeagueOfLegendsAccount>({
    puuid: '',
    gameName: '',
    tagLine: '',
    region: '',
  });
  const [formError, setFormError] = useState<string | null>(null);
  const [formLoading, setFormLoading] = useState(false);

  const fetchAccounts = async () => {
    try {
      setLoading(true);
      const filter = new models.LeagueOfLegendsAccount(streamerId=streamerId);
      const res = await RiotClient.ListAccounts(filter);
      setAccounts(res);
      setError(null);
    } catch (err) {
      setError(String(err));
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    const timer = setTimeout(fetchAccounts, 100);
    return () => clearTimeout(timer);
  }, []);

  const openDialog = (mode: DialogMode, a?: models.LeagueOfLegendsAccount) => {
    if (mode === 'edit' && a) {
      setFormData({
        puuid: a.puuid,
        gameName: a.gameName,
        tagLine: a.tagLine,
        region: a.region || '',
      });
    } else {
      setFormData({ puuid: '', gameName: '', tagLine: '', region: 'NA' });
    }
    setFormError(null);
    setDialogMode(mode);
  };
  const closeDialog = () => setDialogMode(null);

  const handleDelete = async (puuid: string) => {
    if (!confirm('Are you sure you want to delete this account?')) return;
    try {
      await RiotClient.DeleteAccount(puuid);
      await fetchAccounts();
    } catch (err) {
      alert(`Failed to delete account: ${err}`);
    }
  };

  const submit = async () => {
    if (!formData.gameName || !formData.tagLine || !formData.region || (dialogMode === 'edit' && !formData.puuid)) {
      setFormError('All fields are required');
      return;
    }
    setFormLoading(true);
    setFormError(null);

    try {
      if (dialogMode === 'add') {
        await RiotClient.AddAccount(formData.region, formData.gameName, formData.tagLine);
        toast('Added account', { description: `${formData.gameName}#${formData.tagLine} ${formData.region}` });
      } else if (dialogMode === 'edit') {
        await RiotClient.UpdateAccount(formData.region, formData.gameName, formData.tagLine, formData.puuid);
        toast('Updated account', { description: `${formData.gameName}#${formData.tagLine} ${formData.region}` });
      }
      closeDialog();
      await fetchAccounts();
    } catch (err) {
      setFormError(String(err));
    } finally {
      setFormLoading(false);
    }
  };

  if (loading) {
    return (
      <div className="p-6">
        <SectionHeader title="League of Legends" icon={<LeagueOfLegendsIcon />} />
        <div className="flex flex-row p-6">
          <SkeletonCard />
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex flex-col items-center justify-center h-full p-8">
        <Alert variant="destructive" className="max-w-md w-full">
          <AlertCircleIcon className="h-6 w-6 mr-2" />
          <AlertTitle>Failed to load accounts</AlertTitle>
          <AlertDescription>{error}</AlertDescription>
        </Alert>
        <Button className="mt-6" onClick={fetchAccounts}>Retry</Button>
      </div>
    );
  }

  return (
    <div className="p-6">
      <div className="flex items-center justify-between mb-6">
        <SectionHeader title="League of Legends" icon={<LeagueOfLegendsIcon />} />
        <Dialog open={dialogMode === 'add'} onOpenChange={(open) => { if (!open) closeDialog(); }}>
          <DialogTrigger asChild>
            <Button onClick={() => openDialog('add')}>
              <PlusIcon className="h-4 w-4 mr-2" /> Add Account
            </Button>
          </DialogTrigger>
          {/* empty DialogContent here because form is shared below in single dialog */}
        </Dialog>
      </div>

      {accounts.length === 0 ? (
        <p className="text-muted-foreground">No accounts found. Add an account to get started.</p>
      ) : (
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
          {accounts.map((account) => (
            <div key={account.puuid} className="border rounded-lg p-4 hover:shadow-md transition-shadow max-w-md">
              <div className="flex items-start justify-between">
                <div className="space-y-1 flex-1 min-w-0">
                  <h3 className="font-semibold text-lg">
                    {account.gameName}
                    <span className="text-muted-foreground">#{account.tagLine}</span>
                  </h3>
                  <p className="text-sm text-muted-foreground">Region: {account.region}</p>
                  <p className="text-xs text-muted-foreground font-mono mt-2 truncate">{account.puuid}</p>
                </div>
                <div className="flex gap-2 ml-2">
                  <Button variant="ghost" size="icon" onClick={() => openDialog('edit', account)} className="h-8 w-8">
                    <PencilIcon className="h-4 w-4" />
                  </Button>
                  <Button variant="ghost" size="icon" onClick={() => handleDelete(account.puuid)} className="h-8 w-8 text-destructive hover:text-destructive">
                    <TrashIcon className="h-4 w-4" />
                  </Button>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}

      {/* Shared dialog for add/edit */}
      <Dialog open={dialogMode !== null} onOpenChange={(open) => { if (!open) closeDialog(); }}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{dialogMode === 'add' ? 'Add New Account' : 'Edit Account'}</DialogTitle>
            <DialogDescription>
              {dialogMode === 'add'
                ? 'Enter the account details. The account will be validated with Riot servers.'
                : 'Update the account details. Changes will be validated with Riot servers.'}
            </DialogDescription>
          </DialogHeader>

          <div className="grid gap-4 py-4">
            <div className="grid gap-2">
              <Label htmlFor="gameName">Game Name</Label>
              <Input id="gameName" value={formData.gameName} onChange={(e) => setFormData({ ...formData, gameName: e.target.value })} placeholder="Enter game name" />
            </div>
            <div className="grid gap-2">
              <Label htmlFor="tagLine">Tag Line</Label>
              <Input id="tagLine" value={formData.tagLine} onChange={(e) => setFormData({ ...formData, tagLine: e.target.value })} placeholder="Enter tag line" />
            </div>
            <div className="grid gap-2">
              <Label htmlFor="region">Region</Label>
              <Input id="region" value={formData.region} onChange={(e) => setFormData({ ...formData, region: e.target.value })} placeholder="e.g, NA, EUW, KR" />
            </div>

            {formError && (
              <Alert variant="destructive">
                <AlertCircleIcon className="h-4 w-4 mr-2" />
                <AlertTitle>{dialogMode === 'add' ? 'Failed to add account.' : 'Failed to update account.'}</AlertTitle>
                <AlertDescription>{formError}</AlertDescription>
              </Alert>
            )}
          </div>

          <DialogFooter>
            <Button variant="outline" onClick={closeDialog}>Cancel</Button>
            <Button onClick={submit} disabled={formLoading}>
              {formLoading ? (dialogMode === 'add' ? 'Adding...' : 'Updating...') : (dialogMode === 'add' ? 'Add Account' : 'Update Account')}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}