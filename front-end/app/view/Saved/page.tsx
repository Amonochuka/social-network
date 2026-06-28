"use client";

import { useEffect, useState } from "react";
import DefaultLayout from "@/components/layouts/defaultLayout";
import UserProfileImage from "@/components/header/profile/userProfile";
import { BookMarkedIcon, Heart, MessageCircleReply, Eye, Lock as LockIcon, Trash2 } from "lucide-react";
import Image from "next/image";

interface SavedPost {
  id: string;
  author_name: string;
  author_avatar?: string;
  privacy: string;
  created_at: string;
  content: string;
  media_path?: string;
  comment_count: number;
}

export default function SavedPostsPage() {
  const [savedPosts, setSavedPosts] = useState<SavedPost[]>([]);

  useEffect(() => {
    try {
      const posts = JSON.parse(localStorage.getItem("saved_posts") || "[]");
      setSavedPosts(posts);
    } catch {
      setSavedPosts([]);
    }
  }, []);

  const handleRemove = (postId: string) => {
    try {
      const posts = JSON.parse(localStorage.getItem("saved_posts") || "[]");
      const updated = posts.filter((p: any) => p.id !== postId);
      localStorage.setItem("saved_posts", JSON.stringify(updated));
      setSavedPosts(updated);
    } catch (err) {
      console.error(err);
    }
  };

  return (
    <DefaultLayout>
      <main className="h-full w-full px-8 py-6 bg-[#181818] overflow-y-auto">
        <div className="mb-6">
          <h1 className="text-3xl font-bold text-white flex items-center gap-2">
            <BookMarkedIcon className="text-[--primary-theme]" /> Saved Posts
          </h1>
          <p className="mt-1 text-sm text-gray-400">
            Your bookmarked posts for quick reading and reference.
          </p>
        </div>

        {savedPosts.length === 0 ? (
          <div className="flex flex-col items-center justify-center py-20 rounded-2xl border border-white/5 bg-[#1e1e1e]">
            <BookMarkedIcon size={48} className="text-gray-600 mb-4" />
            <p className="text-gray-400 text-lg font-medium">No saved posts yet</p>
            <p className="text-gray-600 text-sm mt-1">Bookmark posts in the main feed to see them here.</p>
          </div>
        ) : (
          <div className="flex flex-col gap-6 max-w-2xl">
            {savedPosts.map((post) => {
              const avatarUrl = post.author_avatar
                ? post.author_avatar.startsWith("http")
                  ? post.author_avatar
                  : `http://localhost:8080/${post.author_avatar}`
                : undefined;

              const mediaUrl = post.media_path
                ? post.media_path.startsWith("http")
                  ? post.media_path
                  : `http://localhost:8080/${post.media_path}`
                : undefined;

              return (
                <div key={post.id} className="rounded-2xl border border-white/10 bg-[#1e1e1e] p-5 shadow-lg relative group">
                  {/* Remove button */}
                  <button
                    onClick={() => handleRemove(post.id)}
                    className="absolute top-5 right-5 p-2 rounded-xl text-gray-500 hover:text-red-400 hover:bg-white/5 opacity-0 group-hover:opacity-100 transition-opacity"
                    title="Remove from saved"
                  >
                    <Trash2 size={18} />
                  </button>

                  {/* Header */}
                  <div className="flex items-center gap-3 mb-4">
                    <UserProfileImage url={avatarUrl} name={post.author_name} />
                    <div>
                      <p className="text-sm font-semibold text-white">{post.author_name}</p>
                      <div className="flex items-center gap-2 mt-0.5 text-xs text-gray-500">
                        <span>{new Date(post.created_at).toLocaleDateString()}</span>
                        <span>•</span>
                        <div className="flex items-center gap-1">
                          {post.privacy === "public" ? <Eye size={12} /> : <LockIcon size={12} />}
                          <span className="capitalize">{post.privacy}</span>
                        </div>
                      </div>
                    </div>
                  </div>

                  {/* Content */}
                  <p className="text-sm text-gray-300 leading-relaxed whitespace-pre-wrap mb-4">
                    {post.content}
                  </p>

                  {mediaUrl && (
                    <div className="relative w-full h-64 rounded-xl overflow-hidden mb-4 border border-white/5 bg-[#252525]">
                      <Image
                        src={mediaUrl}
                        fill
                        sizes="(max-width: 768px) 100vw, 50vw"
                        style={{ objectFit: "cover" }}
                        alt="post media"
                      />
                    </div>
                  )}

                  {/* Footer interactions */}
                  <div className="flex items-center gap-4 text-xs text-gray-500 pt-3 border-t border-white/5">
                    <div className="flex items-center gap-1.5">
                      <Heart size={16} />
                      <span>Likes</span>
                    </div>
                    <div className="flex items-center gap-1.5">
                      <MessageCircleReply size={16} />
                      <span>{post.comment_count} comments</span>
                    </div>
                  </div>
                </div>
              );
            })}
          </div>
        )}
      </main>
    </DefaultLayout>
  );
}