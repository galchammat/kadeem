import { Card, CardHeader, CardTitle, CardContent, CardFooter } from "@/components/ui/card"
import { Avatar, AvatarImage, AvatarFallback } from "@/components/ui/avatar"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { models } from "@wails/go/models";
import { Twitch, Youtube } from "lucide-react";

type Props = {
  streamer: models.StreamerView
  isActive?: boolean
  onSetActive: (id: string) => void
}

export default function StreamerCard({ streamer, isActive, onSetActive }: Props) {
  let avatarURL: string = "";
  streamer.streams.forEach((stream) => {
    if (stream.AvatarURL !== "") {
      avatarURL = stream.AvatarURL;
      return;
    }
  });

  return (
    <Card className="flex flex-row justify-between p-4">
      <CardHeader className="flex gap-4 p-0">
        <Avatar>
            <AvatarImage src={avatarURL} />
            <AvatarFallback>{streamer.name.charAt(0).toUpperCase()}</AvatarFallback>
        </Avatar>
        <div className="min-w-0 mt-1 flex-1">
          <CardTitle>{streamer.name}</CardTitle>
        </div>
      </CardHeader>

      <CardContent className="flex flex-row items-center justify-center p-0 gap-2">
        <div className="flex flex-wrap gap-1 justify-center">
          {streamer.streams.map((stream) => (
            <Badge key={stream.ID} variant="outline">
              {stream.Platform === "twitch" ? <Twitch  /> : <Youtube className="inline-block mr-1 h-3 w-3" />}
            </Badge>
          ))}
        </div>
      </CardContent>

      <CardFooter className="p-0">
        <Button
          size="sm"
          variant={isActive ? "outline" : "default"}
          onClick={() => onSetActive(streamer.name)}
          className="w-20"
        >
          {isActive ? "Active" : "Set active"}
        </Button>
      </CardFooter>
    </Card>
  )
}