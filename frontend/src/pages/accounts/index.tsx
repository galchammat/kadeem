import { Separator } from "@/components/ui/separator";
import LeagueOfLegendsAccounts from "./leagueOfLegends";
import Channels from "./channel";

export function AccountsPage() {
    return (
        <div>
            <Channels />
            <Separator />
            <LeagueOfLegendsAccounts />
        </div>
    );
}

export default AccountsPage;
