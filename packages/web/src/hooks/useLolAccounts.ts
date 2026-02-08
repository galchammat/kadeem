import { useState } from "react"
import type { LeagueOfLegendsAccount } from "@/types"
import * as api from "@/lib/api"

export interface UseLolAccountsReturn {
  accounts: LeagueOfLegendsAccount[]
  loading: boolean
  error: string | null
  fetchAccounts: () => Promise<void>
  addAccount: (region: string, gameName: string, tagLine: string, streamerId: number) => Promise<void>
  updateAccount: (accountId: string, region: string, gameName: string, tagLine: string) => Promise<void>
  deleteAccount: (accountId: string) => Promise<void>
}

export function useLolAccounts(): UseLolAccountsReturn {
  const [accounts, setAccounts] = useState<LeagueOfLegendsAccount[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const fetchAccounts = async () => {
    try {
      setLoading(true)
      setError(null)
      const result = await api.listAccounts()
      setAccounts(result ?? [])
    } catch (err) {
      setError(`Failed to load accounts: ${err}`)
      setAccounts([])
    } finally {
      setLoading(false)
    }
  }

  const addAccount = async (region: string, gameName: string, tagLine: string, streamerId: number) => {
    await api.addAccount(region, gameName, tagLine, streamerId)
  }

  const updateAccount = async (accountId: string, region: string, gameName: string, tagLine: string) => {
    await api.updateAccount(accountId, region, gameName, tagLine)
  }

  const deleteAccount = async (accountId: string) => {
    await api.deleteAccount(accountId)
  }

  return {
    accounts,
    loading,
    error,
    fetchAccounts,
    addAccount,
    updateAccount,
    deleteAccount,
  }
}
