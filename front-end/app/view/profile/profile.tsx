"use client";

import { useEffect, useState } from "react";
import {
  MoreHorizontal,
  Pencil,
  Calendar,
  X,
  Camera,
} from "lucide-react";
import { useAppSelector, useAppDispatch } from "@/store/hooks";
import { authSelector, setSession } from "@/store/features/authSlice";
import { Api } from "@/services/axios";
import Image from "next/image";
import Link from "next/link";

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
}

interface FollowerProfile {
  user_id: string;
  first_name: string;
  last_name: string;
  avatar: string;
  nickname: string;
}

export default function ProfilePage() {
  const { user } = useAppSelector(authSelector);
  const dispatch = useAppDispatch();

  const [isPublic, setIsPublic] = useState(user?.is_public ?? true);
  const [updating, setUpdating] = useState(false);
  const [showEditModal, setShowEditModal] = useState(false);

  // stats
  const [followers, setFollowers] = useState<FollowerProfile[]>([]);
  const [following, setFollowing] = useState<FollowerProfile[]>([]);
  const [posts, setPosts] = useState<FeedPost[]>([]);
  const [statsLoading, setStatsLoading] = useState(true);

  useEffect(() => {
    if (!user) return;
    const fetchStats = async () => {
      try {
        const [followersRes, followingRes, postsRes] = await Promise.all([
          Api.get<FollowerProfile[]>("/followers"),
          Api.get<FollowerProfile[]>("/following"),
          Api.get<FeedPost[]>(`/users/${user.id}/posts`),
        ]);
        setFollowers(followersRes.data ?? []);
        setFollowing(followingRes.data ?? []);
        setPosts(postsRes.data ?? []);
      } catch (err) {
        console.error("Failed to load profile stats:", err);
      } finally {
        setStatsLoading(false);
      }
    };
    fetchStats();
  }, [user]);

  if (!user) {
    return (
      <div className="flex h-screen w-full items-center justify-center bg-[#181818] text-gray-400">
        Loading profile...
      </div>
    );
  }

  const handlePrivacyToggle = async () => {
    if (updating) return;
    setUpdating(true);

    const newValue = !isPublic;
    try {
      await Api.put("/profile/privacy", { is_public: newValue });
      setIsPublic(newValue);
    } catch (err) {
      console.error("Failed to update privacy:", err);
    } finally {
      setUpdating(false);
    }
  };

  const initials = `${user.first_name?.[0] ?? ""}${user.last_name?.[0] ?? ""}`.toUpperCase();

  return (
    <div className="flex h-screen w-full flex-col overflow-y-auto bg-[#181818] pr-20">

      {/* ── cover ── */}
      <div className="relative h-60 w-full shrink-0 overflow-hidden sm:h-70">
        <div className="absolute inset-0 bg-gradient-to-br from-slate-900 via-purple-950 to-indigo-900">
          <div className="absolute inset-0 bg-[radial-gradient(ellipse_at_60%_40%,rgba(139,92,252,0.35),transparent_70%)]" />
        </div>

        <div className="absolute right-4 top-5 flex items-center gap-2">
          <button
            onClick={() => setShowEditModal(true)}
            className="flex items-center gap-2 rounded-xl border border-white/20 bg-black/30 px-4 py-2 text-sm font-semibold text-white backdrop-blur-sm transition hover:bg-black/50"
          >
            <Pencil size={14} />
            Edit Profile
          </button>
          <button className="grid h-9 w-9 place-items-center rounded-xl border border-white/20 bg-black/30 text-white backdrop-blur-sm transition hover:bg-black/50">
            <MoreHorizontal size={17} />
          </button>

          {/* private profile toggle */}
          <div className="flex items-center gap-2 rounded-xl border border-white/20 bg-black/30 px-3 py-1.5 backdrop-blur-sm">
            <span className="text-xs font-medium text-white/80">
              {isPublic ? "Public Profile" : "Private Profile"}
            </span>
            <button
              role="switch"
              aria-checked={!isPublic}
              disabled={updating}
              onClick={handlePrivacyToggle}
              className={`relative h-5 w-9 shrink-0 rounded-full transition-colors duration-200 ${
                isPublic ? "bg-white/20" : "bg-[--primary-theme]"
              } ${updating ? "opacity-50" : ""}`}
            >
              <span
                className={`absolute top-0.5 left-0.5 h-4 w-4 rounded-full bg-white shadow transition-transform duration-200 ${
                  isPublic ? "translate-x-0" : "translate-x-4"
                }`}
              />
            </button>
          </div>
        </div>
      </div>

      {/* ── avatar + identity ── */}
      <div className="w-full bg-[#1e1e1e] px-6 pb-5 pt-0">
        <div className="-mt-14 mb-4 flex items-end justify-between">
          <div className="relative h-28 w-28 shrink-0 rounded-full border-4 border-[#1e1e1e] bg-gradient-to-br from-violet-500 to-fuchsia-500 shadow-xl overflow-hidden">
            {user.avatar ? (
              <img
                src={user.avatar.startsWith("http") ? user.avatar : `http://localhost:8080/${user.avatar}`}
                alt={`${user.first_name} ${user.last_name}`}
                className="h-full w-full object-cover"
              />
            ) : (
              <div className="flex h-full w-full items-center justify-center text-3xl font-bold text-white">
                {initials}
              </div>
            )}
            <span className="absolute bottom-1.5 right-1.5 h-4 w-4 rounded-full border-2 border-[#1e1e1e] bg-emerald-400" />
          </div>
        </div>

        <div className="flex items-baseline gap-2.5">
          <h1 className="text-2xl font-bold text-white">
            {user.first_name} {user.last_name}
          </h1>
          {user.nickname && (
            <span className="text-sm text-gray-500">@{user.nickname}</span>
          )}
        </div>

        <p className="mt-2 text-sm leading-relaxed text-gray-400">
          {user.about_me || "No bio yet."}
        </p>

        {user.date_of_birth && user.date_of_birth !== "" && (
          <div className="mt-3 flex items-center gap-2 text-sm text-gray-500">
            <Calendar size={14} />
            Born {new Date(user.date_of_birth).toLocaleDateString()}
          </div>
        )}

        {/* Stats row */}
        <div className="mt-5 flex gap-6">
          <div className="flex flex-col items-center">
            <span className="text-xl font-bold text-white">{posts.length}</span>
            <span className="text-xs text-gray-500">Posts</span>
          </div>
          <Link href="/view/Network" className="flex flex-col items-center hover:opacity-80 transition">
            <span className="text-xl font-bold text-white">{followers.length}</span>
            <span className="text-xs text-gray-500">Followers</span>
          </Link>
          <Link href="/view/Network" className="flex flex-col items-center hover:opacity-80 transition">
            <span className="text-xl font-bold text-white">{following.length}</span>
            <span className="text-xs text-gray-500">Following</span>
          </Link>
        </div>
      </div>

      {/* ── own posts ── */}
      <div className="flex-1 p-6">
        {statsLoading ? (
          <div className="rounded-2xl bg-[#222] p-8 text-center text-sm text-gray-500">
            Loading posts...
          </div>
        ) : posts.length === 0 ? (
          <div className="rounded-2xl bg-[#222] p-8 text-center text-sm text-gray-500">
            No posts yet. Create your first post from the Home page.
          </div>
        ) : (
          <div className="flex flex-col gap-4">
            {posts.map((post) => (
              <div key={post.id} className="rounded-xl bg-[#222] p-5 border border-white/5">
                <p className="text-sm text-gray-300 leading-relaxed">{post.content}</p>
                {post.media_path && (
                  <div className="mt-3 relative w-full h-48 rounded-lg overflow-hidden">
                    <Image
                      src={`http://localhost:8080/${post.media_path}`}
                      fill
                      sizes="(max-width: 768px) 100vw, 50vw"
                      style={{ objectFit: "cover" }}
                      alt="post media"
                    />
                  </div>
                )}
                <div className="mt-3 flex items-center gap-4 text-xs text-gray-500">
                  <span>{new Date(post.created_at).toLocaleDateString()}</span>
                  <span>{post.comment_count} comments</span>
                  <span className="capitalize">{post.privacy}</span>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>

      {showEditModal && (
        <EditProfileModal
          onClose={() => setShowEditModal(false)}
          onSaved={(updatedUser) => {
            dispatch(setSession(updatedUser));
            setShowEditModal(false);
          }}
        />
      )}
    </div>
  );
}

interface EditProfileModalProps {
  onClose: () => void;
  onSaved: (user: any) => void;
}

function EditProfileModal({ onClose, onSaved }: EditProfileModalProps) {
  const { user } = useAppSelector(authSelector);

  const [firstName, setFirstName] = useState(user?.first_name ?? "");
  const [lastName, setLastName] = useState(user?.last_name ?? "");
  const [nickname, setNickname] = useState(user?.nickname ?? "");
  const [aboutMe, setAboutMe] = useState(user?.about_me ?? "");
  const [dateOfBirth, setDateOfBirth] = useState(
    user?.date_of_birth ? user.date_of_birth.split("T")[0] : ""
  );
  const [avatarFile, setAvatarFile] = useState<File | null>(null);
  const [avatarPreview, setAvatarPreview] = useState<string | null>(null);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleAvatarChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;
    setAvatarFile(file);
    setAvatarPreview(URL.createObjectURL(file));
  };

  const handleSave = async () => {
    setSaving(true);
    setError(null);

    try {
      // 1. update text fields
      const res = await Api.put("/profile", {
        first_name: firstName,
        last_name: lastName,
        nickname: nickname,
        about_me: aboutMe,
        date_of_birth: dateOfBirth,
      });

      let updatedUser = res.data;

      // 2. upload avatar separately if a new one was selected
      if (avatarFile) {
        const formData = new FormData();
        formData.append("avatar", avatarFile);
        const avatarRes = await Api.post("/profile/avatar", formData, {
          headers: { "Content-Type": "multipart/form-data" },
        });
        updatedUser = { ...updatedUser, avatar: avatarRes.data.avatar };
      }

      onSaved(updatedUser);
    } catch (err) {
      console.error("Failed to update profile:", err);
      setError("Failed to save changes. Please try again.");
    } finally {
      setSaving(false);
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 p-4">
      <div className="w-full max-w-md rounded-2xl bg-[#1e1e1e] p-6 shadow-2xl">
        <div className="mb-5 flex items-center justify-between">
          <h2 className="text-lg font-bold text-white">Edit Profile</h2>
          <button
            onClick={onClose}
            className="rounded-lg p-1.5 text-gray-400 transition hover:bg-white/5 hover:text-white"
          >
            <X size={18} />
          </button>
        </div>

        {error && (
          <div className="mb-4 rounded-lg bg-red-500/10 px-3 py-2 text-sm text-red-400">
            {error}
          </div>
        )}

        {/* avatar */}
        <div className="mb-5 flex items-center gap-4">
          <div className="relative h-20 w-20 shrink-0 overflow-hidden rounded-full bg-gradient-to-br from-violet-500 to-fuchsia-500">
            {avatarPreview || user?.avatar ? (
              <img
                src={
                  avatarPreview
                    ? avatarPreview
                    : user!.avatar.startsWith("http")
                    ? user!.avatar
                    : `http://localhost:8080/${user!.avatar}`
                }
                alt="avatar preview"
                className="h-full w-full object-cover"
              />
            ) : (
              <div className="flex h-full w-full items-center justify-center text-xl font-bold text-white">
                {firstName[0]}
                {lastName[0]}
              </div>
            )}
          </div>
          <label className="flex cursor-pointer items-center gap-2 rounded-lg border border-white/15 px-3 py-2 text-xs font-semibold text-white transition hover:bg-white/5">
            <Camera size={14} />
            Change photo
            <input
              type="file"
              accept="image/jpeg,image/png,image/gif"
              onChange={handleAvatarChange}
              hidden
            />
          </label>
        </div>

        {/* fields */}
        <div className="space-y-3">
          <div className="grid grid-cols-2 gap-3">
            <div>
              <label className="mb-1 block text-xs text-gray-400">First name</label>
              <input
                value={firstName}
                onChange={(e) => setFirstName(e.target.value)}
                className="w-full rounded-lg bg-[#262626] px-3 py-2 text-sm text-white outline-none focus:ring-1 focus:ring-[--primary-theme]"
              />
            </div>
            <div>
              <label className="mb-1 block text-xs text-gray-400">Last name</label>
              <input
                value={lastName}
                onChange={(e) => setLastName(e.target.value)}
                className="w-full rounded-lg bg-[#262626] px-3 py-2 text-sm text-white outline-none focus:ring-1 focus:ring-[--primary-theme]"
              />
            </div>
          </div>

          <div>
            <label className="mb-1 block text-xs text-gray-400">Nickname</label>
            <input
              value={nickname}
              onChange={(e) => setNickname(e.target.value)}
              className="w-full rounded-lg bg-[#262626] px-3 py-2 text-sm text-white outline-none focus:ring-1 focus:ring-[--primary-theme]"
            />
          </div>

          <div>
            <label className="mb-1 block text-xs text-gray-400">Date of birth</label>
            <input
              type="date"
              value={dateOfBirth}
              onChange={(e) => setDateOfBirth(e.target.value)}
              className="w-full rounded-lg bg-[#262626] px-3 py-2 text-sm text-white outline-none focus:ring-1 focus:ring-[--primary-theme]"
            />
          </div>

          <div>
            <label className="mb-1 block text-xs text-gray-400">About me</label>
            <textarea
              value={aboutMe}
              onChange={(e) => setAboutMe(e.target.value)}
              rows={3}
              className="w-full resize-none rounded-lg bg-[#262626] px-3 py-2 text-sm text-white outline-none focus:ring-1 focus:ring-[--primary-theme]"
            />
          </div>
        </div>

        <div className="mt-6 flex justify-end gap-2">
          <button
            onClick={onClose}
            className="rounded-lg px-4 py-2 text-sm font-semibold text-gray-300 transition hover:bg-white/5"
          >
            Cancel
          </button>
          <button
            onClick={handleSave}
            disabled={saving}
            className="rounded-lg bg-[--primary-theme] px-4 py-2 text-sm font-semibold text-black transition hover:opacity-90 disabled:opacity-50"
          >
            {saving ? "Saving..." : "Save changes"}
          </button>
        </div>
      </div>
    </div>
  );
}