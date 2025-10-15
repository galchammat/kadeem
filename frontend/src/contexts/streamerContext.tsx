import React, { createContext, useState, useEffect, type ReactNode } from 'react'
import { type Streamer, Stream, Broadcast } from '@/wailsjs/go/models'

type EntityContextType = {
  streamers: Streamer[]
  selectedEntity: Streamer | null
  loading: boolean
  error: string | null
  setSelectedStreamerName: (streamerName: string) => void
  isStreamerSelected: (streamerName: string) => boolean
  refetchStreamers: () => Promise<void>
}

export const EntityContext = createContext<EntityContextType | undefined>(undefined)

// LocalStorage helpers
const SELECTED_STREAMER_KEY = 'kadeem:selectedStreamer'

const getStoredStreamer = (): string | null => {
  try {
    return localStorage.getItem(SELECTED_STREAMER_KEY)
  } catch {
    return null
  }
}

const setStoredStreamer = (streamerName: string) => {
  try {
    localStorage.setItem(SELECTED_STREAMER_KEY, streamerName)
  } catch {
    // Ignore localStorage errors (private browsing, quota exceeded, etc.)
  }
}

export function EntityProvider({ children }: { children: ReactNode }) {
  const [streamers, setStreamers] = useState<Streamer[]>([])
  const [selectedStreamer, setSelectedStreamer] = useState<Streamer | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const fetchStreamers = async () => {
    try {
      setLoading(true)
      setError(null)

      const response = await fetch('http://localhost:8081/entities', {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        }
      })

      if (!response.ok) {
        throw new Error(`Failed to fetch entities: ${response.status}`)
      }

      const data = await response.json()
      const transformedData: Streamer[] = data.map((item: any) => ({
        Name: item.Name,
        DisplayName: item.DisplayName,
        Storage: item.Storage,
        ItemCount: item.ItemCount,
        Sources: item.Sources
      }))
      setStreamers(transformedData)

      // Auto-select streamer based on localStorage or default to first streamer
      const savedStreamerName = getStoredStreamer()
      let streamerToSelect: Streamer | null = null

      if (savedStreamerName) {
        streamerToSelect = data.find((e: Streamer) => e.Name === savedStreamerName) || null
      }

      // If no saved streamer or saved streamer not found, select first streamer
      if (!streamerToSelect && data.length > 0) {
        streamerToSelect = data[0]
      }

      if (streamerToSelect) {
        setSelectedStreamer(streamerToSelect)
        setStoredStreamer(streamerToSelect.Name)
      }

    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch streamers')
      console.error('Error fetching streamers:', err)

    } finally {
      setLoading(false)
    }
  }

  const setSelectedStreamerName = (streamerName: string) => {
    const streamer = streamers.find(e => e.Name === streamerName)
    if (streamer) {
      setSelectedStreamer(streamer)
      setStoredStreamer(streamerName)
    }
  }

  const isStreamerSelected = (streamerName: string) => {
    return selectedStreamer?.Name === streamerName
  }

  // Fetch streamers on mount
  useEffect(() => {
    fetchStreamers()
  }, [])

  return (
    <EntityContext.Provider value={{
      streamers,
      selectedStreamer,
      loading,
      error,
      setSelectedStreamerName,
      isStreamerSelected,
      refetchStreamers: fetchStreamers
    }}>
      {children}
    </EntityContext.Provider>
  )
}