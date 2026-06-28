import { Home, Compass, Bell, MessageSquare, Users, Calendar, Settings, LogOut, User, UserPlus2Icon, SaveIcon, MoreHorizontal } from 'lucide-react';

export const sidebarNav = [
    {id: 1, icon: Home, text: "Home", href: "/view/Home"},
    {id: 2, icon: Compass, text: "Explore", href: "/view/Explore"},
    {id: 3, icon: Bell, text: "Notifications", href:"/view/Notifications"},
    {id: 4, icon: MessageSquare, text: "Messages", href: "/view/Messages"},
    {id: 10, icon: User, text: "Profile", href: "/view/profile"},
];

export const sidebarMoreNav = [
    {id: 5, icon: Users, text: "Groups", href: "/view/Groups"},
    {id: 6, icon: UserPlus2Icon, text: "Network", href: "/view/Network"},
    {id: 7, icon: SaveIcon, text: "Saved", href: "/view/Saved"},
    {id: 8, icon: Calendar, text: "Events", href: "/view/Events"},
    {id: 9, icon: Settings, text: "Settings", href: "/view/Settings"},
];

export const sidebarMoreTrigger = {id: 12, icon: MoreHorizontal, text: "More"};

export const sidebarLogout = {id: 11, icon: LogOut, text: "Logout", href: "#"};