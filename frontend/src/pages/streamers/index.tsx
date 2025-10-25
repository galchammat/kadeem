import { useStreamer } from "@/hooks/useStreamer";
import StreamerCard from "./streamerCard";
import { SkeletonCard } from "@/components/skeletonCard";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { useState } from "react";

export function StreamersPage() {
  const { streamers, loading, error, selectedStreamer, setSelectedStreamerName, addStreamer, refetchStreamers } = useStreamer();
  const [newStreamerInput, setNewStreamerInput] = useState("");

  const handleAddStreamer = async () => {
    const trimmedName = newStreamerInput.trim();
    if (trimmedName === "") {
      return;
    }
    const alreadyExists = streamers.some(s => s.name?.toLowerCase() === trimmedName.toLowerCase());
    if (alreadyExists) {
      alert("Streamer already exists. If you meant to add channels/accounts to this streamer, please do so from the accounts page after setting this streamer as 'active'.");
      return;
    }
    const success = await addStreamer(trimmedName);
    if (success) {
      setNewStreamerInput("");
      refetchStreamers();
    } else {
      alert("Failed to add streamer. They may already exist.");
    }
  }

  if (loading) {
    return <SkeletonCard />
  }

  if (error) {
    return <div>Error: {error}</div>;
  }

  return (
    <div className="p-6">
      <div className="flex flex-row gap-4">
        <Input type="text" placeholder="Type in streamer name..." className="input w-full max-w-xs" onChange={(e) => setNewStreamerInput(e.target.value)} />
        <Button onClick={handleAddStreamer}>Add Streamer</Button>
      </div>
      <ul>
        {streamers.map((streamer) => (
          <div className="max-w-md">
            <StreamerCard
              isActive={selectedStreamer ? streamer.name === selectedStreamer.name : false}
              onSetActive={() => { setSelectedStreamerName(streamer.name) }}
              streamer={streamer}
            />
          </div>
        ))}
      </ul>
    </div>
  );
}

export default StreamersPage;