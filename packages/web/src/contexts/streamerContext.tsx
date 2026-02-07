import React, { createContext, useState, useEffect, type ReactNode } from "react"
import type { StreamerView, Channel } from "@/types"
import * as api from "@/lib/api"

type StreamerContextType = {
  streamers: StreamerView[]
  selectedStreamer: StreamerView | null
  loading: boolean
  error: string | null
  setSelectedStreamerName: (streamerName: string) => void
  isStreamerSelected: (streamerName: string) => boolean
  refetchStreamers: () => Promise<void>
  addChannel: (channel: Pick<Channel, "streamerId" | "channelName" | "platform">) => Promise<void>
  deleteChannel: (id: string) => Promise<void>
  addStreamer: (streamerName: string) => Promise<void>
  deleteStreamer: (streamerName: string) => Promise<void>
}

export const StreamerContext = createContext<StreamerContextType | undefined>(undefined)

// LocalStorage helpers
const SELECTED_STREAMER_KEY = "kadeem:selectedStreamer"

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
  const [streamers, setStreamers] = useState<StreamerView[]>([])
  const [selectedStreamer, setSelectedStreamer] = useState<StreamerView | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const fetchStreamers = async () => {
    try {
      setLoading(true)
      setError(null)

      const data = await api.listStreamers()
      if (!data) {
        setStreamers([])
        setSelectedStreamer(null)
        return
      }
      setStreamers(data)

      // Auto-select streamer based on localStorage or default to first streamer
      const savedStreamerName = getStoredStreamer()
      let streamerToSelect: StreamerView | null = null

      if (savedStreamerName) {
        streamerToSelect = data.find((e: StreamerView) => e.name === savedStreamerName) || null
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
      setError(err instanceof Error ? err.message : "Failed to fetch streamers")
      console.error("Error fetching streamers:", err)
    } finally {
      setLoading(false)
    }
  }

  const setSelectedStreamerName = (streamerName: string) => {
    const streamer = streamers.find((e) => e.name === streamerName)
    if (streamer) {
      setSelectedStreamer(streamer)
      setStoredStreamer(streamerName)
    }
  }

  const isStreamerSelected = (streamerName: string) => {
    return selectedStreamer?.name === streamerName
  }

  const handleAddChannel = async (channel: Pick<Channel, "streamerId" | "channelName" | "platform">) => {
    await api.addChannel(channel)
  }

  const handleDeleteChannel = async (id: string) => {
    await api.deleteChannel(id)
  }

  const handleAddStreamer = async (name: string) => {
    await api.addStreamer(name)
  }

  const handleDeleteStreamer = async (name: string) => {
    await api.deleteStreamer(name)
  }

  // Fetch streamers on mount
  useEffect(() => {
    fetchStreamers()
  }, [])

  return (
    <StreamerContext.Provider
      value={{
        streamers,
        selectedStreamer,
        loading,
        error,
        setSelectedStreamerName,
        isStreamerSelected,
        refetchStreamers: fetchStreamers,
        addChannel: handleAddChannel,
        deleteChannel: handleDeleteChannel,
        addStreamer: handleAddStreamer,
        deleteStreamer: handleDeleteStreamer,
      }}
    >
      {children}
    </StreamerContext.Provider>
  )
}
