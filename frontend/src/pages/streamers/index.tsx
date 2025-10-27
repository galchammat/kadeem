import { useStreamer } from "@/hooks/useStreamer";
import StreamerCard from "./streamerCard";
import { SkeletonCard } from "@/components/skeletonCard";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { useState } from "react";
import { toast } from "sonner";
import { useConfirm } from "@/components/confirmDialog";

export function StreamersPage() {
  const { streamers, loading, error, selectedStreamer, setSelectedStreamerName, addStreamer, deleteStreamer, refetchStreamers } = useStreamer();
  const [newStreamerInput, setNewStreamerInput] = useState("");
  const { confirm, ConfirmDialog } = useConfirm();

  const handleAddStreamer = async () => {
    const trimmedName = newStreamerInput.trim();
    if (trimmedName === "") return;

    const alreadyExists = streamers.some(s => s.name?.toLowerCase() === trimmedName.toLowerCase());
    if (alreadyExists) {
      alert("Streamer already exists. If you meant to add channels/accounts to this streamer, please do so from the accounts page after setting this streamer as 'active'.");
      return;
    }

    const success = await addStreamer(trimmedName);
    if (success) {
      setNewStreamerInput("");
      await refetchStreamers();
      toast.success(`Streamer "${trimmedName}" added successfully. You can now add channels/accounts to them from the accounts page after setting them as 'active'.`, { duration: 10000 });
    } else {
      alert("Failed to add streamer. They may already exist.");
    }
  };

  const handleDeleteStreamer = async (name: string): Promise<boolean> => {
    const ok = await confirm({
      title: "Delete streamer?",
      description: "Deleting a streamer will permanently remove all channels, accounts, and highlights associated with them. This action cannot be undone.",
      confirmLabel: "Delete",
      cancelLabel: "Cancel",
    });
    if (!ok) return false;

    try {
      const deleted = await deleteStreamer(name);
      if (deleted) {
        await refetchStreamers();
        if (selectedStreamer && selectedStreamer.name === name) {
          setSelectedStreamerName("");
        }
      } else {
        alert("Failed to delete streamer.");
      }
      return deleted;
    } catch (e) {
      alert(`Failed to delete streamer: ${e}`);
      return false;
    }
  };

  if (loading) return <SkeletonCard />;
  if (error) return <div>Error: {error}</div>;

  return (
    <div className="p-6">
      <div className="flex flex-row gap-4">
        <Input type="text" value={newStreamerInput} placeholder="Type in streamer name..." className="input w-full max-w-xs" onChange={(e) => setNewStreamerInput(e.target.value)} />
        <Button onClick={handleAddStreamer}>Add Streamer</Button>
      </div>

      {ConfirmDialog}

      <ul className="mt-6 flex flex-col gap-4">
        {streamers.map((streamer) => (
          <div className="max-w-md" key={streamer.id}>
            <StreamerCard
              isActive={selectedStreamer ? streamer.name === selectedStreamer.name : false}
              onSetActive={() => { setSelectedStreamerName(streamer.name) }}
              onDelete={handleDeleteStreamer}
              streamer={streamer}
            />
          </div>
        ))}
      </ul>
    </div>
  );
}

export default StreamersPage;