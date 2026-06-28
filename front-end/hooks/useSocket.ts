"use client"
import { ChatMessages } from "@/types"
import { useEffect, useRef, useState } from "react"
import { Api } from "@/services/axios"

export interface SocketType {
    connected: boolean;
    messages: ChatMessages[];
    sendMessage: (receiverId: string, content: string) => void;
}

interface IncomingPayload {
    type: string;
    message: {
        id: string;
        sender_id: string;
        sender_name?: string;
        sender_avatar?: string;
        receiver_id: string;
        content: string;
        created_at: string;
    };
}

export default function useSocket(url: string, otherUserId: string): SocketType {
    const [connected, setConnected] = useState<boolean>(false)
    const [messages, setMessages] = useState<ChatMessages[]>([])
    const socketRef = useRef<WebSocket | null>(null)

    // load existing conversation history when the chat opens
    useEffect(() => {
        const fetchHistory = async () => {
            try {
                const res = await Api.get(`/chat/private/${otherUserId}`)
                const history: ChatMessages[] = (res.data ?? []).map((m: any) => ({
                    messageId: m.id,
                    senderId: m.sender_id,
                    senderName: m.sender_name,
                    senderAvatar: m.sender_avatar,
                    receiverId: m.receiver_id,
                    content: m.content,
                    createdAt: m.created_at,
                }))
                setMessages(history)
            } catch (err) {
                console.error("Failed to load chat history:", err)
            }
        }
        fetchHistory()
    }, [otherUserId])

    // open the live websocket connection — receive-only on this side
    useEffect(() => {
        const socket = new WebSocket(url)
        socketRef.current = socket

        socket.onopen = () => setConnected(true)
        socket.onclose = () => setConnected(false)
        socket.onerror = (err: Event) => console.error("ws error:", err)

        socket.onmessage = (event: MessageEvent) => {
            try {
                const payload: IncomingPayload = JSON.parse(event.data)
                if (payload.type !== "private_message") return

                const m = payload.message
                // only add it if it belongs to the conversation currently open
                if (m.sender_id !== otherUserId && m.receiver_id !== otherUserId) return

                setMessages((prev) => [
                    ...prev,
                    {
                        messageId: m.id,
                        senderId: m.sender_id,
                        senderName: m.sender_name,
                        senderAvatar: m.sender_avatar,
                        receiverId: m.receiver_id,
                        content: m.content,
                        createdAt: m.created_at,
                    },
                ])
            } catch (err) {
                console.error("Failed to parse incoming message:", err)
            }
        }

        return () => socket.close()
    }, [url, otherUserId])

    // sending goes through the real REST endpoint — the backend handles
    // permission checks and pushes the message to the recipient's socket.
    // We optimistically add it to our own list immediately too.
    const sendMessage = async (receiverId: string, content: string) => {
        try {
            const res = await Api.post("/chat/private", {
                receiver_id: receiverId,
                content,
            })
            const m = res.data
            setMessages((prev) => [
                ...prev,
                {
                    messageId: m.id,
                    senderId: m.sender_id,
                    receiverId: m.receiver_id,
                    content: m.content,
                    createdAt: m.created_at,
                },
            ])
        } catch (err) {
            console.error("Failed to send message:", err)
        }
    }

    return {
        connected,
        messages,
        sendMessage,
    }
}