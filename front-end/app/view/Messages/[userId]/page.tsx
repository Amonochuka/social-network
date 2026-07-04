"use client";

import { useParams } from "next/navigation";
import { useEffect, useState } from "react";
import ChatContentLayout from "@/components/chats/chatContentLayout";
import { useSocketContext } from "@/contexts/SocketContext";
import { Api } from "@/services/axios";

interface ChatProfile {
  id: string;
  first_name: string;
  last_name: string;
  avatar: string;
}

export default function ChatUser() {
  const { userId } = useParams();
  const [profile, setProfile] = useState<ChatProfile | null>(null);
  const [loading, setLoading] = useState(true);
  const userIdString = Array.isArray(userId) ? userId[0] : userId;

  const { connected, onlineUsers, getMessages, loadConversation, sendMessage } = useSocketContext();

  const isUserOnline = userIdString ? onlineUsers.has(userIdString) : false;

  useEffect(() => {
    if (!userId) return;
    Api.get(`/profile/${userId}`)
      .then((res) => {
        setProfile({
          id: res.data.id,
          first_name: res.data.first_name,
          last_name: res.data.last_name,
          avatar: res.data.avatar,
        });
      })
      .catch((err) => console.error("Failed to load chat profile:", err))
      .finally(() => setLoading(false));
  }, [userId]);

  useEffect(() => {
    if (!userIdString) return;
    loadConversation(userIdString);
  }, [userIdString, loadConversation]);

  if (!userIdString) return null;

  const messages = getMessages(userIdString);

  if (loading) {
    return (
      <div className="flex items-center justify-center h-screen text-gray-500 bg-[#181818]">
        Loading chat...
      </div>
    );
  }

  if (!profile) {
    return (
      <div className="flex items-center justify-center h-screen text-gray-500 bg-[#181818]">
        User not found
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
