"use client";

import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useRef,
  useState,
} from "react";
import { ChatMessages } from "@/types";
import { Api } from "@/services/axios";
import loadEnvFile from "@/config/config";
import { useAppSelector } from "@/store/hooks";
import { authSelector } from "@/store/features/authSlice";

interface SocketContextValue {
  connected: boolean;
  onlineUsers: Set<string>;
  getMessages: (conversationId: string) => ChatMessages[];
  loadConversation: (otherUserId: string) => Promise<void>;
  sendMessage: (receiverId: string, content: string) => Promise<void>;
}

const SocketContext = createContext<SocketContextValue | undefined>(undefined);

export function WebSocketProvider({
  children,
}: {
  children: React.ReactNode;
}) {
  const [connected, setConnected] = useState(false);
  const [onlineUsers, setOnlineUsers] = useState<Set<string>>(new Set());
  const [messagesByUser, setMessagesByUser] = useState<
    Record<string, ChatMessages[]>
  >({});

  const socketRef = useRef<WebSocket | null>(null);

  const auth = useAppSelector(authSelector);
  const currentUserId = auth.user?.id;

  const config = loadEnvFile({
    urlType: "socket-url",
    urlUsage: "chat",
  });

  const socketUrl = config.sockectUrl ?? "";

  const addMessageToConversation = useCallback(
    (conversationId: string, message: ChatMessages) => {
      setMessagesByUser((prev) => {
        const existing = prev[conversationId] ?? [];

        if (existing.some((m) => m.messageId === message.messageId)) {
          return prev;
        }

        return {
          ...prev,
          [conversationId]: [...existing, message],
        };
      });
    },
    []
  );

  const getMessages = useCallback(
    (conversationId: string) => messagesByUser[conversationId] ?? [],
    [messagesByUser]
  );

  const loadConversation = useCallback(
    async (otherUserId: string) => {
      if (!otherUserId) return;

      if (messagesByUser[otherUserId]) return;

      try {
        const res = await Api.get(`/chat/private/${otherUserId}`);

        const history: ChatMessages[] = (res.data ?? []).map((m: any) => ({
          messageId: m.id,
          senderId: String(m.sender_id),
          senderName: m.sender_name,
          senderAvatar: m.sender_avatar,
          receiverId: String(m.receiver_id),
          content: m.content,
          createdAt: m.created_at,
        }));

        setMessagesByUser((prev) => ({
          ...prev,
          [otherUserId]: history,
        }));
      } catch (err) {
        console.error("Failed to load conversation", err);
      }
    },
    [messagesByUser]
  );

  const sendMessage = useCallback(
    async (receiverId: string, content: string) => {
      if (!receiverId || !content.trim()) return;

      try {
        const res = await Api.post("/chat/private", {
          receiver_id: receiverId,
          content,
        });

        const m = res.data;

        addMessageToConversation(receiverId, {
          messageId: m.id,
          senderId: String(m.sender_id),
          senderName: m.sender_name,
          senderAvatar: m.sender_avatar,
          receiverId: String(m.receiver_id),
          content: m.content,
          createdAt: m.created_at,
        });
      } catch (err) {
        console.error("Failed to send message", err);
      }
    },
    [addMessageToConversation]
  );

  useEffect(() => {
    if (!socketUrl) return;
    if (socketRef.current) return;
    if (typeof window === "undefined") return;

    console.log("Opening websocket:", socketUrl);

    const socket = new WebSocket(socketUrl);

    socketRef.current = socket;

    socket.onopen = () => {
      console.log("✅ WebSocket connected");
      setConnected(true);
    };

    socket.onclose = (event) => {
      console.warn(
        "❌ WebSocket closed",
        event.code,
        event.reason
      );

      setConnected(false);
      socketRef.current = null;
    };

    socket.onerror = (event) => {
      console.error("WebSocket error", event);
    };

    socket.onmessage = (event) => {
      try {
        const payload = JSON.parse(event.data);

        switch (payload.type) {
          case "private_message": {
            const m = payload.message;

            if (!m) return;

            const senderId = String(m.sender_id);
            const receiverId = String(m.receiver_id);

            const conversationId =
              senderId === currentUserId
                ? receiverId
                : senderId;

            addMessageToConversation(conversationId, {
              messageId: m.id,
              senderId,
              senderName: m.sender_name,
              senderAvatar: m.sender_avatar,
              receiverId,
              content: m.content,
              createdAt: m.created_at,
            });

            break;
          }

          case "presence": {
            const userId = String(payload.user_id);

            setOnlineUsers((prev) => {
              const next = new Set(prev);

              if (payload.is_online) {
                next.add(userId);
              } else {
                next.delete(userId);
              }

              return next;
            });

            break;
          }

          default:
            console.log("Unknown websocket event:", payload);
        }
      } catch (err) {
        console.error("Failed to parse websocket payload", err);
      }
    };

    return () => {
      socket.close();
      socketRef.current = null;
    };
  }, [socketUrl, currentUserId, addMessageToConversation]);

  return (
    <SocketContext.Provider
      value={{
        connected,
        onlineUsers,
        getMessages,
        loadConversation,
        sendMessage,
      }}
    >
      {children}
    </SocketContext.Provider>
  );
}

export function useSocketContext() {
  const context = useContext(SocketContext);

  if (!context) {
    throw new Error(
      "useSocketContext must be used inside WebSocketProvider"
    );
  }

  return context;
}