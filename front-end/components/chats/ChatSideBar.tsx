"use client";

import { useEffect, useState } from "react";
import UserProfileImage from "@/components/header/profile/userProfile";
import { SearchUI } from "@/components/main/header";
import { Logo } from "@/components/sidebar/sidebar";
import { Api } from "@/services/axios";
import Link from "next/link";

interface ConversationPreview {
  user_id: string;
  first_name: string;
  last_name: string;
  avatar: string;
  last_message: string;
  last_message_at: string;
}

function timeAgo(dateString: string): string {
  const date = new Date(dateString);
  const seconds = Math.floor((Date.now() - date.getTime()) / 1000);

  if (seconds < 60) return `${seconds}s`;
  const minutes = Math.floor(seconds / 60);
  if (minutes < 60) return `${minutes}m`;
  const hours = Math.floor(minutes / 60);
  if (hours < 24) return `${hours}h`;
  const days = Math.floor(hours / 24);
  return `${days}d`;
}

export default function ChatSideBar() {
  const [conversations, setConversations] = useState<ConversationPreview[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchConversations = async () => {
      try {
        const res = await Api.get<ConversationPreview[]>("/chat/conversations");
        setConversations(res.data ?? []);
      } catch (err) {
        console.error("Failed to load conversations:", err);
        setConversations([]);
      } finally {
        setLoading(false);
      }
    };

    fetchConversations();
  }, []);

  return (
    <aside className="w-[29%] min-w-90.5 py-2 px-6 h-screen flex flex-col gap-1 bg-[#222222]">
      <Logo />
      <div className="flex flex-col gap-2">
        <p className="font-bold py-1 text-2xl text-[#14afa7]">Chat Messages</p>
        <SearchUI placeholder="search friends...." />
      </div>
      <div className="mt-2.5 flex flex-col gap-3 flex-1 overflow-y-auto">
        {loading ? (
          <p className="text-sm text-gray-500 px-2 py-3">Loading conversations...</p>
        ) : conversations.length === 0 ? (
          <p className="text-sm text-gray-500 px-2 py-3">
            No conversations yet. Follow someone and start chatting.
          </p>
        ) : (
          conversations.map((convo) => (
            <Link
              href={`/view/Messages/${convo.user_id}`}
              key={convo.user_id}
              className="flex justify-between w-full rounded-b-md px-2 py-3 cursor-pointer hover:rounded-lg hover:bg-[#383737]"
            >
              <div className="flex gap-3 items-start flex-1 min-w-0">
                <UserProfileImage
                  url={
                    convo.avatar
                      ? convo.avatar.startsWith("http")
                        ? convo.avatar
                        : `http://localhost:8080/${convo.avatar}`
                      : undefined
                  }
                  name={`${convo.first_name} ${convo.last_name}`}
                />
                <div className="flex-1 min-w-0">
                  <p className="font-bold truncate min-w-0">
                    {convo.first_name} {convo.last_name}
                  </p>
                  <p className="text-sm text-gray-400 truncate">
                    {convo.last_message}
                  </p>
                </div>
              </div>
              <div>
                <p className="text-sm shrink-0 text-gray-300 font-bold">
                  {timeAgo(convo.last_message_at)}
                </p>
              </div>
            </Link>
          ))
        )}
      </div>
    </aside>
  );
}