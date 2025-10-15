import { useContext } from 'react'
import { StreamerContext } from '@/contexts/streamerContext'

export function useStreamer() {
    const context = useContext(StreamerContext)
    if (context === undefined) {
        throw new Error('useStreamer must be used within a StreamerProvider')
    }
    return context
}