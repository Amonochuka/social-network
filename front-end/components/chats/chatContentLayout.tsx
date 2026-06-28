import SingleUserChatProfile from "./singleUserChatProfile";
import Messages from "./messages";
import SendTextMessage from "@/components/ui/sendUI";
import ProtectedRoute from "@/components/protected/protected";
import { ChatMessages } from "@/types";

interface ChatProfile {
  id: string;
  first_name: string;
  last_name: string;
  avatar: string;
}

interface Props {
  data: ChatProfile;
  connected: boolean;
  messages: ChatMessages[];
  sendMessage: (receiverId: string, content: string) => void;
}

export default function ChatContentLayout({
  data,
  connected,
  messages,
  sendMessage,
}: Props) {
  return (
    <ProtectedRoute>
      <div className="w-full flex flex-col h-screen bg-[#181818]">
        <SingleUserChatProfile data={data} />

        <div className="flex-1 overflow-y-auto px-5 py-6">
          <Messages messages={messages} />
        </div>

        <SendTextMessage sendMessage={sendMessage} />

        {!connected && (
          <p className="text-red-500 absolute left-300 top-100">
            Connection offline...
          </p>
        )}
      </div>
    </ProtectedRoute>
  );
}