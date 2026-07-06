"use client";

import NavSideBar from "@/components/sidebar/sidebar";
import ChatSideBar from "@/components/chats/ChatSideBar";
import { usePathname } from "next/navigation";
import { NewChat } from "@/components/chats/messages";

export default function ChatLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const path = usePathname();

  return (
    <div className="min-h-screen w-full bg-black text-white flex">
      {/* Main navigation */}
      <div className="hidden sm:flex w-[88px] xl:w-[275px] shrink-0 border-r border-white/10">
        <div className="fixed h-screen w-[88px] xl:w-[275px]">
          <NavSideBar />
        </div>
      </div>

      {/* Chat sidebar */}
      <div className="w-[380px] shrink-0 border-r border-white/10">
        <ChatSideBar />
      </div>

      {/* Conversation */}
      <main className="flex-1">
        {path === "/view/Messages" ? <NewChat /> : children}
      </main>
    </div>
  );
}