import React, { createContext, useState, useEffect, type ReactNode } from 'react'
import { models } from '@wails/go/models' 
import { ListStreamersWithDetails, AddChannel } from '@wails/go/livestream/StreamerClient'

type StreamerContextType = {
  streamers: models.StreamerView[]
  selectedStreamer: models.StreamerView | null
  loading: boolean
  error: string | null
  setSelectedStreamerName: (streamerName: string) => void
  isStreamerSelected: (streamerName: string) => boolean
  refetchStreamers: () => Promise<void>
  addChannel: (channel: models.Channel) => Promise<boolean>
}

export const StreamerContext = createContext<StreamerContextType | undefined>(undefined)

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

export function StreamerProvider({ children }: { children: ReactNode }) {
  const [streamers, setStreamers] = useState<models.StreamerView[]>([])
  const [selectedStreamer, setSelectedStreamer] = useState<models.StreamerView | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const fetchStreamers = async () => {
    try {
      setLoading(true)
      setError(null)

      const data = await ListStreamersWithDetails()
      setStreamers(data)

      // Auto-select streamer based on localStorage or default to first streamer
      const savedStreamerName = getStoredStreamer()
      let streamerToSelect: models.StreamerView | null = null

      if (savedStreamerName) {
        streamerToSelect = data.find((e: models.StreamerView) => e.name === savedStreamerName) || null
      }

      // If no saved streamer or saved streamer not found, select first streamer
      if (!streamerToSelect && data.length > 0) {
        streamerToSelect = data[0]
      }

      if (streamerToSelect) {
        setSelectedStreamer(streamerToSelect)
        setStoredStreamer(streamerToSelect.name)
      }

    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch streamers')
      console.error('Error fetching streamers:', err)

    } finally {
      setLoading(false)
    }
  }

  const setSelectedStreamerName = (streamerName: string) => {
    const streamer = streamers.find(e => e.name === streamerName)
    if (streamer) {
      setSelectedStreamer(streamer)
      setStoredStreamer(streamerName)
    }
  }

  const isStreamerSelected = (streamerName: string) => {
    return selectedStreamer?.name === streamerName
  }

  // Fetch streamers on mount
  useEffect(() => {
    fetchStreamers()
  }, [])

  return (
    <StreamerContext.Provider value={{
      streamers,
      selectedStreamer,
      loading,
      error,
      setSelectedStreamerName,
      isStreamerSelected,
      refetchStreamers: fetchStreamers,
      addChannel: AddChannel,
    }}>
      {children}
    </StreamerContext.Provider>
  )
}