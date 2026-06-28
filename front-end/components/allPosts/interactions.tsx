"use client";

import { useEffect, useState } from "react";
import { BookMarkedIcon, Heart, MessageCircleReply } from "lucide-react";
import { Api } from "@/services/axios";
import UserProfileImage from "../header/profile/userProfile";

interface Props {
  likes: number;
  comments: number;
  postId: string;
  postDetails?: {
    id: string;
    author_name: string;
    author_avatar?: string;
    privacy: string;
    created_at: string;
    content: string;
    media_path?: string;
    comment_count: number;
  };
}

export function PostInteractions({ likes, comments, postId, postDetails }: Props) {
  const [showComments, setShowComments] = useState(false);
  const [isBookmarked, setIsBookmarked] = useState(false);

  useEffect(() => {
    try {
      const bookmarks = JSON.parse(localStorage.getItem("saved_posts") || "[]");
      setIsBookmarked(bookmarks.some((b: any) => b.id === postId));
    } catch {
      setIsBookmarked(false);
    }
  }, [postId]);

  const handleBookmark = () => {
    try {
      const bookmarks = JSON.parse(localStorage.getItem("saved_posts") || "[]");
      const index = bookmarks.findIndex((b: any) => b.id === postId);
      let updated = [];
      if (index > -1) {
        updated = bookmarks.filter((b: any) => b.id !== postId);
        setIsBookmarked(false);
      } else {
        const itemToSave = postDetails || {
          id: postId,
          content: "",
          author_name: "User",
          created_at: new Date().toISOString(),
          comment_count: comments,
          privacy: "public"
        };
        updated = [...bookmarks, itemToSave];
        setIsBookmarked(true);
      }
      localStorage.setItem("saved_posts", JSON.stringify(updated));
    } catch (err) {
      console.error("Failed to update bookmark:", err);
    }
  };

  return (
    <div className="w-full">
      <div className="flex items-center justify-between mt-1 text-gray-500 w-full max-w-md">
        <div className="flex items-center gap-12">
          {/* Like Button */}
          <button className="group flex items-center gap-1.5 transition-colors">
            <div className="p-2 rounded-full group-hover:bg-pink-500/10 group-hover:text-pink-500 transition-colors">
              <Heart size={18} />
            </div>
            <span className="text-sm font-medium group-hover:text-pink-500 transition-colors">{likes}</span>
          </button>
          
          {/* Comment Button */}
          <button 
            onClick={() => setShowComments((prev) => !prev)}
            className="group flex items-center gap-1.5 transition-colors"
          >
            <div className="p-2 rounded-full group-hover:bg-[--primary-theme]/10 group-hover:text-[--primary-theme] transition-colors">
              <MessageCircleReply size={18} />
            </div>
            <span className="text-sm font-medium group-hover:text-[--primary-theme] transition-colors">{comments}</span>
          </button>
        </div>

        {/* Bookmark Button */}
        <button 
          onClick={handleBookmark}
          className="group flex items-center transition-colors"
        >
          <div className="p-2 rounded-full group-hover:bg-blue-500/10 group-hover:text-blue-500 transition-colors">
            <BookMarkedIcon
              size={18}
              className={isBookmarked ? "text-blue-500" : "text-gray-500 group-hover:text-blue-500"}
              fill={isBookmarked ? "currentColor" : "none"}
            />
          </div>
        </button>
      </div>

      {showComments && <CommentsSectionUI postId={postId} />}
    </div>
  );
}

interface CommentDetail {
  id: string;
  post_id: string;
  user_id: string;
  author_name: string;
  author_avatar: string;
  content: string;
  media_path: string;
  media_type: string;
  created_at: string;
}

export function CommentsSectionUI({ postId }: { postId: string }) {
  const [comments, setComments] = useState<CommentDetail[]>([]);
  const [loading, setLoading] = useState(true);
  const [text, setText] = useState("");
  const [posting, setPosting] = useState(false);

  useEffect(() => {
    const fetchComments = async () => {
      try {
        const res = await Api.get<CommentDetail[]>(
          `/posts/${postId}/comments/all`
        );
        setComments(res.data ?? []);
      } catch (err) {
        console.error("Failed to load comments:", err);
        setComments([]);
      } finally {
        setLoading(false);
      }
    };

    fetchComments();
  }, [postId]);

  const handleAddComment = async () => {
    if (!text.trim() || posting) return;
    setPosting(true);

    try {
      const formData = new FormData();
      formData.append("content", text);

      const res = await Api.post(`/posts/${postId}/comments`, formData, {
        headers: { "Content-Type": "multipart/form-data" },
      });

      setComments((prev) => [...prev, res.data]);
      setText("");
    } catch (err) {
      console.error("Failed to add comment:", err);
    } finally {
      setPosting(false);
    }
  };

  return (
    <div className="mt-3 pt-3 border-t border-white/10 flex flex-col gap-4">
      {loading ? (
        <p className="text-sm text-gray-500 animate-pulse text-center">Loading comments...</p>
      ) : comments.length === 0 ? (
        <p className="text-sm text-gray-500 text-center">No comments yet.</p>
      ) : (
        <div className="flex flex-col gap-4">
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
              <div className="flex flex-col min-w-0 bg-white/5 rounded-2xl px-4 py-2 text-sm max-w-full">
                <p className="font-bold text-white truncate">{c.author_name}</p>
                <p className="text-white mt-0.5 break-words">{c.content}</p>
              </div>
            </div>
          ))}
        </div>
      )}

      {/* Add Comment Input */}
      <div className="flex items-center gap-3 mt-2">
        <input
          className="flex-1 bg-transparent border border-white/10 rounded-full px-4 py-2 text-sm text-white outline-none focus:border-[--primary-theme]"
          type="text"
          placeholder="Post your reply..."
          value={text}
          onChange={(e) => setText(e.target.value)}
          onKeyDown={(e) => e.key === "Enter" && handleAddComment()}
        />
        <button
          className="bg-[--primary-theme] hover:bg-[#129c94] text-white font-bold text-sm px-5 py-2 rounded-full transition-colors disabled:opacity-50"
          onClick={handleAddComment}
          disabled={posting || !text.trim()}
        >
          {posting ? "..." : "Reply"}
        </button>
      </div>
    </div>
  );
}