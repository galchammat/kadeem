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
import { PencilIcon, TrashIcon } from 'lucide-react';
import { Twitch, Youtube, PlusIcon, AlertCircleIcon } from 'lucide-react';
import { useState } from 'react';
import { models } from '@wails/go/models';
import { type DialogMode } from '@/types';
import { useStreamer } from '@/hooks/useStreamer';
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
  const { selectedStreamer, addChannel, loading } = useStreamer();
  if (selectedStreamer === null) {
    return <p className="p-6 text-muted-foreground">No streamer selected.</p>;
  }

  const defaultFormData: models.Channel = {
    id: 0,
    streamerId: selectedStreamer.id,
    platform: '',
    channelName: '',
    channelId: '',
    avatarUrl: '',
  }

  const [dialogMode, setDialogMode] = useState<DialogMode>(null);
  const [formData, setFormData] = useState<models.Channel>(defaultFormData);
  const [formError, setFormError] = useState<string | null>(null);



  return (
    <div className="p-6">
      <div className="flex items-center justify-between mb-6">
        <Sectionheader title="Livestream" icon={<TvMinimalPlay size={32} />} />
        <Dialog open={dialogMode === 'add'} onOpenChange={(open) => { if (!open) closeDialog(); }}>
          <DialogTrigger asChild>
            <Button onClick={() => openDialog('add')}>
              <PlusIcon className="h-4 w-4 mr-2" /> Add Channel
            </Button>
          </DialogTrigger>
          {/* empty DialogContent here because form is shared below in single dialog */}
        </Dialog>
      </div>
      {
        selectedStreamer.channels.length === 0 ? (
          <p className="text-muted-foreground">No accounts found. Add an account to get started.</p>
        ) : (
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
            {selectedStreamer.channels.map((channel) => (
              <div
                key={channel.channelId}
                className="border rounded-lg p-4 hover:shadow-md transition-shadow"
              >
                <div className="flex items-start justify-between">
                  <div className="space-y-1 flex-1">
                    <h3 className="font-semibold text-lg">
                      {channel.channelName}
                      <span className="text-muted-foreground">#{channel.channelId}</span>
                    </h3>
                    <p className="flex flex-row gap-2 text-sm text-muted-foreground">Platform: {channel.platform} {getPlatformIcon(channel.platform)}</p>
                    <p className="text-xs text-muted-foreground font-mono mt-2 break-all">
                      Last Streamed: {null}
                    </p>
                  </div>
                  <div className="flex gap-2 ml-2">
                    <Button variant="ghost" size="icon" onClick={() => openDialog('edit', account)} className="h-8 w-8">
                      <PencilIcon className="h-4 w-4" />
                    </Button>
                    <Button variant="ghost" size="icon" onClick={() => handleDelete(account.puuid)} className="h-8 w-8 text-destructive hover:text-destructive">
                      <TrashIcon className="h-4 w-4" />
                    </Button>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      {/* Shared dialog for add/edit */}
      <Dialog open={dialogMode !== null} onOpenChange={(open) => { if (!open) closeDialog(); }}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{dialogMode === 'add' ? 'Add New Channel' : 'Edit Channel'}</DialogTitle>
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
                <AlertTitle>{dialogMode === 'add' ? 'Failed to add account.' : 'Failed to update account.'}</AlertTitle>
                <AlertDescription>{formError}</AlertDescription>
              </Alert>
            )}
          </div>

          <DialogFooter>
            <Button variant="outline" onClick={closeDialog}>Cancel</Button>
            <Button onClick={submit} disabled={formLoading}>
              {formLoading ? (dialogMode === 'add' ? 'Adding...' : 'Updating...') : (dialogMode === 'add' ? 'Add Account' : 'Update Account')}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>

  )
}

export default Channels;