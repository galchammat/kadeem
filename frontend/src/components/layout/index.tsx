import { Outlet, useNavigate, useLocation } from 'react-router';
import {
  SidebarProvider,
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarMenu,
  SidebarMenuItem,
  SidebarMenuButton,
  SidebarInset,
  SidebarTrigger,
} from '@/components/ui/sidebar';
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from '@/components/ui/accordion';
import { Separator } from '@/components/ui/separator';
import { HomeIcon, UsersIcon, GamepadIcon, LogOutIcon, UserIcon, ArrowLeftRight, SunIcon, MoonIcon } from 'lucide-react';
import { Toaster } from 'sonner';
import { useStreamer } from '@/hooks/useStreamer';
import { Avatar, AvatarFallback, AvatarImage } from '../ui/avatar';
import { Switch } from '../ui/switch';
import { useTheme } from '@/components/themeProvider';

function getPageTitle(pathname: string) {
  switch (pathname) {
    case '/':
      return 'Home';
    case '/accounts':
      return 'Accounts';
    case '/matches':
      return 'Matches';
    case '/streamers':
      return 'Streamers';
    default:
      return 'Page';
  }
}

export function Layout() {
  const { theme, setTheme } = useTheme()
  const { selectedStreamer } = useStreamer();
  const navigate = useNavigate();
  const location = useLocation();

  const isActive = (path: string) => location.pathname === path;

  return (
    <SidebarProvider>
      <Sidebar>
        <SidebarContent>
          <SidebarGroup>
            <div className="px-4 py-4">
              <h2 className="text-lg font-semibold">K A D E E M</h2>
            </div>
            <Separator />
            <SidebarGroupContent className="pt-4">
              <SidebarMenu>
                <SidebarMenuItem>
                  <SidebarMenuButton
                    onClick={() => navigate('/')}
                    isActive={isActive('/')}
                  >
                    <HomeIcon className="h-4 w-4" />
                    <span>Home</span>
                  </SidebarMenuButton>
                </SidebarMenuItem>
                <SidebarMenuItem>
                  <SidebarMenuButton
                    onClick={() => navigate('/accounts')}
                    isActive={isActive('/accounts')}
                  >
                    <UsersIcon className="h-4 w-4" />
                    <span>Accounts</span>
                  </SidebarMenuButton>
                </SidebarMenuItem>
                <SidebarMenuItem>
                  <SidebarMenuButton
                    onClick={() => navigate('/matches')}
                    isActive={isActive('/matches')}
                  >
                    <GamepadIcon className="h-4 w-4" />
                    <span>Matches</span>
                  </SidebarMenuButton>
                </SidebarMenuItem>
              </SidebarMenu>
            </SidebarGroupContent>
          </SidebarGroup>
        </SidebarContent>

        <SidebarFooter>
          <SidebarMenuItem>
            <SidebarMenuButton className="justify-between"
              onClick={() => navigate('/streamers')}
              isActive={isActive('/streamers')}
            >
              <div className="flex flex-row items-center gap-2">
                <Avatar>
                  <AvatarImage src={selectedStreamer?.avatarUrl} />
                  <AvatarFallback>{selectedStreamer?.name.charAt(0).toUpperCase()}</AvatarFallback>
                </Avatar>
                <span>{selectedStreamer?.name}</span>
              </div>
              <ArrowLeftRight />
            </SidebarMenuButton>
          </SidebarMenuItem>
          <Separator />
          <Accordion type="single" collapsible className="w-full px-2">
            <AccordionItem value="user-settings" className="border-0">
              <AccordionTrigger className="hover:no-underline">
                <div className="flex flex-row items-center gap-2">
                  <UserIcon />
                  <span>yog404</span>
                </div>
              </AccordionTrigger>
              <AccordionContent>
                <div className="space-y-1 pl-6">
                  <button
                    className="flex items-center gap-2 w-full text-sm py-2 px-2 hover:bg-sidebar-accent rounded-md transition-colors"
                    onClick={() => {
                      console.log('Logout clicked');
                      // Add logout functionality here
                    }}
                  >
                    <LogOutIcon className="h-4 w-4" />
                    <span>Logout</span>
                  </button>
                </div>
                <div className="flex flex-row space-y-1 pl-8 py-2 gap-1 w-full">
                  <MoonIcon className="h-4 w-4" />
                  <Switch id="dark-mode" checked={theme === 'light'} onCheckedChange={(checked: boolean) => { setTheme(checked ? 'light' : 'dark') }} />
                  <SunIcon className="h-4 w-4" />
                </div>
              </AccordionContent>
            </AccordionItem>
          </Accordion>
        </SidebarFooter>
      </Sidebar>

      <SidebarInset>
        <header className="flex h-17 items-center gap-4 border-b px-4">
          <SidebarTrigger />
          <Separator orientation="vertical" className="h-6" />
          <h1 className="text-3xl font-semibold">{getPageTitle(location.pathname)}</h1>
        </header>
        <main className="flex-1 overflow-auto">
          <Toaster theme='dark' position="top-center" richColors={true} />
          <Outlet />
        </main>
      </SidebarInset>
    </SidebarProvider>
  );
}