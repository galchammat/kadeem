import { useStreamer } from "@/hooks/useStreamer";
import StreamerCard from "./streamerCard";
import { SkeletonCard } from "@/components/skeletonCard";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { useRef, useState } from "react";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog"
import { toast } from "sonner";

export function StreamersPage() {
  const { streamers, loading, error, selectedStreamer, setSelectedStreamerName, addStreamer, deleteStreamer, refetchStreamers } = useStreamer();
  const [newStreamerInput, setNewStreamerInput] = useState("");
  // State for AlertDialog-driven delete confirmation
  const [alertOpen, setAlertOpen] = useState(false);
  const [pendingDeleteName, setPendingDeleteName] = useState<string | null>(null);
  // Promise resolver for the in-flight delete request
  const deleteResolveRef = useRef<((v: boolean) => void) | null>(null);

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
      await refetchStreamers();
      toast.success(`Streamer "${trimmedName}" added successfully. You can now add channels/accounts to them from the accounts page after setting them as 'active'.`, { duration: 10000 });
    } else {
      alert("Failed to add streamer. They may already exist.");
    }
  }

  const handleDeleteStreamer = async (name: string): Promise<boolean> => {
    return new Promise<boolean>((resolve) => {
      setPendingDeleteName(name);
      deleteResolveRef.current = resolve;
      setAlertOpen(true);
    });
  }

  const onConfirmDelete = async () => {
    if (!pendingDeleteName) {
      if (deleteResolveRef.current) deleteResolveRef.current(false);
      cleanupDeleteState();
      return;
    }
    try {
      const deleted = await deleteStreamer(pendingDeleteName);
      if (deleted) {
        // wait for refetch to complete so UI updates before resolving the caller
        await refetchStreamers();

        // if the deleted streamer was selected, clear the selection
        if (selectedStreamer && selectedStreamer.name === pendingDeleteName) {
          setSelectedStreamerName('');
        }
      } else {
        alert("Failed to delete streamer.");
      }
      if (deleteResolveRef.current) deleteResolveRef.current(deleted);
    } catch (e) {
      if (deleteResolveRef.current) deleteResolveRef.current(false);
      alert(`Failed to delete streamer: ${e}`);
    } finally {
      cleanupDeleteState();
    }
  }
  // ...existing code...

  const onCancelDelete = () => {
    if (deleteResolveRef.current) deleteResolveRef.current(false);
    cleanupDeleteState();
  }

  const cleanupDeleteState = () => {
    setPendingDeleteName(null);
    setAlertOpen(false);
    deleteResolveRef.current = null;
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
        <Input type="text" value={newStreamerInput} placeholder="Type in streamer name..." className="input w-full max-w-xs" onChange={(e) => setNewStreamerInput(e.target.value)} />
        <Button onClick={handleAddStreamer}>Add Streamer</Button>
      </div>
      <AlertDialog open={alertOpen} onOpenChange={(open) => { if (!open) onCancelDelete(); }}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Delete streamer?</AlertDialogTitle>
            <AlertDialogDescription>
              Deleting a streamer will permanently remove all channels, accounts, and highlights associated with them. This action cannot be undone.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel onClick={onCancelDelete}>Cancel</AlertDialogCancel>
            <AlertDialogAction onClick={onConfirmDelete}>Delete</AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
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