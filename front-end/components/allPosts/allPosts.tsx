"use client";

import { useEffect, useState, useRef } from "react";
import UserProfileImage from "../header/profile/userProfile";
import style from "@/styles/all-post.module.css";
import { Button } from "../ui/button";
import { ButtonData } from "@/types";
import { createPostBtn } from "@/styles/style";
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
} from "lucide-react";
import Image from "next/image";
import { CSSProperties } from "react";
import { PostInteractions } from "./interactions";
import { Api } from "@/services/axios";
import { useAppSelector } from "@/store/hooks";
import { authSelector } from "@/store/features/authSlice";

const data: ButtonData = {
  text: "create post",
  type: "button",
  style: createPostBtn,
};

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

  // load followers for private post picker
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
    if (!postText.trim() || posting) return;
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
    public: { icon: Globe, label: "Public" },
    almost_private: { icon: Users, label: "Followers Only" },
    private: { icon: Lock, label: "Selected Followers" },
  };

  const PrivacyIcon = privacyLabels[privacy].icon;

  return (
    <div className={style.postParentCont}>
      <div className={style.postHomeCont}>
        <div className={style.postHomeMain}>
          <UserProfileImage url={myAvatarUrl} name={myFullName} />
          <input
            className={style.postData}
            type="text"
            placeholder="what's on your mind"
            name="postData"
            value={postText}
            onChange={(e) => setPostText(e.target.value)}
            onKeyDown={(e) => e.key === "Enter" && handleCreatePost()}
          />
          <Button
            data={{
              ...data,
              text: posting ? "posting..." : "create post",
              onClick: handleCreatePost,
            }}
          />
        </div>

        {/* Media preview */}
        {mediaPreview && (
          <div style={{ position: "relative", margin: "0.75rem 0.5rem", maxWidth: "200px" }}>
            <img
              src={mediaPreview}
              alt="preview"
              style={{ width: "100%", borderRadius: "0.5rem", objectFit: "cover", maxHeight: "150px" }}
            />
            <button
              onClick={removeMedia}
              style={{
                position: "absolute", top: 4, right: 4,
                background: "rgba(0,0,0,0.7)", border: "none",
                borderRadius: "50%", padding: "3px", cursor: "pointer", color: "white",
              }}
            >
              <X size={14} />
            </button>
          </div>
        )}

        {/* Icons row */}
        <div className={style.iconPostDisplay}>
          <div
            className={style.iconFlex}
            onClick={() => fileInputRef.current?.click()}
            style={{ cursor: "pointer" }}
          >
            <LucideImage />
            <span>photos</span>
          </div>
          <div
            className={style.iconFlex}
            onClick={() => fileInputRef.current?.click()}
            style={{ cursor: "pointer" }}
          >
            <Clapperboard />
            <span>videos</span>
          </div>

          {/* Privacy selector */}
          <div style={{ position: "relative", marginLeft: "auto" }}>
            <button
              onClick={() => setShowPrivacyMenu((p) => !p)}
              style={{
                display: "flex", alignItems: "center", gap: "0.4rem",
                background: "rgba(255,255,255,0.06)", border: "1px solid rgba(255,255,255,0.1)",
                borderRadius: "0.5rem", padding: "0.4rem 0.75rem",
                color: "var(--primary-theme)", fontSize: "0.8rem", fontWeight: 600,
                cursor: "pointer",
              }}
            >
              <PrivacyIcon size={14} />
              {privacyLabels[privacy].label}
            </button>

            {showPrivacyMenu && (
              <div
                style={{
                  position: "absolute", top: "100%", right: 0, marginTop: "0.25rem",
                  background: "#2a2a2a", border: "1px solid rgba(255,255,255,0.1)",
                  borderRadius: "0.5rem", overflow: "hidden", zIndex: 50, minWidth: "180px",
                }}
              >
                {(["public", "almost_private", "private"] as const).map((p) => {
                  const Icon = privacyLabels[p].icon;
                  return (
                    <button
                      key={p}
                      onClick={() => { setPrivacy(p); setShowPrivacyMenu(false); }}
                      style={{
                        display: "flex", alignItems: "center", gap: "0.5rem",
                        width: "100%", padding: "0.6rem 0.75rem",
                        background: privacy === p ? "rgba(255,255,255,0.08)" : "transparent",
                        border: "none", color: "#e5e7eb", fontSize: "0.85rem",
                        cursor: "pointer", textAlign: "left",
                      }}
                    >
                      <Icon size={14} />
                      {privacyLabels[p].label}
                      {privacy === p && <Check size={14} style={{ marginLeft: "auto", color: "var(--primary-theme)" }} />}
                    </button>
                  );
                })}
              </div>
            )}
          </div>
        </div>

        {/* Follower picker for private posts */}
        {showFollowerPicker && (
          <div style={{
            margin: "0.75rem 0.5rem", padding: "0.75rem",
            background: "rgba(255,255,255,0.04)", borderRadius: "0.5rem",
            border: "1px solid rgba(255,255,255,0.08)",
          }}>
            <p style={{ fontSize: "0.8rem", color: "#9ca3af", marginBottom: "0.5rem" }}>
              Select followers who can see this post:
            </p>
            <div style={{ display: "flex", flexWrap: "wrap", gap: "0.4rem" }}>
              {followers.length === 0 ? (
                <p style={{ fontSize: "0.75rem", color: "#666" }}>No followers yet.</p>
              ) : (
                followers.map((f) => {
                  const selected = selectedFollowers.includes(f.user_id);
                  return (
                    <button
                      key={f.user_id}
                      onClick={() => toggleFollower(f.user_id)}
                      style={{
                        display: "flex", alignItems: "center", gap: "0.35rem",
                        padding: "0.3rem 0.6rem", borderRadius: "1rem",
                        background: selected ? "var(--primary-theme)" : "rgba(255,255,255,0.08)",
                        color: selected ? "#000" : "#e5e7eb",
                        border: "none", fontSize: "0.78rem", fontWeight: 500,
                        cursor: "pointer", transition: "all 0.15s",
                      }}
                    >
                      {f.first_name} {f.last_name}
                      {selected && <Check size={12} />}
                    </button>
                  );
                })
              )}
            </div>
          </div>
        )}

        {/* Hidden file input */}
        <input
          ref={fileInputRef}
          type="file"
          accept="image/jpeg,image/png,image/gif,video/mp4"
          onChange={handleMediaSelect}
          hidden
        />
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
    return <div className={style.AllPostLayout}>Loading feed...</div>;
  }

  if (posts.length === 0) {
    return (
      <div className={style.AllPostLayout}>
        <p style={{ textAlign: "center", color: "#888", padding: "2rem" }}>
          No posts yet. Follow people to see their posts here.
        </p>
      </div>
    );
  }

  return (
    <div className={style.AllPostLayout}>
      {posts.map((post) => (
        <div key={post.id} className={style.userPostDisplay}>
          <UserPostProfile
            userImage={
              post.author_avatar
                ? post.author_avatar.startsWith("http")
                  ? post.author_avatar
                  : `http://localhost:8080/${post.author_avatar}`
                : undefined
            }
            fullName={post.author_name}
            status={post.privacy}
            datePosted={new Date(post.created_at).toLocaleDateString()}
            privacy={post.privacy}
          />
          <UserPostContent
            postId={post.id}
            description={post.content}
            postImage={
              post.media_path
                ? `http://localhost:8080/${post.media_path}`
                : undefined
            }
            likes={0}
            comments={post.comment_count}
          />
        </div>
      ))}
    </div>
  );
}

interface Props {
  userImage?: string;
  fullName: string;
  status: string;
  datePosted: string;
  privacy: string;
}

export function UserPostProfile({
  userImage,
  fullName,
  status,
  datePosted,
  privacy,
}: Props) {
  return (
    <div className={style.postUserProfile}>
      <div className={style.userImageName}>
        <UserProfileImage url={userImage} name={fullName} />
        <span className={style.userPostProfile}>
          <p>{fullName}</p>
          <p>posted on {datePosted}</p>
        </span>
      </div>
      <div className={style.userPrivacy}>
        {privacy === "public" ? <Eye /> : <LockIcon />}
        {status}
      </div>
    </div>
  );
}

interface ContentInterface {
  postId: string;
  description: string;
  postImage?: string;
  likes: number;
  comments: number;
}

export function UserPostContent({ postId, description, postImage, likes, comments }: ContentInterface) {
  const imageStyle: CSSProperties = {
    objectFit: "cover",
    borderRadius: "0.5rem",
  };

  return (
    <div className={style.postContMain}>
      <div className={style.postContMain}>
        {postImage ? (
          <>
            <div>
              <p className={style.postDesscription}>{description}</p>
            </div>

            <div className={style.postsImageCont}>
              <Image
                src={postImage}
                fill
                sizes="(max-width: 768px) 100vw, 50vw"
                style={imageStyle}
                alt="post image"
              />
            </div>
          </>
        ) : (
          <PostDescriptionUI description={description} />
        )}

        <div>
          <PostInteractions postId={postId} comments={comments} likes={likes} />
        </div>
      </div>
    </div>
  );
}

interface DescriptionProps {
  description: string;
}

export function PostDescriptionUI({ description }: DescriptionProps) {
  return (
    <div>
      <p className={style.shoutDescription}>{description}</p>
    </div>
  );
}