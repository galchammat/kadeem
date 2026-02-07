import { Avatar, AvatarImage, AvatarFallback } from "@/components/ui/avatar";

function LeagueOfLegendsIcon() {
    return (
        <Avatar>
            <AvatarImage src="icons/league-of-legends-icon.webp" className="h-8 w-8 mr-2" />
            <AvatarFallback>LOFallback</AvatarFallback>
        </Avatar>
    )
}

export default LeagueOfLegendsIcon;