"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import { Api } from "@/services/axios";
import { useAppSelector, useAppDispatch } from "@/store/hooks";
import { authSelector } from "@/store/features/authSlice";
import DefaultLayout from "@/components/layouts/defaultLayout";
import { Calendar, Lock, UserPlus, UserMinus } from "lucide-react";
import Image from "next/image";
import {sendFollowRequest,unfollowUser,} from "@/store/features/followerSlice";

interface UserProfile {
  id: string;
  email: string;
  first_name: string;
  last_name: string;
  date_of_birth: string;
  avatar: string;
  nickname: string;
  about_me: string;
  is_public: boolean;
}

interface FeedPost {
  id: string;
  content: string;
  media_path: string;
  media_type: string;
  privacy: string;
  comment_count: number;
  created_at: string;
}


export default function OtherUserProfilePage() {
  return (
    <DefaultLayout>
      <OtherUserProfile />
    </DefaultLayout>
  );
}

function OtherUserProfile() {
  const params = useParams();
  const userId = params.userId as string;
  const { user: currentUser } = useAppSelector(authSelector);
  const dispatch = useAppDispatch();

  const [profile, setProfile] = useState<UserProfile | null>(null);
  const [loading, setLoading] = useState(true);
  const [followStatus, setFollowStatus] = useState<"none" | "requested" | "following">("none");
  const [followLoading, setFollowLoading] = useState(false);
  const [posts, setPosts] = useState<FeedPost[]>([]);
  const [postsLoading, setPostsLoading] = useState(false);
  const [canViewProfile, setCanViewProfile] = useState(false);

  useEffect(() => {
    const fetchProfile = async () => {
      setLoading(true);
      try {
        const res = await Api.get<UserProfile>(`/profile/${userId}`);
        setProfile(res.data);

        // Check if we follow this user
        const statusRes = await Api.get<{ status: "none" | "requested" | "following" }>(`/follow/${userId}/status`);
        const status = statusRes.data.status;
        setFollowStatus(status);
       
        // Can view if public profile or following

        const canView =
          res.data.is_public ||
          status === "following" ||
          currentUser?.id === userId;

        setCanViewProfile(canView);

        // Load posts if allowed
        // Load posts if allowed
        if (canView) {
          setPostsLoading(true);

        try {
          const postsRes = await Api.get<FeedPost[]>(`/users/${userId}/posts`);
          setPosts(postsRes.data ?? []);
        } catch {
          setPosts([]);
        } finally {
          setPostsLoading(false);
        }
      } else {
          setPosts([]);
      }
      } catch (err) {
          console.error(err);
          setProfile(null);
        } finally {
        setLoading(false);
      }
    };

    if (userId) fetchProfile();
  }, [userId, currentUser?.id, dispatch]);

  const handleFollow = async () => {
    setFollowLoading(true);

    try {
      await dispatch(sendFollowRequest(userId)).unwrap();

      if (profile?.is_public) {
        setFollowStatus("following");
      } else {
        setFollowStatus("requested");
      }

      if (profile?.is_public) {
        setCanViewProfile(true);

        const postsRes = await Api.get<FeedPost[]>(`/users/${userId}/posts`);
        setPosts(postsRes.data ?? []);
      }
    } catch (err) {
      console.error("Failed to follow:", err);
    } finally {
      setFollowLoading(false);
    }
  };

  const handleUnfollow = async () => {
    setFollowLoading(true);

    try {
      await dispatch(unfollowUser(userId)).unwrap();

      setFollowStatus("none");

      if (!profile?.is_public) {
        setCanViewProfile(false);
        setPosts([]);
      }
    } catch (err) {
      console.error("Failed to unfollow:", err);
    } finally {
      setFollowLoading(false);
    }
  };

  if (loading) {
    return (
      <div className="flex h-screen w-full items-center justify-center text-gray-400">
        Loading profile...
      </div>
    );
  }
    if (!profile) {
      return (
       <div className="flex h-screen w-full items-center justify-center text-gray-400">
      User not found.
    </div>
      );
    }

  const initials = `${profile.first_name?.[0] ?? ""}${profile.last_name?.[0] ?? ""}`.toUpperCase();
  const avatarUrl = profile.avatar
    ? profile.avatar.startsWith("http")
      ? profile.avatar
      : `http://localhost:8080/${profile.avatar}`
    : null;

  const isOwnProfile = currentUser?.id === userId;

  return (
    <div className="flex h-screen w-full flex-col overflow-y-auto bg-[#181818]">
      {/* Cover */}
      <div className="relative h-52 w-full shrink-0 overflow-hidden">
        <div className="absolute inset-0 bg-gradient-to-br from-slate-900 via-indigo-950 to-purple-900">
          <div className="absolute inset-0 bg-[radial-gradient(ellipse_at_50%_50%,rgba(99,102,241,0.3),transparent_70%)]" />
        </div>
      </div>

      {/* Profile info */}
      <div className="w-full bg-[#1e1e1e] px-6 pb-5">
        <div className="-mt-14 mb-4 flex items-end justify-between">
          <div className="relative h-28 w-28 shrink-0 rounded-full border-4 border-[#1e1e1e] bg-gradient-to-br from-violet-500 to-fuchsia-500 shadow-xl overflow-hidden">
            {avatarUrl ? (
              <img src={avatarUrl} alt={`${profile.first_name}`} className="h-full w-full object-cover" />
            ) : (
              <div className="flex h-full w-full items-center justify-center text-3xl font-bold text-white">
                {initials}
              </div>
            )}
          </div>

          {/* Follow/Unfollow button */}
          {!isOwnProfile && (
            <button
              onClick={followStatus === "following"? handleUnfollow: handleFollow}
              disabled={followLoading || followStatus === "requested"}
             className={`flex items-center gap-2 rounded-xl px-5 py-2.5 text-sm font-semibold transition ${
  followStatus === "following"
    ? "border border-white/15 bg-transparent text-gray-300 hover:bg-white/5 hover:text-red-400"
    : followStatus === "requested"
    ? "border border-yellow-500/30 bg-yellow-500/10 text-yellow-300"
    : "bg-[--primary-theme] text-black hover:opacity-90"
} disabled:opacity-50`}
            >
              {followStatus === "following" ? (
    <>
        <UserMinus size={16} />
        {followLoading ? "..." : "Unfollow"}
    </>
) : followStatus === "requested" ? (
    <>
        <UserPlus size={16} />
        Requested
    </>
) : (
    <>
        <UserPlus size={16} />
        {followLoading ? "..." : "Follow"}
    </>
)}
            </button>
          )}
        </div>

        <div className="flex items-baseline gap-2.5">
          <h1 className="text-2xl font-bold text-white">
            {profile.first_name} {profile.last_name}
          </h1>
          {profile.nickname && (
            <span className="text-sm text-gray-500">@{profile.nickname}</span>
          )}
        </div>

        {/* Private badge */}
        {!profile.is_public && (
        <div className="mt-3">
          <span className="inline-flex items-center gap-2 rounded-full border border-white/10 bg-white/5 px-3 py-1 text-xs text-gray-400">
            <Lock size={12} />
            Private Account
            </span>
          </div>
        )}

        {canViewProfile && (
          <>
            <p className="mt-2 text-sm leading-relaxed text-gray-400">
              {profile.about_me || "No bio yet."}
            </p>

            {profile.date_of_birth && (
              <div className="mt-3 flex items-center gap-2 text-sm text-gray-500">
                <Calendar size={14} />
                Born {new Date(profile.date_of_birth).toLocaleDateString()}
              </div>
            )}
          </>
        )}

        
      </div>

      {/* Posts */}
<div className="flex-1 p-6">
  {canViewProfile ? (
    postsLoading ? (
      <div className="text-center text-sm text-gray-500 py-8">
        Loading posts...
      </div>
    ) : posts.length === 0 ? (
      <div className="rounded-2xl bg-[#222] p-8 text-center text-sm text-gray-500">
        No posts yet.
      </div>
    ) : (
      <div className="flex flex-col gap-4">
        {posts.map((post) => (
          <div
            key={post.id}
            className="rounded-xl bg-[#222] p-5 border border-white/5"
          >
            <p className="text-sm text-gray-300 leading-relaxed">
              {post.content}
            </p>

            {post.media_path && (
              <div className="mt-3 relative w-full h-48 rounded-lg overflow-hidden">
                <Image
                  src={`http://localhost:8080/${post.media_path}`}
                  fill
                  sizes="(max-width:768px) 100vw, 50vw"
                  style={{ objectFit: "cover" }}
                  alt="post media"
                />
              </div>
            )}

            <div className="mt-3 flex items-center gap-4 text-xs text-gray-500">
              <span>
                {new Date(post.created_at).toLocaleDateString()}
              </span>

              <span>{post.comment_count} comments</span>
            </div>
          </div>
        ))}
      </div>
    )
  ) : (
    <div className="flex h-72 flex-col items-center justify-center rounded-2xl border border-white/10 bg-[#202020] text-center">
      <div className="mb-4 rounded-full bg-white/5 p-5">
        <Lock size={32} className="text-gray-500" />
      </div>

      <h2 className="text-xl font-semibold text-white">
        Private Account
      </h2>

      <p className="mt-2 max-w-sm text-sm text-gray-400">
        Follow this account to see their posts, photos and profile
        information.
      </p>
    </div>
  )}
</div>
       
      )
    </div>
  );
}
