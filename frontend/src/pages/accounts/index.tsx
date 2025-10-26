import { Separator } from "@/components/ui/separator";
import LeagueOfLegendsAccounts from "./leagueOfLegends";
import Channels from "./channel";
import { useStreamer } from "@/hooks/useStreamer";
import { AlertCircleIcon } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Link } from "react-router";

export function AccountsPage() {
  const { selectedStreamer } = useStreamer();
  if (selectedStreamer === null) {
    return (<div className="p-6">
      <div className="flex flex-col items-center justify-center p-6">
        <AlertCircleIcon className="h-12 w-12 text-muted-foreground mb-4" />
        <p className="text-center text-muted-foreground mb-4">
          Channels and Game-Accounts are managed per streamer. Please select a streamer to manage their accounts.
        </p>
        <Button asChild>
          <Link to="/streamers">
            Go to Streamers
          </Link>
        </Button>
      </div>
    </div>
    );
  }

  return (
    <div>
      <Channels />
      <Separator />
      <LeagueOfLegendsAccounts streamerId={selectedStreamer.id} />
    </div>
  );
}

export default AccountsPage;
