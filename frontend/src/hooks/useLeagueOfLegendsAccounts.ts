import { useState, useEffect } from 'react';
import * as RiotClient from '../../wailsjs/go/riot/RiotClient';
import { models } from '../../wailsjs/go/models';

export interface UseLeagueOfLegendsAccountsReturn {
  accounts: models.LeagueOfLegendsAccount[];
  loading: boolean;
  error: string | null;
  formData: models.LeagueOfLegendsAccount;
  formError: string | null;
  formLoading: boolean;
  editDialogOpen: boolean;
  addDialogOpen: boolean;
  setFormData: (data: models.LeagueOfLegendsAccount) => void;
  setEditDialogOpen: (open: boolean) => void;
  setAddDialogOpen: (open: boolean) => void;
  fetchAccounts: () => Promise<void>;
  handleEdit: (account: models.LeagueOfLegendsAccount) => void;
  handleDelete: (puuid: string) => Promise<void>;
  handleAddAccount: () => void;
  handleSubmitEdit: () => Promise<void>;
  handleSubmitAdd: () => Promise<void>;
}

export function useLeagueOfLegendsAccounts(): UseLeagueOfLegendsAccountsReturn {
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
      await RiotClient.AddAccount(formData.region, formData.gameName, formData.tagLine, 0);
      setAddDialogOpen(false);
      await fetchAccounts();
    } catch (err) {
      setFormError(String(err));
    } finally {
      setFormLoading(false);
    }
  };

  return {
    accounts,
    loading,
    error,
    formData,
    formError,
    formLoading,
    editDialogOpen,
    addDialogOpen,
    setFormData,
    setEditDialogOpen,
    setAddDialogOpen,
    fetchAccounts,
    handleEdit,
    handleDelete,
    handleAddAccount,
    handleSubmitEdit,
    handleSubmitAdd,
  };
}