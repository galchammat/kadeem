import { TvMinimalPlay } from 'lucide-react';
import Sectionheader from '@/components/sectionHeader';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog';
import { Spinner } from '@/components/ui/spinner';
import { PencilIcon, TrashIcon } from 'lucide-react';
import { Twitch, Youtube, PlusIcon, AlertCircleIcon } from 'lucide-react';
import { useState } from 'react';
import { models } from '@wails/go/models';
import { type DialogMode } from '@/types';
import { useStreamer } from '@/hooks/useStreamer';
import { toast } from 'sonner';
import { useConfirm } from '@/components/confirmDialog';
// import ChannelForm from './channelForm';

function getPlatformIcon(platform: string) {
  switch (platform) {
    case 'twitch':
      return <Twitch className="h-5 w-5 text-purple-600" />;
    case 'youtube':
      return <Youtube className="h-5 w-5 text-red-600" />;
    default:
      return null;
  }
}

function Channels() {
  const { selectedStreamer, addChannel, deleteChannel, refetchStreamers, loading, error } = useStreamer();
  if (selectedStreamer === null) {
    return <p className="p-6 text-muted-foreground">No streamer selected.</p>;
  }

  const defaultFormData: Partial<models.Channel> = {
    streamerId: selectedStreamer.id,
    platform: 'twitch',
    channelName: ''
  }

  const [dialogMode, setDialogMode] = useState<DialogMode>(null);
  const [formData, setFormData] = useState<Partial<models.Channel>>(defaultFormData);
  const [formError, setFormError] = useState<string | null>(null);
  const [formLoading, setFormLoading] = useState<boolean>(false);

  const { confirm, ConfirmDialog } = useConfirm();

  const openDialog = (mode: DialogMode, a?: models.Channel) => {
    if (mode === 'edit' && a) {
      setFormData(defaultFormData);
    }
    setFormError(null);
    setDialogMode(mode);
  };
  const closeDialog = () => setDialogMode(null);

  const handleDelete = async (channelId: string) => {
    const ok = await confirm({
      title: "Delete channel?",
      description: "Deleting a channel will permanently delete all associated broadcast data including vods managed by this app."
    });
    if (!ok) return;
    try {
      await deleteChannel(channelId);
      await refetchStreamers();
    } catch (err) {
      alert(`Failed to delete channel: ${err}`);
    }
  };

  const submit = async () => {
    if (!formData.channelName || !formData.platform) {
      setFormError('All fields are required');
      return;
    }
    setFormLoading(true);
    setFormError(null);

    try {
      if (dialogMode === 'add') {
        await addChannel(formData as models.Channel);
        toast('Added channel', { description: `${formData.channelName} ${formData.platform}` });
      }
      closeDialog();
      await refetchStreamers();
    } catch (err) {
      setFormError(String(err));
    } finally {
      setFormLoading(false);
    }
  };


  if (error) {
    return (
      <div className="flex flex-col items-center justify-center h-full p-8">
        <Alert variant="destructive" className="max-w-md w-full">
          <AlertCircleIcon className="h-6 w-6 mr-2" />
          <AlertTitle>Failed to load channels</AlertTitle>
          <AlertDescription>{error}</AlertDescription>
        </Alert>
        <Button className="mt-6" onClick={refetchStreamers}>Retry</Button>
      </div>
    );
  } else if (loading) {
    return (
      <div className="flex items-center gap-6">
        <Spinner className="size-6 text-red-500" />
        <span>Loading Streamers ...</span>
        <Spinner className="size-6 text-purple-500" />
      </div>
    )
  }

  return (
    <div className="p-6">
      <div className="flex items-center justify-between mb-6">
        <Sectionheader title="Livestream Channels" icon={<TvMinimalPlay size={32} />} />
        <Dialog open={dialogMode === 'add'} onOpenChange={(open) => { if (!open) closeDialog(); }}>
          <DialogTrigger asChild>
            <Button onClick={() => openDialog('add')}>
              <PlusIcon className="h-4 w-4 mr-2" /> Add Channel
            </Button>
          </DialogTrigger>
          {/* empty DialogContent here because form is shared below in single dialog */}
        </Dialog>
      </div>
      {ConfirmDialog}
      {
        !selectedStreamer.channels || selectedStreamer.channels.length === 0 ? (
          <p className="text-muted-foreground">{`No channels found for ${selectedStreamer.name}. Add a channel to get started.`}</p>
        ) : (
          <div className="grid gap-4 md:grid-cols-3 lg:grid-cols-4">
            {selectedStreamer.channels.map((channel) => (
              <div
                key={channel.id}
                className="border rounded-lg p-4 hover:shadow-md transition-shadow max-w-md"
              >
                <div className="flex items-start justify-between">
                  <div className="space-y-1 flex-1">
                    <h3 className="font-semibold text-lg">
                      {channel.channelName}
                      <span className="text-muted-foreground">#{channel.id}</span>
                    </h3>
                    <p className="flex flex-row gap-2 text-sm text-muted-foreground">Platform: {channel.platform} {getPlatformIcon(channel.platform)}</p>
                    <p className="text-xs text-muted-foreground font-mono mt-2 break-all">
                      Last Streamed: {null}
                    </p>
                  </div>
                  <div className="flex gap-2 ml-2">
                    <Button variant="ghost" size="icon" onClick={() => handleDelete(channel.id)} className="h-8 w-8 text-destructive hover:text-destructive">
                      <TrashIcon className="h-4 w-4" />
                    </Button>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      {/* Dialog for add */}
      <Dialog open={dialogMode !== null} onOpenChange={(open) => { if (!open) closeDialog(); }}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{dialogMode === 'add' ? 'Add New Channel' : ''}</DialogTitle>
            <DialogDescription>
              {dialogMode === 'add'
                ? 'Enter the channel details. The channel will be validated with the selected platform.'
                : 'Update the channel details. Changes will be validated with the selected platform.'}
            </DialogDescription>
          </DialogHeader>

          <div className="grid gap-4 py-4">
            <div className="grid gap-2">
              <Label htmlFor="channelName">Name</Label>
              <Input id="channelName" value={formData.channelName} onChange={(e) => setFormData({ ...formData, channelName: e.target.value })} placeholder="Enter channel name" />
            </div>
            <div className="grid gap-2">
              <Label htmlFor="region">Platform</Label>
              {formData.platform}
            </div>

            {formError && (
              <Alert variant="destructive">
                <AlertCircleIcon className="h-4 w-4 mr-2" />
                <AlertTitle>{dialogMode === 'add' ? 'Failed to add channel.' : 'Failed to update channel.'}</AlertTitle>
                <AlertDescription>{formError}</AlertDescription>
              </Alert>
            )}
          </div>

          <DialogFooter>
            <Button variant="outline" onClick={closeDialog}>Cancel</Button>
            <Button onClick={submit} disabled={formLoading}>
              {formLoading ? (dialogMode === 'add' ? 'Adding...' : 'Updating...') : (dialogMode === 'add' ? 'Add Channel' : 'Update Channel')}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>

  )
}

export default Channels;