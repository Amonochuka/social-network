"use client";

import { useEffect, useState } from "react";
import UserProfileImage from "@/components/header/profile/userProfile";
import { SearchUI } from "@/components/main/header";
import { searchService, SearchUser } from "@/services/searchService";
import { Api } from "@/services/axios";
import { useSocketContext } from "@/contexts/SocketContext";
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

  const [search, setSearch] = useState("");
  const [searchResults, setSearchResults] = useState<ConversationPreview[]>([]);

  const { onlineUsers } = useSocketContext();

  const mapSearchUsers = (users: SearchUser[]): ConversationPreview[] =>
    users.map((user) => ({
      user_id: user.id,
      first_name: user.first_name,
      last_name: user.last_name,
      avatar: user.avatar,
      last_message: "",
      last_message_at: "",
    }));

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

  useEffect(() => {
    if (!search.trim()) {
      setSearchResults([]);
      return;
    }

    const timer = setTimeout(async () => {
      try {
        const data = await searchService.searchUsers(search);
        setSearchResults(mapSearchUsers(data.users));
      } catch (err) {
        console.error("Search failed:", err);
        setSearchResults([]);
      }
    }, 300);

    return () => clearTimeout(timer);
  }, [search]);

  const displayList = search.trim() ? searchResults : conversations;

  return (
    <aside className="w-full h-screen bg-[#222222] flex flex-col px-5 pt-3 pb-5">
      <h1 className="text-3xl font-bold text-[--primary-theme] mb-3">
        Messages
      </h1>

      <SearchUI
        placeholder="Search friends..."
        value={search}
        onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
          setSearch(e.target.value)
        }
      />

      <div className="mt-3 flex-1 overflow-y-auto flex flex-col gap-2">
        {loading ? (
          <p className="text-sm text-gray-500 px-2 py-3">
            Loading conversations...
          </p>
        ) : displayList.length === 0 ? (
          <p className="text-sm text-gray-500 px-2 py-3">
            {search.trim()
              ? "No users found."
              : "No conversations yet. Follow someone and start chatting."}
          </p>
        ) : (
          displayList.map((convo) => {
            const isOnline = onlineUsers.has(convo.user_id);

            return (
              <Link
                href={`/view/Messages/${convo.user_id}`}
                key={convo.user_id}
                className="flex justify-between rounded-lg px-2 py-3 hover:bg-[#383737] transition-colors"
              >
                <div className="flex gap-3 flex-1 min-w-0">
                  <div className="relative shrink-0">
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

                    {isOnline && (
                      <span className="absolute bottom-0 right-0 h-3.5 w-3.5 rounded-full bg-green-500 border-2 border-[#222222]" />
                    )}
                  </div>

                  <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-2">
                      <p className="font-bold truncate">
                        {convo.first_name} {convo.last_name}
                      </p>

                      {isOnline && (
                        <span className="text-xs text-green-500">
                          Online
                        </span>
                      )}
                    </div>

                    {convo.last_message && (
                      <p className="text-sm text-gray-400 truncate">
                        {convo.last_message}
                      </p>
                    )}
                  </div>
                </div>

                {convo.last_message_at && (
                  <p className="text-xs text-gray-400 shrink-0">
                    {timeAgo(convo.last_message_at)}
                  </p>
                )}
              </Link>
            );
          })
        )}
      </div>
    </aside>
  );
}