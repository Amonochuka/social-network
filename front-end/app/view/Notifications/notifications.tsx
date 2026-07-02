"use client";

import { useEffect, useState } from "react";
import Avatar from "@/components/ui/avatar";
import { Api } from "@/services/axios";

interface NotificationDetail {
  id: string;
  user_id: string;
  actor_id: string;
  actor_name: string;
  actor_avatar: string;
  type: string;
  reference_id: string;
  is_read: boolean;
  created_at: string;
}

function timeAgo(dateString: string): string {
  const date = new Date(dateString);
  const seconds = Math.floor((Date.now() - date.getTime()) / 1000);

  if (seconds < 60) return `${seconds}s ago`;
  const minutes = Math.floor(seconds / 60);
  if (minutes < 60) return `${minutes}m ago`;
  const hours = Math.floor(minutes / 60);
  if (hours < 24) return `${hours}h ago`;
  const days = Math.floor(hours / 24);
  return `${days}d ago`;
}

function actionText(type: string): string {
  switch (type) {
    case "follow_request":
      return "requested to follow you.";
    case "follow_accepted":
      return "accepted your follow request.";
    case "follow_declined":
      return "declined your follow request.";
    case "group_invitation":
      return "invited you to join a group.";
    case "group_join_request":
      return "requested to join your group.";
    case "event_created":
      return "created a new event.";
    case "new_follower":
      return "started following you.";
    default:
      return "sent you a notification.";

  }
}

export default function NotificationsPage() {
  console.log("NotificationsPage rendered");

  const [items, setItems] = useState<NotificationDetail[]>([]);
  const [loading, setLoading] = useState(true);
  const [handled, setHandled] = useState<Record<string, "accepted" | "declined">>({});


  useEffect(() => {
  console.log("useEffect running");

  const fetchNotifications = async () => {
    console.log("Calling /notifications");

    try {
      const res = await Api.get<NotificationDetail[]>("/notifications");

console.log(JSON.stringify(res.data, null, 2));

      setItems(res.data ?? []);
    } catch (err) {
      console.error("Failed to load notifications:", err);
      setItems([]);
    } finally {
      console.log("Finished request");
      setLoading(false);
    }
  };

  fetchNotifications();
}, []);

const markAsRead = async (id: string) => {
  try {
    await Api.put(`/notifications/${id}/read`);

    setItems(prev =>
      prev.map(n =>
        n.id === id ? { ...n, is_read: true } : n
      )
    );
  } catch (err) {
    console.error("Failed to mark as read:", err);
  }
};


  const handleAccept = async (n: NotificationDetail) => {
    try {
      await Api.post(`/follow/requests/${n.reference_id}/accept`);
      setHandled((prev) => ({ ...prev, [n.id]: "accepted" }));
    } catch (err) {
      console.error("Failed to accept follow request:", err);
    }
  };

  const handleDecline = async (n: NotificationDetail) => {
    try {
      await Api.post(`/follow/requests/${n.reference_id}/decline`);
      setHandled((prev) => ({ ...prev, [n.id]: "declined" }));
    } catch (err) {
      console.error("Failed to decline follow request:", err);
    }
  };

   const handleGroupInviteAccept = async (n: NotificationDetail) => {
    try {
      await Api.post(`/groups/invitations/${n.reference_id}/accept`);
      setHandled((prev) => ({ ...prev, [n.id]: "accepted" }));
    } catch (err) {
      console.error("Failed to accept invitation:", err);
    }
  };

  const handleGroupInviteDecline = async (n: NotificationDetail) => {
    try {
      await Api.post(`/groups/invitations/${n.reference_id}/decline`);
      setHandled((prev) => ({ ...prev, [n.id]: "declined" }));
    } catch (err) {
      console.error("Failed to decline invitation:", err);
    }
  };

  if (loading) {
    return (
      <div className="flex h-[70vh] items-center justify-center">
        <p className="text-sm text-gray-500">Loading notifications...</p>
      </div>
    );
  }

 

  return (
    <div className="mx-auto max-w-2xl px-4">
      <div className="divide-y divide-[#1f1f1f]">
        {items.map((n) => (
          <div
            key={n.id}
            className="flex items-start gap-4 py-5 transition-colors hover:bg-white/[0.02]"
          >
            <Avatar name={n.actor_name} size="md" />

            <div className="min-w-0 flex-1">
              <p className="text-sm leading-relaxed text-white">
                <span className="font-semibold">{n.actor_name}</span>{" "}
                <span className="text-gray-400">{actionText(n.type)}</span>
              </p>

              <p className="mt-1 text-xs text-gray-500">
                {timeAgo(n.created_at)}
              </p>
            </div>

            <div className="shrink-0">
              {/* Follow Request */}
              {n.type === "follow_request" &&
                (handled[n.id] ? (
                  <span
                    className={`text-sm font-semibold ${
                      handled[n.id] === "accepted"
                        ? "text-[#14afa7]"
                        : "text-red-500"
                    }`}
                  >
                    {handled[n.id] === "accepted" ? "Accepted" : "Declined"}
                  </span>
                ) : (
                  <div className="flex items-center gap-2">
                    <button
                      onClick={() => handleAccept(n)}
                      className="rounded-lg bg-[#14afa7] px-4 py-2 text-xs font-semibold text-white transition hover:opacity-90"
                    >
                      Accept
                    </button>

                    <button
                      onClick={() => handleDecline(n)}
                      className="rounded-lg border border-[#3a3a3a] px-4 py-2 text-xs font-semibold text-white transition hover:bg-[#242424]"
                    >
                      Decline
                    </button>
                  </div>
                ))}

              {/* Other notification types — just informational for now */}
              {!["follow_request", "group_invitation"].includes(n.type) && !n.is_read && (
              <button
              onClick={() => markAsRead(n.id)}
              className="rounded-lg bg-[#262626] px-4 py-2 text-xs font-semibold text-white transition hover:bg-[#333]"
                  >
              Mark read
              </button> 
            )}
            {n.type === "group_invitation" &&
            (
    handled[n.id] ? (
        <span className={`text-sm font-semibold ${
            handled[n.id] === "accepted"
                ? "text-[#14afa7]"
                : "text-red-500"
        }`}>
            {handled[n.id] === "accepted"
                ? "Accepted"
                : "Declined"}
        </span>
    ) : (
        <div className="flex gap-2">
            <button
                onClick={() => handleGroupInviteAccept(n)}
                className="rounded-lg bg-[#14afa7] px-4 py-2 text-xs font-semibold text-white transition hover:opacity-90"
            >
                Accept
            </button>

            <button
                onClick={() => handleGroupInviteDecline(n)}
                className="rounded-lg border border-[#3a3a3a] px-4 py-2 text-xs font-semibold text-white transition hover:bg-[#242424]"
            >
                Decline
            </button>
        </div>
    )
)}
            </div>
          </div>
        ))}
      </div>

      {/* Empty State */}
      {items.length === 0 && (
        <div className="flex h-[70vh] items-center justify-center">
          <p className="text-sm text-gray-500">You&apos;re all caught up.</p>
        </div>
      )}
    </div>
  );
}