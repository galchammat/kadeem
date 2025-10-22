import { useStreamer } from "@/hooks/useStreamer";
import StreamerCard from "./streamerCard";

export function StreamersPage() {
  const { streamers, loading, error, selectedStreamer, setSelectedStreamerName } = useStreamer();

  if (loading) {
    return <div>Loading...</div>;
  }

  if (error) {
    return <div>Error: {error}</div>;
  }

  return (
    <div className="p-6">
      <ul>
        {streamers.map((streamer) => (
          <div className="max-w-md">
            <StreamerCard 
              isActive={selectedStreamer ? streamer.name === selectedStreamer.name : false} 
              onSetActive={() => {setSelectedStreamerName(streamer.name)}} 
              streamer={streamer} 
            />
          </div>
        ))}
      </ul>
    </div>
  );
}

export default StreamersPage;