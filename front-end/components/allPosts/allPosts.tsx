"use client";

import { useEffect, useState, useRef } from "react";
import UserProfileImage from "../header/profile/userProfile";
import {
  Clapperboard,
  Eye,
  LockIcon,
  Image as LucideImage,
  Globe,
  Users,
  Lock,
  X,
  Check,
  MoreHorizontal,
} from "lucide-react";
import Image from "next/image";
import { CSSProperties } from "react";
import { PostInteractions } from "./interactions";
import { Api } from "@/services/axios";
import { useAppSelector } from "@/store/hooks";
import { authSelector } from "@/store/features/authSlice";

interface FeedPost {
  id: string;
  user_id: string;
  author_name: string;
  author_avatar: string;
  content: string;
  media_path: string;
  media_type: string;
  privacy: string;
  comment_count: number;
  created_at: string;
  updated_at: string;
}

interface FollowerOption {
  user_id: string;
  first_name: string;
  last_name: string;
  avatar: string;
  nickname: string;
}

export default function UserPostUI() {
  const { user } = useAppSelector(authSelector);
  const [postText, setPostText] = useState("");
  const [posting, setPosting] = useState(false);
  const [refreshKey, setRefreshKey] = useState(0);

  // privacy
  const [privacy, setPrivacy] = useState<"public" | "almost_private" | "private">("public");
  const [showPrivacyMenu, setShowPrivacyMenu] = useState(false);

  // private follower picker
  const [followers, setFollowers] = useState<FollowerOption[]>([]);
  const [selectedFollowers, setSelectedFollowers] = useState<string[]>([]);
  const [showFollowerPicker, setShowFollowerPicker] = useState(false);

  // media
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [mediaFile, setMediaFile] = useState<File | null>(null);
  const [mediaPreview, setMediaPreview] = useState<string | null>(null);

  useEffect(() => {
    if (privacy === "private") {
      Api.get<FollowerOption[]>("/followers")
        .then((res) => setFollowers(res.data ?? []))
        .catch(() => setFollowers([]));
      setShowFollowerPicker(true);
    } else {
      setShowFollowerPicker(false);
      setSelectedFollowers([]);
    }
  }, [privacy]);

  const handleMediaSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;
    setMediaFile(file);
    setMediaPreview(URL.createObjectURL(file));
  };

  const removeMedia = () => {
    setMediaFile(null);
    setMediaPreview(null);
    if (fileInputRef.current) fileInputRef.current.value = "";
  };

  const toggleFollower = (id: string) => {
    setSelectedFollowers((prev) =>
      prev.includes(id) ? prev.filter((f) => f !== id) : [...prev, id]
    );
  };

  const handleCreatePost = async () => {
    if ((!postText.trim() && !mediaFile) || posting) return;
    setPosting(true);

    try {
      const formData = new FormData();
      formData.append("content", postText);
      formData.append("privacy", privacy);

      if (mediaFile) {
        formData.append("media", mediaFile);
      }

      if (privacy === "private" && selectedFollowers.length > 0) {
        formData.append("allowed_users", JSON.stringify(selectedFollowers));
      }

      await Api.post("/posts", formData, {
        headers: { "Content-Type": "multipart/form-data" },
      });

      setPostText("");
      setPrivacy("public");
      setSelectedFollowers([]);
      removeMedia();
      setRefreshKey((k) => k + 1);
    } catch (err) {
      console.error("Failed to create post:", err);
    } finally {
      setPosting(false);
    }
  };

  const myAvatarUrl = user?.avatar
    ? user.avatar.startsWith("http")
      ? user.avatar
      : `http://localhost:8080/${user.avatar}`
    : undefined;

  const myFullName = user ? `${user.first_name} ${user.last_name}` : "?";

  const privacyLabels = {
    public: { icon: Globe, label: "Everyone can reply" },
    almost_private: { icon: Users, label: "Followers only" },
    private: { icon: Lock, label: "Selected followers" },
  };

  const PrivacyIcon = privacyLabels[privacy].icon;

  return (
    <div className="flex flex-col w-full border-t border-white/10 sm:border-t-0">
      {/* Create Post Box (X Style) */}
      <div className="flex gap-4 p-4 border-b border-white/10">
        <div className="shrink-0 pt-1">
          <UserProfileImage url={myAvatarUrl} name={myFullName} />
        </div>
        <div className="flex flex-col w-full pt-1">
          <textarea
            className="w-full bg-transparent text-xl outline-none placeholder:text-gray-500 resize-none overflow-hidden min-h-[50px]"
            placeholder="What is happening?!"
            value={postText}
            onChange={(e) => {
              setPostText(e.target.value);
              e.target.style.height = "auto";
              e.target.style.height = e.target.scrollHeight + "px";
            }}
            rows={1}
          />

          {mediaPreview && (
            <div className="relative mt-3 w-full rounded-2xl overflow-hidden border border-white/10">
              <img
                src={mediaPreview}
                alt="preview"
                className="w-full max-h-[500px] object-cover"
              />
              <button
                onClick={removeMedia}
                className="absolute top-2 right-2 bg-black/70 hover:bg-black/90 p-1.5 rounded-full backdrop-blur-md transition-colors"
              >
                <X size={18} className="text-white" />
              </button>
            </div>
          )}

          {/* Privacy & Follower Picker Section */}
          <div className="mt-3 relative w-fit border-b border-white/10 pb-3 mb-2">
             <button
              onClick={() => setShowPrivacyMenu((p) => !p)}
              className="flex items-center gap-1.5 text-[--primary-theme] text-sm font-bold hover:bg-[--primary-theme]/10 rounded-full px-3 py-1 -ml-3 transition-colors"
            >
              <PrivacyIcon size={16} />
              {privacyLabels[privacy].label}
            </button>

            {showPrivacyMenu && (
              <div className="absolute top-full left-0 mt-1 bg-black shadow-[0_0_15px_rgba(255,255,255,0.1)] border border-white/20 rounded-xl overflow-hidden z-50 min-w-[220px]">
                {(["public", "almost_private", "private"] as const).map((p) => {
                  const Icon = privacyLabels[p].icon;
                  return (
                    <button
                      key={p}
                      onClick={() => { setPrivacy(p); setShowPrivacyMenu(false); }}
                      className="flex items-center gap-3 w-full px-4 py-3 hover:bg-white/5 transition-colors text-left"
                    >
                      <div className={`p-2 rounded-full ${privacy === p ? 'bg-[--primary-theme] text-white' : 'bg-white/10 text-white'}`}>
                         <Icon size={18} />
                      </div>
                      <span className="text-sm font-bold text-white flex-1">{privacyLabels[p].label}</span>
                      {privacy === p && <Check size={18} className="text-[--primary-theme]" />}
                    </button>
                  );
                })}
              </div>
            )}
          </div>

          {showFollowerPicker && (
            <div className="mb-3 p-3 bg-white/5 rounded-xl border border-white/10">
              <p className="text-xs font-bold text-gray-400 mb-2 uppercase tracking-wider">Select followers</p>
              <div className="flex flex-wrap gap-2">
                {followers.length === 0 ? (
                  <p className="text-sm text-gray-500">No followers yet.</p>
                ) : (
                  followers.map((f) => {
                    const selected = selectedFollowers.includes(f.user_id);
                    return (
                      <button
                        key={f.user_id}
                        onClick={() => toggleFollower(f.user_id)}
                        className={`flex items-center gap-1.5 px-3 py-1.5 rounded-full text-sm font-medium transition-colors border ${
                          selected 
                            ? "bg-[--primary-theme]/20 border-[--primary-theme] text-[--primary-theme]" 
                            : "bg-transparent border-white/20 text-gray-300 hover:bg-white/5"
                        }`}
                      >
                        {f.first_name} {f.last_name}
                      </button>
                    );
                  })
                )}
              </div>
            </div>
          )}

          <div className="flex items-center justify-between pt-2">
            <div className="flex items-center gap-1 -ml-2">
              <button 
                onClick={() => fileInputRef.current?.click()}
                className="p-2 rounded-full hover:bg-[--primary-theme]/10 text-[--primary-theme] transition-colors"
                title="Media"
              >
                <LucideImage size={20} />
              </button>
              <button 
                onClick={() => fileInputRef.current?.click()}
                className="p-2 rounded-full hover:bg-[--primary-theme]/10 text-[--primary-theme] transition-colors hidden sm:block"
                title="Video"
              >
                <Clapperboard size={20} />
              </button>
              <input
                ref={fileInputRef}
                type="file"
                accept="image/jpeg,image/png,image/gif,video/mp4"
                onChange={handleMediaSelect}
                hidden
              />
            </div>
            
            <button
              onClick={handleCreatePost}
              disabled={posting || (!postText.trim() && !mediaFile)}
              className="bg-[--primary-theme] hover:bg-[#129c94] text-white font-bold py-2 px-6 rounded-full transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {posting ? "Posting..." : "Post"}
            </button>
          </div>
        </div>
      </div>

      <AllPostUI refreshKey={refreshKey} />
    </div>
  );
}

export function AllPostUI({ refreshKey }: { refreshKey: number }) {
  const [posts, setPosts] = useState<FeedPost[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchFeed = async () => {
      try {
        const res = await Api.get<FeedPost[]>("/posts/feed");
        setPosts(res.data ?? []);
      } catch (err) {
        console.error("Failed to load feed:", err);
        setPosts([]);
      } finally {
        setLoading(false);
      }
    };

    fetchFeed();
  }, [refreshKey]);

  if (loading) {
    return <div className="p-8 text-center text-gray-500 font-medium animate-pulse">Loading feed...</div>;
  }

  if (posts.length === 0) {
    return (
      <div className="p-8 text-center border-b border-white/10">
        <h2 className="text-xl font-extrabold mb-2 text-white">Welcome to your timeline</h2>
        <p className="text-gray-500 text-sm">When you follow people, you'll see the posts they share here.</p>
      </div>
    );
  }

  return (
    <div className="flex flex-col w-full">
      {posts.map((post) => (
        <UserPostContent key={post.id} post={post} />
      ))}
    </div>
  );
}

export function UserPostContent({ post }: { post: FeedPost }) {
  const avatarUrl = post.author_avatar
    ? post.author_avatar.startsWith("http")
      ? post.author_avatar
      : `http://localhost:8080/${post.author_avatar}`
    : undefined;

  return (
    <article className="flex gap-3 p-4 border-b border-white/10 hover:bg-white/[0.02] transition-colors cursor-pointer">
      <div className="shrink-0">
        <UserProfileImage url={avatarUrl} name={post.author_name} />
      </div>
      <div className="flex flex-col w-full min-w-0">
        {/* Header Info */}
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-1.5 overflow-hidden">
            <span className="font-bold text-white truncate hover:underline">{post.author_name}</span>
            <span className="text-sm text-gray-500 shrink-0">·</span>
            <span className="text-sm text-gray-500 shrink-0 hover:underline">
               {new Date(post.created_at).toLocaleDateString(undefined, { month: 'short', day: 'numeric' })}
            </span>
          </div>
          
          <button className="text-gray-500 hover:text-[--primary-theme] hover:bg-[--primary-theme]/10 p-1.5 rounded-full transition-colors shrink-0">
            <MoreHorizontal size={18} />
          </button>
        </div>

        {/* Content */}
        <div className="mt-1 text-[15px] leading-relaxed text-white whitespace-pre-wrap break-words">
          {post.content}
        </div>

        {/* Media */}
        {post.media_path && (
          <div className="mt-3 relative w-full rounded-2xl overflow-hidden border border-white/10">
            <Image
              src={`http://localhost:8080/${post.media_path}`}
              width={600}
              height={400}
              className="w-full h-auto object-cover max-h-[500px]"
              alt="Post attachment"
            />
          </div>
        )}

        {/* Interactions */}
        <div className="mt-3">
          <PostInteractions 
            postId={post.id} 
            comments={post.comment_count} 
            likes={0} 
            postDetails={post} 
          />
        </div>
      </div>
    </article>
  );
}