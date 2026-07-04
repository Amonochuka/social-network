"use client"
import { useEffect, useState } from "react";
import UserProfileImage from "@/components/header/profile/userProfile";
import { Phone, Video, MoreVertical } from "lucide-react";
import { useParams, usePathname } from "next/navigation";
import ChatContentLayout from "@/components/chats/chatContentLayout";
import { useSocketContext } from "@/contexts/SocketContext";
import { NewChat } from "@/components/chats/messages";
import { Api } from "@/services/axios";

interface ChatProfile {
  id: string;
  first_name: string;
  last_name: string;
  avatar: string;
}

interface HeaderProps {
  data: ChatProfile;
  connected?: boolean;
  isUserOnline?: boolean;
}

export default function SingleUserChatProfile({ data, connected = true, isUserOnline = false }: HeaderProps) {
  const avatarUrl = data.avatar
    ? data.avatar.startsWith("http")
      ? data.avatar
      : `http://localhost:8080/${data.avatar}`
    : undefined;

  return (
    <header className="h-20 bg-[#222222] border-b border-[#333] px-5 flex items-center justify-between">
      <div className="flex items-center gap-3">
        <UserProfileImage url={avatarUrl} name={`${data.first_name} ${data.last_name}`} />
        <div>
          <p className="font-bold text-lg">
            {data.first_name} {data.last_name}
          </p>
          <div className="flex items-center gap-2">
            <div className={`w-2 h-2 rounded-full ${
              isUserOnline ? "bg-green-500" : "bg-gray-500"
            }`} />
            <span className="text-xs text-gray-400">
              {isUserOnline ? "Online" : "Offline"}
            </span>
            {!connected && <span className="text-xs text-red-500 ml-2">Connection offline</span>}
          </div>
        </div>
      </div>
      <div className="flex gap-6">
        <Phone className="cursor-pointer" size={20} />
        <Video className="cursor-pointer" size={20} />
        <MoreVertical className="cursor-pointer" size={20} />
      </div>
    </header>
  );
}

export function ChatUser() {
  const { userId } = useParams();
  const pathname = usePathname();
  const userIdString = Array.isArray(userId) ? userId[0] : userId;
  const { connected, onlineUsers, getMessages, loadConversation, sendMessage } = useSocketContext();

  const [profile, setProfile] = useState<ChatProfile | null>(null);
  const [loading, setLoading] = useState(true);
  const [forbidden, setForbidden] = useState(false);

  const isUserOnline = userIdString ? onlineUsers.has(userIdString) : false;

  useEffect(() => {
    if (!userId) return;

    const fetchProfile = async () => {
      try {
        const res = await Api.get(`/chat/partner/${userId}`);
        setProfile(res.data);
        setForbidden(false);
      } catch (err: any) {
        console.error("Failed to load chat profile:", err);
        setProfile(null);
        setForbidden(err?.response?.status === 403);
      } finally {
        setLoading(false);
      }
    };

    fetchProfile();
  }, [userId]);

  useEffect(() => {
    if (!userIdString) return;
    loadConversation(userIdString);
  }, [userIdString, loadConversation]);

  const messages = userIdString ? getMessages(userIdString) : [];

  if (pathname === "/view/Messages") {
    return <NewChat />;
  }

  if (!userId) return null;

  if (loading) {
    return (
      <div className="flex items-center justify-center h-screen text-gray-400">
        Loading conversation...
      </div>
    );
  }

  if (!profile) {
    return (
      <div className="flex items-center justify-center h-screen text-gray-400 text-center px-6">
        {forbidden
          ? "You can't message this user. You must follow each other to chat."
          : "User not found"}
      </div>
    );
  }

  return (
    <ChatContentLayout
      data={profile}
      messages={messages}
      connected={connected}
      isUserOnline={isUserOnline}
      sendMessage={sendMessage}
    />
  );
}