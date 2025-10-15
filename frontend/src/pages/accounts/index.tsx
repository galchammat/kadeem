import { Separator } from "@/components/ui/separator";
import LeagueOfLegendsAccounts from "./leagueOfLegends";
import StreamAccounts from "./stream";

export function AccountsPage() {
    return (
        <div>
            <StreamAccounts />
            <Separator />
            <LeagueOfLegendsAccounts />
        </div>
    );
}

export default AccountsPage;
