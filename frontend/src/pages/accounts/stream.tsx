import { TvMinimalPlay } from 'lucide-react';
import Sectionheader from '@/components/sectionHeader';
import { Button } from '@/components/ui/button';
import { PencilIcon, TrashIcon } from 'lucide-react';
import { Twitch, Youtube } from 'lucide-react';


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


function StreamAccounts() {
    const accounts = [
        {
            platform: 'twitch',
            channelName: 'tarzaned',
            channelId: '123782776',
            lastStreamed: '3 days ago',
        }
    ]
    return (
        <div className="p-6">
            <div className="flex items-center justify-between mb-6">
                <Sectionheader title="Livestream" icon={<TvMinimalPlay size={32} />} />
            </div>
            {
                accounts.length === 0 ? (
                    <p className="text-muted-foreground">No accounts found. Add an account to get started.</p>
                ) : (
                    <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
                        {accounts.map((account) => (
                            <div
                                key={account.channelId}
                                className="border rounded-lg p-4 hover:shadow-md transition-shadow"
                            >
                                <div className="flex items-start justify-between">
                                    <div className="space-y-1 flex-1">
                                        <h3 className="font-semibold text-lg">
                                            {account.channelName}
                                            <span className="text-muted-foreground">#{account.channelId}</span>
                                        </h3>
                                        <p className="flex flex-row gap-2 text-sm text-muted-foreground">Platform: {account.platform} {getPlatformIcon(account.platform)}</p>
                                        <p className="text-xs text-muted-foreground font-mono mt-2 break-all">
                                            Last Streamed: {account.lastStreamed}
                                        </p>
                                    </div>
                                    <div className="flex gap-2 ml-2">
                                        <Button
                                            variant="ghost"
                                            size="icon"
                                            onClick={() => { }}
                                            className="h-8 w-8"
                                        >
                                            <PencilIcon className="h-4 w-4" />
                                        </Button>
                                        <Button
                                            variant="ghost"
                                            size="icon"
                                            onClick={() => { }}
                                            className="h-8 w-8 text-destructive hover:text-destructive"
                                        >
                                            <TrashIcon className="h-4 w-4" />
                                        </Button>
                                    </div>
                                </div>
                            </div>
                        ))}
                    </div>
                )
            }
        </div>

    )
}

export default StreamAccounts;