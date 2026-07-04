"use client"

import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useRef,
  useState,
} from "react"
import { ChatMessages } from "@/types"
import { Api } from "@/services/axios"
import loadEnvFile from "@/config/config"
import { useAppSelector } from "@/store/hooks"
import { authSelector } from "@/store/features/authSlice"

interface SocketContextValue {
  connected: boolean
  onlineUsers: string[]
  getMessages: (conversationId: string) => ChatMessages[]
  loadConversation: (otherUserId: string) => Promise<void>
  sendMessage: (receiverId: string, content: string) => Promise<void>
}

const SocketContext = createContext<SocketContextValue | undefined>(undefined)

export function WebSocketProvider({ children }: { children: React.ReactNode }) {
  const [connected, setConnected] = useState(false)
  const [onlineUsers, setOnlineUsers] = useState<string[]>([])
  const [messagesByUser, setMessagesByUser] = useState<
    Record<string, ChatMessages[]>
  >({})
  const socketRef = useRef<WebSocket | null>(null)
  const auth = useAppSelector(authSelector)
  const currentUserId = auth.user?.id

  const config = loadEnvFile({ urlType: "socket-url", urlUsage: "chat" })
  const socketUrl = config.sockectUrl ?? ""

  const addMessageToConversation = useCallback(
    (conversationId: string, message: ChatMessages) => {
      setMessagesByUser((prev) => {
        const existing = prev[conversationId] ?? []
        return {
          ...prev,
          [conversationId]: [...existing, message],
        }
      })
    },
    []
  )

  const getMessages = useCallback(
    (conversationId: string) => messagesByUser[conversationId] ?? [],
    [messagesByUser]
  )

  const loadConversation = useCallback(async (otherUserId: string) => {
    if (!otherUserId) return

    try {
      const res = await Api.get(`/chat/private/${otherUserId}`)
      const history: ChatMessages[] = (res.data ?? []).map((m: any) => ({
        messageId: m.id,
        senderId: String(m.sender_id),
        senderName: m.sender_name,
        senderAvatar: m.sender_avatar,
        receiverId: String(m.receiver_id),
        content: m.content,
        createdAt: m.created_at,
      }))

      setMessagesByUser((prev) => ({
        ...prev,
        [otherUserId]: history,
      }))
    } catch (err) {
      console.error("WebSocketProvider: failed to load conversation", err)
    }
  }, [])

  const sendMessage = useCallback(
    async (receiverId: string, content: string) => {
      if (!receiverId || !content.trim()) return

      try {
        const res = await Api.post("/chat/private", {
          receiver_id: receiverId,
          content,
        })
        const m = res.data

        addMessageToConversation(receiverId, {
          messageId: m.id,
          senderId: String(m.sender_id),
          senderName: m.sender_name,
          senderAvatar: m.sender_avatar,
          receiverId: String(m.receiver_id),
          content: m.content,
          createdAt: m.created_at,
        })
      } catch (err) {
        console.error("WebSocketProvider: failed to send message", err)
      }
    },
    [addMessageToConversation]
  )

  useEffect(() => {
    if (!socketUrl || socketRef.current) return
    if (typeof window === "undefined") return

    const sessionId = window.localStorage.getItem("session_id")
    if (!sessionId) {
      console.warn(
        "WebSocketProvider: session_id not found, websocket will start after authentication"
      )
      return
    }

    const url = socketUrl.includes("?")
      ? `${socketUrl}&session_id=${encodeURIComponent(sessionId)}`
      : `${socketUrl}?session_id=${encodeURIComponent(sessionId)}`

    const socket = new WebSocket(url)
    socketRef.current = socket

    socket.onopen = () => {
      setConnected(true)
      console.log("WebSocketProvider: connected")
    }

    socket.onclose = (event) => {
      console.warn("WebSocketProvider: closed", event.code, event.reason)
      setConnected(false)
    }

    socket.onerror = (event) => {
      console.error("WebSocketProvider: error", event)
    }

    socket.onmessage = (event) => {
      try {
        const payload = JSON.parse(event.data)

        if (payload.type === "private_message" && payload.message) {
          const m = payload.message
          const senderId = String(m.sender_id)
          const receiverId = String(m.receiver_id)

          const conversationId =
            currentUserId && senderId === currentUserId
              ? receiverId
              : senderId

          addMessageToConversation(conversationId, {
            messageId: m.id,
            senderId,
            senderName: m.sender_name,
            senderAvatar: m.sender_avatar,
            receiverId,
            content: m.content,
            createdAt: m.created_at,
          })
          return
        }

        if (payload.type === "presence_update") {
          if (Array.isArray(payload.online_users)) {
            setOnlineUsers(payload.online_users.map(String))
            return
          }

          if (Array.isArray(payload.users)) {
            setOnlineUsers(payload.users.map(String))
            return
          }
        }
      } catch (err) {
        console.error("WebSocketProvider: failed to parse incoming message", err)
      }
    }

    return () => {
      socketRef.current?.close()
      socketRef.current = null
    }
  }, [socketUrl, currentUserId, addMessageToConversation])

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
  )
}

export function useSocketContext() {
  const context = useContext(SocketContext)
  if (!context) {
    throw new Error(
      "useSocketContext must be used within a WebSocketProvider"
    )
  }
  return context
}
