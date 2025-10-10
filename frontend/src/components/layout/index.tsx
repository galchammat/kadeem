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
import { HomeIcon, UsersIcon, GamepadIcon, LogOutIcon, SettingsIcon, UserIcon } from 'lucide-react';

export function Layout() {
  const navigate = useNavigate();
  const location = useLocation();

  const isActive = (path: string) => location.pathname === path;

  return (
    <SidebarProvider>
      <Sidebar>
        <SidebarContent>
          <SidebarGroup>
            <div className="px-4 py-4">
              <h2 className="text-lg font-semibold">Kadeem</h2>
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
                <SidebarMenuItem>
                  <SidebarMenuButton
                    onClick={() => navigate('/about')}
                    isActive={isActive('/about')}
                  >
                    <SettingsIcon className="h-4 w-4" />
                    <span>About</span>
                  </SidebarMenuButton>
                </SidebarMenuItem>
              </SidebarMenu>
            </SidebarGroupContent>
          </SidebarGroup>
        </SidebarContent>
        
        <SidebarFooter>
          <Accordion type="single" collapsible className="w-full px-2">
            <AccordionItem value="user-settings" className="border-0">
              <AccordionTrigger className="hover:no-underline">
                <div className="flex items-center gap-2">
                  <UserIcon className="h-4 w-4" />
                  <span className="text-sm">User Profile</span>
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
              </AccordionContent>
            </AccordionItem>
          </Accordion>
        </SidebarFooter>
      </Sidebar>
      
      <SidebarInset>
        <header className="flex h-14 items-center gap-4 border-b px-4">
          <SidebarTrigger />
          <Separator orientation="vertical" className="h-6" />
          <h1 className="text-lg font-semibold">League of Legends Tracker</h1>
        </header>
        <main className="flex-1 overflow-auto">
          <Outlet />
        </main>
      </SidebarInset>
    </SidebarProvider>
  );
}