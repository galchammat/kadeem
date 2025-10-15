import { Avatar, AvatarImage } from "@/components/ui/avatar";

type SectionHeaderProps = {
  title: string;
  icon?: React.ReactNode;
}

function SectionHeader({ title, icon }: SectionHeaderProps) {
  return (
    <div className="flex items-center gap-4">
      {icon && (
        <span className="flex-shrink-0">{icon}</span>
      )}
      <h1 className="text-3xl font-bold">{title}</h1>
    </div>
  );
}

export default SectionHeader;
