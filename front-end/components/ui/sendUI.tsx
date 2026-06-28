"use client"

import { Send, Smile } from "lucide-react"
import { useParams } from "next/navigation";
import { useState } from "react";

interface SendMessageType {
  sendMessage: (receiverId: string, content: string) => void;
}

export default function SendTextMessage({ sendMessage }: SendMessageType) {
    const [message, setMessage] = useState<string>("");
    const [showEmojis, setShowEmojis] = useState(false);
    const { userId } = useParams()

    const emojis = ["😀", "😂", "🥰", "😍", "😎", "😭", "😡", "👍", "👎", "🔥", "🎉", "❤️", "✨", "🙌", "🤔", "👏"];

    function HandleSend() {
      if (!message.trim() || userId === undefined) return;

      sendMessage(String(userId), message.trim())
      setMessage("")
      setShowEmojis(false)
    }

    const addEmoji = (emoji: string) => {
      setMessage((prev) => prev + emoji);
    };

  return (
    <div className="bg-[#222222] border-t border-[#333] p-4 relative">
      {showEmojis && (
        <div className="absolute bottom-[80px] left-4 bg-[#2a2a2a] border border-[#444] rounded-xl p-3 grid grid-cols-8 gap-2 shadow-2xl z-50">
          {emojis.map((emoji) => (
            <button
              key={emoji}
              onClick={() => addEmoji(emoji)}
              className="text-xl hover:scale-125 transition cursor-pointer"
            >
              {emoji}
            </button>
          ))}
        </div>
      )}

      <div className="flex items-center gap-3">
        <button
          onClick={() => setShowEmojis((prev) => !prev)}
          className="p-2 text-gray-400 hover:text-white transition cursor-pointer"
        >
          <Smile size={22} />
        </button>

        <input
         onChange={(event: React.ChangeEvent<HTMLInputElement>)=> setMessage(event.target.value)}
          onKeyDown={(event: React.KeyboardEvent<HTMLInputElement>) => {
            if (event.key === "Enter") HandleSend();
          }}
          type="text"
          value={message}
          placeholder="Type a message..."
          className="
            flex-1
            bg-[#303030]
            rounded-full
            px-5
            py-3
            outline-none
            text-white
            placeholder:text-gray-400
          "
        />

        <button
          onClick={HandleSend}
          className="
            p-3
            rounded-full
            bg-(--primary-theme)
            cursor-pointer
            transition
            hover:scale-105
          "
        >
          <Send size={18} />
        </button>
      </div>
    </div>
  );
}