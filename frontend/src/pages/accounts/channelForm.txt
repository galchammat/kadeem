import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { PlusIcon, AlertCircleIcon } from 'lucide-react';
import { models } from '@wails/go/models';

type ChannelFormProps = {
    addDialogOpen: boolean;
    setAddDialogOpen: (open: boolean) => void;
    formData: models.Channel;
    setFormData: (data: models.Channel) => void;
    formLoading: boolean;
    formError: string | null;
    handleSubmitAdd: () => void;
    handleAddAccount: () => void;
    handleEditAccount: () => void;
    handleSubmitEdit: () => void;
};

function ChannelForm({
    addDialogOpen,
    setAddDialogOpen,
    formData,
    setFormData,
    formLoading,
    formError,
    handleSubmitAdd,
    handleAddAccount,
    handleEditAccount,
    handleSubmitEdit
}: ChannelFormProps) {
    return (
        <Dialog open={addDialogOpen} onOpenChange={setAddDialogOpen}>
            <DialogTrigger asChild>
                <Button onClick={handleAddAccount}>
                    <PlusIcon className="h-4 w-4 mr-2" />
                    Add Account
                </Button>
            </DialogTrigger>
            <DialogContent>
                <DialogHeader>
                    <DialogTitle>Add New Account</DialogTitle>
                    <DialogDescription>
                        Enter the account details. The account will be validated with Riot servers.
                    </DialogDescription>
                </DialogHeader>
                <div className="grid gap-4 py-4">
                    <div className="grid gap-2">
                        <Label htmlFor="add-gameName">Game Name</Label>
                        <Input
                            id="add-gameName"
                            value={formData.gameName}
                            onChange={(e) => setFormData({ ...formData, gameName: e.target.value })}
                            placeholder="Enter game name"
                        />
                    </div>
                    <div className="grid gap-2">
                        <Label htmlFor="add-tagLine">Tag Line</Label>
                        <Input
                            id="add-tagLine"
                            value={formData.tagLine}
                            onChange={(e) => setFormData({ ...formData, tagLine: e.target.value })}
                            placeholder="Enter tag line"
                        />
                    </div>
                    <div className="grid gap-2">
                        <Label htmlFor="add-region">Region</Label>
                        <Input
                            id="add-region"
                            value={formData.region}
                            onChange={(e) => setFormData({ ...formData, region: e.target.value })}
                            placeholder="e.g, NA, EUW, KR"
                        />
                    </div>
                    {formError && (
                        <Alert variant="destructive">
                            <AlertCircleIcon className="h-4 w-4 mr-2" />
                            <AlertTitle>Failed to add account.</AlertTitle>
                            <AlertDescription>{formError}</AlertDescription>
                        </Alert>
                    )}
                </div>
                <DialogFooter>
                    <Button variant="outline" onClick={() => setAddDialogOpen(false)}>
                        Cancel
                    </Button>
                    <Button onClick={handleSubmitAdd} disabled={formLoading}>
                        {formLoading ? 'Adding...' : 'Add Account'}
                    </Button>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    )
}

export default ChannelForm;



