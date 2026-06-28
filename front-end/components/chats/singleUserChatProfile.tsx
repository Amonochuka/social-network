"use client"
import { useEffect, useState } from "react";
import UserProfileImage from "@/components/header/profile/userProfile";
import { Phone, Video, MoreVertical } from "lucide-react";
import { useParams, usePathname } from "next/navigation";
import ChatContentLayout from "@/components/chats/chatContentLayout";
import useSocket, { SocketType } from "@/hooks/useSocket";
import loadEnvFile from "@/config/config";
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
}

export default function SingleUserChatProfile({ data }: HeaderProps) {
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
  const config = loadEnvFile({ urlType: "socket-url", urlUsage: "chat" });

  const [profile, setProfile] = useState<ChatProfile | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!userId) return;

    const fetchProfile = async () => {
      try {
        const res = await Api.get(`/profile/${userId}`);
        setProfile(res.data);
      } catch (err) {
        console.error("Failed to load chat profile:", err);
        setProfile(null);
      } finally {
        setLoading(false);
      }
    };

    fetchProfile();
  }, [userId]);

  const { connected, messages, sendMessage }: SocketType = useSocket(
    config.sockectUrl ?? "",
    String(userId ?? "")
  );

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
      <div className="flex items-center justify-center h-screen">
        User not found
      </div>
    );
  }

  return (
    <ChatContentLayout
      data={profile}
      messages={messages}
      connected={connected}
      sendMessage={sendMessage}
    />
  );
}