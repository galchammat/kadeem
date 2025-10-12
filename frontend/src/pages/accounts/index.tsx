import { useState, useEffect } from 'react';
import * as RiotClient from '../../../wailsjs/go/riot/RiotClient';
import { models } from '../../../wailsjs/go/models';
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
import { PencilIcon, TrashIcon, PlusIcon, AlertCircle, AlertCircleIcon, CheckCircle2Icon } from 'lucide-react';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { toast, Toaster } from 'sonner';

export function AccountsPage() {
  const [accounts, setAccounts] = useState<models.LeagueOfLegendsAccount[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [editDialogOpen, setEditDialogOpen] = useState(false);
  const [addDialogOpen, setAddDialogOpen] = useState(false);
  const [formData, setFormData] = useState<models.LeagueOfLegendsAccount>({
    puuid: "",
    gameName: '',
    tagLine: '',
    region: '',
  });
  const [formError, setFormError] = useState<string | null>(null);
  const [formLoading, setFormLoading] = useState(false);

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

  useEffect(() => {
    fetchAccounts();
  }, []);

  const handleEdit = (account: models.LeagueOfLegendsAccount) => {
    setFormData({
      gameName: account.gameName,
      tagLine: account.tagLine,
      region: account.region || '',
      puuid: account.puuid,
    });
    setFormError(null);
    setEditDialogOpen(true);
  };

  const handleDelete = async (puuid: string) => {
    if (!confirm('Are you sure you want to delete this account?')) {
      return;
    }
    try {
      await RiotClient.DeleteAccount(puuid);
      await fetchAccounts();
    } catch (err) {
      alert(`Failed to delete account: ${err}`);
    }
  };

  const handleAddAccount = () => {
    setFormData({
      puuid: '',
      gameName: '',
      tagLine: '',
      region: 'NA',
    });
    setFormError(null);
    setAddDialogOpen(true);
  };

  const handleSubmitEdit = async () => {
    if (!formData.gameName || !formData.tagLine || !formData.region || !formData.puuid) {
      setFormError('All fields are required');
      return;
    }

    setFormLoading(true);
    setFormError(null);

    try {
      await RiotClient.UpdateAccount(formData.region, formData.gameName, formData.tagLine, formData.puuid);
      setEditDialogOpen(false);
      await fetchAccounts();
      toast(
        <div>
          <span><CheckCircle2Icon className="h-6 w-6 mr-2 text-green-500 inline" /></span>
          Updated Account
        </div>, {
          description:
            <div>
              {formData.gameName}
              <span className="text-muted-foreground">#{formData.tagLine}</span>
              <span> {formData.region}</span>
            </div>
      });
    } catch (err) {
      setFormError(String(err));
    } finally {
      setFormLoading(false);
    }
  };

  const handleSubmitAdd = async () => {
    if (!formData.gameName || !formData.tagLine || !formData.region) {
      setFormError('All fields are required');
      return;
    }

    setFormLoading(true);
    setFormError(null);

    try {
      await RiotClient.AddAccount(formData.region, formData.gameName, formData.tagLine);
      setAddDialogOpen(false);
      await fetchAccounts();
      toast(
        <div>
          <span><CheckCircle2Icon className="h-6 w-6 mr-2 text-green-500 inline" /></span>
          Added Account
        </div>, {
          description:
            <div>
              {formData.gameName}
              <span className="text-muted-foreground">#{formData.tagLine}</span>
              <span> {formData.region}</span>
            </div>
      });
    } catch (err) {
      setFormError(String(err));
    } finally {
      setFormLoading(false);
    }
  };

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
      <div className="flex flex-col items-center justify-center h-full p-8">
        <Alert variant="destructive" className="max-w-md w-full">
          <AlertCircleIcon className="h-6 w-6 mr-2" />
          <AlertTitle>Failed to load accounts</AlertTitle>
          <AlertDescription>{error}</AlertDescription>
        </Alert>
        <Button className="mt-6" onClick={fetchAccounts}>
          Retry
        </Button>
      </div>
    );
  }

  return (
    <div className="p-6">
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-3xl font-bold">Accounts</h1>
        <Toaster />
        <Dialog open={addDialogOpen} onOpenChange={setAddDialogOpen}>
          <DialogTrigger asChild>
            <Button onClick={handleAddAccount}>
              <PlusIcon className="h-4 w-4 mr-2" />
              Add Account
            </Button>
          </DialogTrigger>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Add New Account</DialogTitle>
              <DialogDescription>
                Enter the account details. The account will be validated with Riot servers.
              </DialogDescription>
            </DialogHeader>
            <div className="grid gap-4 py-4">
              <div className="grid gap-2">
                <Label htmlFor="add-gameName">Game Name</Label>
                <Input
                  id="add-gameName"
                  value={formData.gameName}
                  onChange={(e) => setFormData({ ...formData, gameName: e.target.value })}
                  placeholder="Enter game name"
                />
              </div>
              <div className="grid gap-2">
                <Label htmlFor="add-tagLine">Tag Line</Label>
                <Input
                  id="add-tagLine"
                  value={formData.tagLine}
                  onChange={(e) => setFormData({ ...formData, tagLine: e.target.value })}
                  placeholder="Enter tag line"
                />
              </div>
              <div className="grid gap-2">
                <Label htmlFor="add-region">Region</Label>
                <Input
                  id="add-region"
                  value={formData.region}
                  onChange={(e) => setFormData({ ...formData, region: e.target.value })}
                  placeholder="e.g, NA, EUW, KR"
                />
              </div>
              {formError && (
                <Alert variant="destructive">
                  <AlertCircleIcon className="h-4 w-4 mr-2" />
                  <AlertTitle>Failed to add account.</AlertTitle>
                  <AlertDescription>{formError}</AlertDescription>
                </Alert>
              )}
            </div>
            <DialogFooter>
              <Button variant="outline" onClick={() => setAddDialogOpen(false)}>
                Cancel
              </Button>
              <Button onClick={handleSubmitAdd} disabled={formLoading}>
                {formLoading ? 'Adding...' : 'Add Account'}
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>

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
                <div className="space-y-1 flex-1">
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
                <div className="flex gap-2 ml-2">
                  <Button
                    variant="ghost"
                    size="icon"
                    onClick={() => handleEdit(account)}
                    className="h-8 w-8"
                  >
                    <PencilIcon className="h-4 w-4" />
                  </Button>
                  <Button
                    variant="ghost"
                    size="icon"
                    onClick={() => handleDelete(account.puuid)}
                    className="h-8 w-8 text-destructive hover:text-destructive"
                  >
                    <TrashIcon className="h-4 w-4" />
                  </Button>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}

      <Dialog open={editDialogOpen} onOpenChange={setEditDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Edit Account</DialogTitle>
            <DialogDescription>
              Update the account details. Changes will be validated with Riot servers.
            </DialogDescription>
          </DialogHeader>
          <div className="grid gap-4 py-4">
            <div className="grid gap-2">
              <Label htmlFor="edit-gameName">Game Name</Label>
              <Input
                id="edit-gameName"
                value={formData.gameName}
                onChange={(e) => setFormData({ ...formData, gameName: e.target.value })}
                placeholder="Enter game name"
              />
            </div>
            <div className="grid gap-2">
              <Label htmlFor="edit-tagLine">Tag Line</Label>
              <Input
                id="edit-tagLine"
                value={formData.tagLine}
                onChange={(e) => setFormData({ ...formData, tagLine: e.target.value })}
                placeholder="Enter tag line"
              />
            </div>
            <div className="grid gap-2">
              <Label htmlFor="edit-region">Region</Label>
              <Input
                id="edit-region"
                value={formData.region}
                onChange={(e) => setFormData({ ...formData, region: e.target.value })}
                placeholder="e.g, NA, EUW, KR"
              />
            </div>
            {formError && (
              <Alert variant="destructive">
                <AlertCircleIcon className="h-4 w-4 mr-2" />
                <AlertTitle>Failed to update account.</AlertTitle>
                <AlertDescription>{formError}</AlertDescription>
              </Alert>
            )}
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setEditDialogOpen(false)}>
              Cancel
            </Button>
            <Button onClick={handleSubmitEdit} disabled={formLoading}>
              {formLoading ? 'Updating...' : 'Update Account'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}

export default AccountsPage;
