"use client";

import { useEffect, useState, useRef } from "react";
import { useParams } from "next/navigation";
import { Api } from "@/services/axios";
import { useAppSelector } from "@/store/hooks";
import { authSelector } from "@/store/features/authSlice";
import DefaultLayout from "@/components/layouts/defaultLayout";
import UserProfileImage from "@/components/header/profile/userProfile";
import {
  Users, Calendar, MessageSquare, FileText,
  Send, Plus, X, Check, XCircle,
  UserPlus, Clock, Heart,
} from "lucide-react";

interface GroupDetail {
  id: string;
  creator_id: string;
  title: string;
  description: string;
  created_at: string;
  members: MemberProfile[];
  is_member: boolean;
  is_creator: boolean;
  invite_status: string;
  request_status: string;
}

interface MemberProfile {
  user_id: string;
  first_name: string;
  last_name: string;
  avatar: string;
  nickname: string;
}

interface GroupPost {
  id: string;
  group_id: string;
  user_id: string;
  author_name: string;
  author_avatar: string;
  content: string;
  media_path: string;
  media_type: string;
  comment_count: number;
  like_count: number;
  liked_by_me: boolean;
  created_at: string;
}

interface EventDetail {
  id: string;
  group_id: string;
  creator_id: string;
  creator_name: string;
  title: string;
  description: string;
  event_time: string;
  created_at: string;
  going_count: number;
  not_going_count: number;
  my_response: string;
}

interface ChatMessage {
  id: string;
  group_id: string;
  sender_id: string;
  sender_name: string;
  sender_avatar: string;
  content: string;
  created_at: string;
}

interface FollowerProfile {
  user_id: string;
  first_name: string;
  last_name: string;
  avatar: string;
  nickname: string;
}

type Tab = "posts" | "events" | "chat" | "members";

export default function GroupDetailPage() {
  return (
    <DefaultLayout>
      <GroupDetailContent />
    </DefaultLayout>
  );
}

function GroupDetailContent() {
  const params = useParams();
  const groupId = params.grooupId as string;
  const { user } = useAppSelector(authSelector);

  const [group, setGroup] = useState<GroupDetail | null>(null);
  const [loading, setLoading] = useState(true);
  const [activeTab, setActiveTab] = useState<Tab>("posts");

  // posts
  const [posts, setPosts] = useState<GroupPost[]>([]);
  const [newPostText, setNewPostText] = useState("");
  const [postingPost, setPostingPost] = useState(false);

  // events
  const [events, setEvents] = useState<EventDetail[]>([]);
  const [showCreateEvent, setShowCreateEvent] = useState(false);

  // chat
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [chatText, setChatText] = useState("");
  const [sendingChat, setSendingChat] = useState(false);
  const chatEndRef = useRef<HTMLDivElement>(null);

  // invite
  const [showInviteModal, setShowInviteModal] = useState(false);

  const fetchGroup = async () => {
    try {
      const res = await Api.get<GroupDetail>(`/groups/${groupId}`);
      setGroup(res.data);
    } catch (err) {
      console.error("Failed to load group:", err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (groupId) fetchGroup();
  }, [groupId]);

  // Load tab data
  useEffect(() => {
    if (!group?.is_member) return;

    if (activeTab === "posts") {
      Api.get<GroupPost[]>(`/groups/${groupId}/posts`)
        .then((res) => setPosts(res.data ?? []))
        .catch(() => setPosts([]));
    } else if (activeTab === "events") {
      Api.get<EventDetail[]>(`/groups/${groupId}/events`)
        .then((res) => setEvents(res.data ?? []))
        .catch(() => setEvents([]));
    } else if (activeTab === "chat") {
      Api.get<ChatMessage[]>(`/groups/${groupId}/chat`)
        .then((res) => {
          setMessages(res.data ?? []);
          setTimeout(() => chatEndRef.current?.scrollIntoView({ behavior: "smooth" }), 100);
        })
        .catch(() => setMessages([]));
    }
  }, [activeTab, group?.is_member, groupId]);

  const handleJoinRequest = async () => {
    try {
      await Api.post(`/groups/${groupId}/join`);
      fetchGroup();
    } catch (err) {
      console.error("Failed to request join:", err);
    }
  };

  const handleAcceptInvitation = async () => {
    if (!group) return;
    try {
      // We need the invitation ID — for now use the endpoint that accepts by group
      // The invite_status tells us there's a pending invitation
      const res = await Api.get<any>(`/groups/${groupId}`);
      // Actually we need a different approach — let's refetch and use the invitation endpoint
      await Api.post(`/groups/invitations/${group.invite_status}/accept`);
      fetchGroup();
    } catch (err) {
      console.error("Failed to accept invitation:", err);
    }
  };

  const handleCreatePost = async () => {
    if (!newPostText.trim() || postingPost) return;
    setPostingPost(true);
    try {
      const formData = new FormData();
      formData.append("content", newPostText);
      await Api.post(`/groups/${groupId}/posts`, formData, {
        headers: { "Content-Type": "multipart/form-data" },
      });
      setNewPostText("");
      const res = await Api.get<GroupPost[]>(`/groups/${groupId}/posts`);
      setPosts(res.data ?? []);
    } catch (err) {
      console.error("Failed to create post:", err);
    } finally {
      setPostingPost(false);
    }
  };

  const handleSendChat = async () => {
    if (!chatText.trim() || sendingChat) return;
    setSendingChat(true);
    try {
      await Api.post(`/groups/${groupId}/chat`, { content: chatText });
      setChatText("");
      const res = await Api.get<ChatMessage[]>(`/groups/${groupId}/chat`);
      setMessages(res.data ?? []);
      setTimeout(() => chatEndRef.current?.scrollIntoView({ behavior: "smooth" }), 100);
    } catch (err) {
      console.error("Failed to send message:", err);
    } finally {
      setSendingChat(false);
    }
  };

  const handleRSVP = async (eventId: string, response: string) => {
    try {
      await Api.post(`/groups/events/${eventId}/rsvp`, { response });
      const res = await Api.get<EventDetail[]>(`/groups/${groupId}/events`);
      setEvents(res.data ?? []);
    } catch (err) {
      console.error("Failed to RSVP:", err);
    }
  };

  if (loading) {
    return <div className="flex h-screen w-full items-center justify-center text-gray-400">Loading group...</div>;
  }

  if (!group) {
    return <div className="flex h-screen w-full items-center justify-center text-gray-400">Group not found.</div>;
  }

  return (
    <div className="flex h-screen w-full flex-col overflow-y-auto bg-[#181818]">
      {/* Header */}
      <div className="bg-gradient-to-r from-violet-950 to-indigo-950 px-8 py-8">
        <div className="flex items-start justify-between">
          <div>
            <h1 className="text-3xl font-bold text-white">{group.title}</h1>
            <p className="mt-2 text-sm text-gray-300">{group.description}</p>
            <div className="mt-3 flex items-center gap-4 text-sm text-gray-400">
              <span className="flex items-center gap-1">
                <Users size={14} /> {group.members?.length ?? 0} members
              </span>
              <span>Created {new Date(group.created_at).toLocaleDateString()}</span>
            </div>
          </div>

          <div className="flex gap-2">
            {group.is_member && (
              <button
                onClick={() => setShowInviteModal(true)}
                className="flex items-center gap-2 rounded-xl bg-white/10 px-4 py-2 text-sm font-semibold text-white transition hover:bg-white/15"
              >
                <UserPlus size={16} />
                Invite
              </button>
            )}

            {!group.is_member && group.request_status === "pending" && (
              <div className="flex items-center gap-2 rounded-xl bg-white/10 px-4 py-2 text-sm text-gray-300">
                <Clock size={16} />
                Request Pending
              </div>
            )}

            {!group.is_member && group.invite_status === "pending" && (
              <button
                onClick={handleAcceptInvitation}
                className="flex items-center gap-2 rounded-xl bg-[--primary-theme] px-4 py-2 text-sm font-semibold text-black transition hover:opacity-90"
              >
                <Check size={16} />
                Accept Invitation
              </button>
            )}

            {!group.is_member && !group.request_status && !group.invite_status && (
              <button
                onClick={handleJoinRequest}
                className="flex items-center gap-2 rounded-xl bg-[--primary-theme] px-4 py-2 text-sm font-semibold text-black transition hover:opacity-90"
              >
                <Plus size={16} />
                Request to Join
              </button>
            )}
          </div>
        </div>

        {/* Tabs */}
        {group.is_member && (
          <div className="mt-6 flex gap-1">
            {(["posts", "events", "chat", "members"] as Tab[]).map((tab) => {
              const icons = { posts: FileText, events: Calendar, chat: MessageSquare, members: Users };
              const Icon = icons[tab];
              return (
                <button
                  key={tab}
                  onClick={() => setActiveTab(tab)}
                  className={`flex items-center gap-2 rounded-lg px-4 py-2 text-sm font-medium transition ${
                    activeTab === tab
                      ? "bg-white/15 text-white"
                      : "text-gray-400 hover:text-white hover:bg-white/5"
                  }`}
                >
                  <Icon size={15} />
                  {tab.charAt(0).toUpperCase() + tab.slice(1)}
                </button>
              );
            })}
          </div>
        )}
      </div>

      {/* Content */}
      {!group.is_member ? (
        <div className="flex flex-1 items-center justify-center text-gray-500">
          Join this group to see its content.
        </div>
      ) : (
        <div className="flex-1 p-6">
          {/* Posts tab */}
          {activeTab === "posts" && (
            <div className="space-y-4">
              {/* Create post */}
              <div className="rounded-xl bg-[#222] p-4 border border-white/5">
                <div className="flex gap-3">
                  <input
                    value={newPostText}
                    onChange={(e) => setNewPostText(e.target.value)}
                    onKeyDown={(e) => e.key === "Enter" && handleCreatePost()}
                    placeholder="Write something to the group..."
                    className="flex-1 bg-[#2a2a2a] rounded-lg px-4 py-2.5 text-sm text-white outline-none"
                  />
                  <button
                    onClick={handleCreatePost}
                    disabled={postingPost}
                    className="rounded-lg bg-[--primary-theme] px-4 py-2 text-sm font-semibold text-black transition hover:opacity-90 disabled:opacity-50"
                  >
                    {postingPost ? "..." : "Post"}
                  </button>
                </div>
              </div>

              {posts.length === 0 ? (
                <div className="text-center text-sm text-gray-500 py-8">No posts yet. Be the first to post!</div>
              ) : (
                posts.map((post) => (
                  <GroupPostCard
                    key={post.id}
                    post={post}
                    currentUserId={user?.id ?? ""}
                    onDeleted={(id) => setPosts((prev) => prev.filter((p) => p.id !== id))}
                  />
                ))
              )}
            </div>
          )}

          {/* Events tab */}
          {activeTab === "events" && (
            <div className="space-y-4">
              <div className="flex justify-between items-center">
                <h2 className="text-lg font-bold text-white">Events</h2>
                <button
                  onClick={() => setShowCreateEvent(true)}
                  className="flex items-center gap-2 rounded-lg bg-[--primary-theme] px-4 py-2 text-sm font-semibold text-black transition hover:opacity-90"
                >
                  <Plus size={14} /> Create Event
                </button>
              </div>

              {events.length === 0 ? (
                <div className="text-center text-sm text-gray-500 py-8">No events yet.</div>
              ) : (
                events.map((event) => (
                  <div key={event.id} className="rounded-xl bg-[#222] p-5 border border-white/5">
                    <div className="flex justify-between items-start">
                      <div>
                        <h3 className="text-base font-bold text-white">{event.title}</h3>
                        <p className="mt-1 text-sm text-gray-400">{event.description}</p>
                        <div className="mt-2 flex items-center gap-3 text-xs text-gray-500">
                          <span className="flex items-center gap-1">
                            <Calendar size={12} />
                            {new Date(event.event_time).toLocaleString()}
                          </span>
                          <span>by {event.creator_name}</span>
                        </div>
                        <div className="mt-2 flex items-center gap-3 text-xs">
                          <span className="text-emerald-400">{event.going_count} going</span>
                          <span className="text-red-400">{event.not_going_count} not going</span>
                        </div>
                      </div>
                      <div className="flex gap-2">
                        <button
                          onClick={() => handleRSVP(event.id, "going")}
                          className={`rounded-lg px-3 py-1.5 text-xs font-semibold transition ${
                            event.my_response === "going"
                              ? "bg-emerald-500/20 text-emerald-400 border border-emerald-500/30"
                              : "bg-white/5 text-gray-400 hover:bg-white/10"
                          }`}
                        >
                          <Check size={14} />
                        </button>
                        <button
                          onClick={() => handleRSVP(event.id, "not_going")}
                          className={`rounded-lg px-3 py-1.5 text-xs font-semibold transition ${
                            event.my_response === "not_going"
                              ? "bg-red-500/20 text-red-400 border border-red-500/30"
                              : "bg-white/5 text-gray-400 hover:bg-white/10"
                          }`}
                        >
                          <XCircle size={14} />
                        </button>
                      </div>
                    </div>
                  </div>
                ))
              )}

              {showCreateEvent && (
                <CreateEventModal
                  groupId={groupId}
                  onClose={() => setShowCreateEvent(false)}
                  onCreated={() => {
                    setShowCreateEvent(false);
                    Api.get<EventDetail[]>(`/groups/${groupId}/events`)
                      .then((res) => setEvents(res.data ?? []));
                  }}
                />
              )}
            </div>
          )}

          {/* Chat tab */}
          {activeTab === "chat" && (
            <div className="flex flex-col h-[calc(100vh-320px)]">
              <div className="flex-1 overflow-y-auto space-y-3 pb-4">
                {messages.length === 0 ? (
                  <div className="text-center text-sm text-gray-500 py-8">No messages yet. Say hello!</div>
                ) : (
                  messages.map((msg) => {
                    const isMine = msg.sender_id === user?.id;
                    return (
                      <div key={msg.id} className={`flex gap-3 ${isMine ? "flex-row-reverse" : ""}`}>
                        <UserProfileImage
                          url={msg.sender_avatar ? (msg.sender_avatar.startsWith("http") ? msg.sender_avatar : `http://localhost:8080/${msg.sender_avatar}`) : undefined}
                          name={msg.sender_name}
                        />
                        <div className={`max-w-[60%] rounded-xl px-4 py-2.5 ${isMine ? "bg-[--primary-theme] text-black" : "bg-[#2a2a2a] text-white"}`}>
                          {!isMine && <p className="text-xs font-semibold mb-1 opacity-70">{msg.sender_name}</p>}
                          <p className="text-sm">{msg.content}</p>
                          <p className={`text-[10px] mt-1 ${isMine ? "text-black/50" : "text-gray-500"}`}>
                            {new Date(msg.created_at).toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" })}
                          </p>
                        </div>
                      </div>
                    );
                  })
                )}
                <div ref={chatEndRef} />
              </div>

              <div className="flex gap-3 pt-3 border-t border-white/5">
                <input
                  value={chatText}
                  onChange={(e) => setChatText(e.target.value)}
                  onKeyDown={(e) => e.key === "Enter" && handleSendChat()}
                  placeholder="Type a message..."
                  className="flex-1 bg-[#2a2a2a] rounded-xl px-4 py-3 text-sm text-white outline-none"
                />
                <button
                  onClick={handleSendChat}
                  disabled={sendingChat}
                  className="rounded-xl bg-[--primary-theme] px-4 py-2 text-black transition hover:opacity-90 disabled:opacity-50"
                >
                  <Send size={18} />
                </button>
              </div>
            </div>
          )}

          {/* Members tab */}
          {activeTab === "members" && (
            <div className="space-y-3">
              {(group.members ?? []).map((member) => (
                <div key={member.user_id} className="flex items-center gap-3 rounded-xl bg-[#222] p-4 border border-white/5">
                  <UserProfileImage
                    url={member.avatar ? (member.avatar.startsWith("http") ? member.avatar : `http://localhost:8080/${member.avatar}`) : undefined}
                    name={`${member.first_name} ${member.last_name}`}
                  />
                  <div>
                    <p className="text-sm font-semibold text-white">{member.first_name} {member.last_name}</p>
                    {member.nickname && <p className="text-xs text-gray-500">@{member.nickname}</p>}
                  </div>
                  {member.user_id === group.creator_id && (
                    <span className="ml-auto text-xs text-[--primary-theme] font-semibold">Creator</span>
                  )}
                </div>
              ))}
            </div>
          )}
        </div>
      )}

      {/* Invite modal */}
      {showInviteModal && (
        <InviteModal groupId={groupId} onClose={() => setShowInviteModal(false)} />
      )}
    </div>
  );
}

function GroupPostCard({
  post,
  currentUserId,
  onDeleted,
}: {
  post: GroupPost;
  currentUserId: string;
  onDeleted: (id: string) => void;
}) {
  const [liked, setLiked] = useState(post.liked_by_me);
  const [likeCount, setLikeCount] = useState(post.like_count);
  const [liking, setLiking] = useState(false);
  const [showComments, setShowComments] = useState(false);
  const [deleting, setDeleting] = useState(false);

  const handleLike = async () => {
    if (liking) return;
    const prevLiked = liked;
    const prevCount = likeCount;
    setLiked(!liked);
    setLikeCount(liked ? likeCount - 1 : likeCount + 1);
    setLiking(true);
    try {
      const res = await Api.post<{ liked: boolean; like_count: number }>(
        `/group-posts/${post.id}/like`
      );
      setLiked(res.data.liked);
      setLikeCount(res.data.like_count);
    } catch {
      setLiked(prevLiked);
      setLikeCount(prevCount);
    } finally {
      setLiking(false);
    }
  };

  const handleDelete = async () => {
    if (!confirm("Delete this post?")) return;
    setDeleting(true);
    try {
      await Api.delete(`/group-posts/${post.id}`);
      onDeleted(post.id);
    } catch (err) {
      console.error("Failed to delete post:", err);
    } finally {
      setDeleting(false);
    }
  };

  return (
    <div className="rounded-xl bg-[#222] p-5 border border-white/5">
      <div className="flex items-center gap-3 mb-3">
        <UserProfileImage
          url={
            post.author_avatar
              ? post.author_avatar.startsWith("http")
                ? post.author_avatar
                : `http://localhost:8080/${post.author_avatar}`
              : undefined
          }
          name={post.author_name}
        />
        <div className="flex-1">
          <p className="text-sm font-semibold text-white">{post.author_name}</p>
          <p className="text-xs text-gray-500">{new Date(post.created_at).toLocaleDateString()}</p>
        </div>
        {post.user_id === currentUserId && (
          <button
            onClick={handleDelete}
            disabled={deleting}
            className="text-xs text-gray-500 hover:text-red-400 transition-colors px-2 py-1 rounded-lg hover:bg-red-400/10"
          >
            {deleting ? "..." : "Delete"}
          </button>
        )}
      </div>

      <p className="text-sm text-gray-300 leading-relaxed">{post.content}</p>

      {/* Interactions */}
      <div className="flex items-center gap-8 mt-4 pt-3 border-t border-white/5">
        <button
          onClick={handleLike}
          disabled={liking}
          className="group flex items-center gap-1.5 transition-colors"
        >
          <Heart
            size={17}
            fill={liked ? "currentColor" : "none"}
            className={liked ? "text-pink-500" : "text-gray-500 group-hover:text-pink-500 transition-colors"}
          />
          <span className={`text-sm font-medium transition-colors ${liked ? "text-pink-500" : "text-gray-500 group-hover:text-pink-500"}`}>
            {likeCount}
          </span>
        </button>

        <button
          onClick={() => setShowComments((p) => !p)}
          className="group flex items-center gap-1.5 transition-colors"
        >
          <MessageSquare
            size={17}
            className="text-gray-500 group-hover:text-sky-400 transition-colors"
          />
          <span className="text-sm font-medium text-gray-500 group-hover:text-sky-400 transition-colors">
            Comment
          </span>
        </button>
      </div>

      {showComments && <GroupCommentSection postId={post.id} />}
    </div>
  );
}

interface GroupCommentDetail {
  id: string;
  group_post_id: string;
  user_id: string;
  author_name: string;
  author_avatar: string;
  content: string;
  media_path: string;
  media_type: string;
  created_at: string;
}

function GroupCommentSection({ postId }: { postId: string }) {
  const [comments, setComments] = useState<GroupCommentDetail[]>([]);
  const [loading, setLoading] = useState(true);
  const [text, setText] = useState("");
  const [posting, setPosting] = useState(false);

  useEffect(() => {
    Api.get<GroupCommentDetail[]>(`/group-posts/${postId}/comments`)
      .then((res) => setComments(res.data ?? []))
      .catch(() => setComments([]))
      .finally(() => setLoading(false));
  }, [postId]);

  const handleAddComment = async () => {
    if (!text.trim() || posting) return;
    setPosting(true);
    try {
      const formData = new FormData();
      formData.append("content", text);
      const res = await Api.post<GroupCommentDetail>(
        `/group-posts/${postId}/comments`,
        formData,
        { headers: { "Content-Type": "multipart/form-data" } }
      );
      setComments((prev) => [...prev, res.data]);
      setText("");
    } catch (err) {
      console.error("Failed to add comment:", err);
    } finally {
      setPosting(false);
    }
  };

  return (
    <div className="mt-3 pt-3 border-t border-white/10 flex flex-col gap-3">
      {loading ? (
        <p className="text-sm text-gray-500 animate-pulse text-center">Loading...</p>
      ) : comments.length === 0 ? (
        <p className="text-sm text-gray-500 text-center">No comments yet.</p>
      ) : (
        <div className="flex flex-col gap-3">
          {comments.map((c) => (
            <div key={c.id} className="flex gap-2">
              <div className="shrink-0 pt-0.5">
                <UserProfileImage
                  url={
                    c.author_avatar
                      ? c.author_avatar.startsWith("http")
                        ? c.author_avatar
                        : `http://localhost:8080/${c.author_avatar}`
                      : undefined
                  }
                  name={c.author_name}
                />
              </div>
              <div className="bg-white/5 rounded-2xl px-4 py-2 text-sm min-w-0 max-w-full">
                <p className="font-bold text-white">{c.author_name}</p>
                <p className="text-gray-300 mt-0.5 break-words">{c.content}</p>
              </div>
            </div>
          ))}
        </div>
      )}
      <div className="flex items-center gap-3 mt-1">
        <input
          className="flex-1 bg-transparent border border-white/10 rounded-full px-4 py-2 text-sm text-white outline-none focus:border-[--primary-theme]"
          type="text"
          placeholder="Write a reply..."
          value={text}
          onChange={(e) => setText(e.target.value)}
          onKeyDown={(e) => e.key === "Enter" && handleAddComment()}
        />
        <button
          className="bg-[--primary-theme] text-white font-bold text-sm px-5 py-2 rounded-full transition-colors hover:opacity-90 disabled:opacity-50"
          onClick={handleAddComment}
          disabled={posting || !text.trim()}
        >
          {posting ? "..." : "Reply"}
        </button>
      </div>
    </div>
  );
}

function CreateEventModal({ groupId, onClose, onCreated }: { groupId: string; onClose: () => void; onCreated: () => void }) {
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [eventTime, setEventTime] = useState("");
  const [creating, setCreating] = useState(false);

  const handleCreate = async () => {
    if (!title.trim() || !eventTime) return;
    setCreating(true);
    try {
      await Api.post(`/groups/${groupId}/events`, { title, description, event_time: eventTime });
      onCreated();
    } catch (err) {
      console.error("Failed to create event:", err);
    } finally {
      setCreating(false);
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 p-4">
      <div className="w-full max-w-md rounded-2xl bg-[#1e1e1e] p-6 shadow-2xl">
        <div className="mb-5 flex items-center justify-between">
          <h2 className="text-lg font-bold text-white">Create Event</h2>
          <button onClick={onClose} className="rounded-lg p-1.5 text-gray-400 hover:text-white"><X size={18} /></button>
        </div>
        <div className="space-y-3">
          <div>
            <label className="mb-1 block text-xs text-gray-400">Event name</label>
            <input value={title} onChange={(e) => setTitle(e.target.value)}
              className="w-full rounded-lg bg-[#262626] px-3 py-2.5 text-sm text-white outline-none focus:ring-1 focus:ring-[--primary-theme]" />
          </div>
          <div>
            <label className="mb-1 block text-xs text-gray-400">Description</label>
            <textarea value={description} onChange={(e) => setDescription(e.target.value)} rows={2}
              className="w-full resize-none rounded-lg bg-[#262626] px-3 py-2.5 text-sm text-white outline-none focus:ring-1 focus:ring-[--primary-theme]" />
          </div>
          <div>
            <label className="mb-1 block text-xs text-gray-400">Date & Time</label>
            <input type="datetime-local" value={eventTime} onChange={(e) => setEventTime(e.target.value)}
              className="w-full rounded-lg bg-[#262626] px-3 py-2.5 text-sm text-white outline-none focus:ring-1 focus:ring-[--primary-theme]" />
          </div>
        </div>
        <div className="mt-6 flex justify-end gap-2">
          <button onClick={onClose} className="rounded-lg px-4 py-2 text-sm font-semibold text-gray-300 hover:bg-white/5">Cancel</button>
          <button onClick={handleCreate} disabled={creating}
            className="rounded-lg bg-[--primary-theme] px-4 py-2 text-sm font-semibold text-black hover:opacity-90 disabled:opacity-50">
            {creating ? "Creating..." : "Create"}
          </button>
        </div>
      </div>
    </div>
  );
}

function InviteModal({ groupId, onClose }: { groupId: string; onClose: () => void }) {
  const [followers, setFollowers] = useState<FollowerProfile[]>([]);
  const [inviting, setInviting] = useState<string | null>(null);
  const [invited, setInvited] = useState<Set<string>>(new Set());

  useEffect(() => {
    Api.get<FollowerProfile[]>("/following")
      .then((res) => setFollowers(res.data ?? []))
      .catch(() => setFollowers([]));
  }, []);

  const handleInvite = async (userId: string) => {
    setInviting(userId);
    try {
      await Api.post(`/groups/${groupId}/invite`, { invitee_id: userId });
      setInvited((prev) => new Set(prev).add(userId));
    } catch (err) {
      console.error("Failed to invite:", err);
    } finally {
      setInviting(null);
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 p-4">
      <div className="w-full max-w-md rounded-2xl bg-[#1e1e1e] p-6 shadow-2xl max-h-[80vh] flex flex-col">
        <div className="mb-5 flex items-center justify-between">
          <h2 className="text-lg font-bold text-white">Invite to Group</h2>
          <button onClick={onClose} className="rounded-lg p-1.5 text-gray-400 hover:text-white"><X size={18} /></button>
        </div>
        <div className="flex-1 overflow-y-auto space-y-2">
          {followers.length === 0 ? (
            <p className="text-sm text-gray-500 text-center py-4">No followers to invite.</p>
          ) : (
            followers.map((f) => (
              <div key={f.user_id} className="flex items-center justify-between rounded-lg bg-[#262626] px-4 py-3">
                <div className="flex items-center gap-3">
                  <UserProfileImage
                    url={f.avatar ? (f.avatar.startsWith("http") ? f.avatar : `http://localhost:8080/${f.avatar}`) : undefined}
                    name={`${f.first_name} ${f.last_name}`}
                  />
                  <span className="text-sm text-white">{f.first_name} {f.last_name}</span>
                </div>
                <button
                  onClick={() => handleInvite(f.user_id)}
                  disabled={invited.has(f.user_id) || inviting === f.user_id}
                  className={`rounded-lg px-3 py-1.5 text-xs font-semibold transition ${
                    invited.has(f.user_id)
                      ? "bg-emerald-500/20 text-emerald-400"
                      : "bg-[--primary-theme] text-black hover:opacity-90"
                  } disabled:opacity-50`}
                >
                  {invited.has(f.user_id) ? "Invited" : inviting === f.user_id ? "..." : "Invite"}
                </button>
              </div>
            ))
          )}
        </div>
      </div>
    </div>
  );
}