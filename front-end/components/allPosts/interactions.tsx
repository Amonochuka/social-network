"use client";

import { useEffect, useRef, useState } from "react";
import { BookMarkedIcon, Heart, Image as ImageIcon, MessageCircleReply, X } from "lucide-react";
import { Api } from "@/services/axios";
import UserProfileImage from "../header/profile/userProfile";

interface Props {
  likes: number;
  likedByMe?: boolean;
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

export function PostInteractions({ likes, likedByMe = false, comments, postId, postDetails }: Props) {
  const [showComments, setShowComments] = useState(false);
  const [isBookmarked, setIsBookmarked] = useState(false);

  // Like state — initialised from props, updated optimistically
  const [liked, setLiked] = useState(likedByMe);
  const [likeCount, setLikeCount] = useState(likes);
  const [liking, setLiking] = useState(false);

  // Keep in sync if the parent re-renders with fresh data
  useEffect(() => {
    setLiked(likedByMe);
    setLikeCount(likes);
  }, [likedByMe, likes]);

  useEffect(() => {
    try {
      const bookmarks = JSON.parse(localStorage.getItem("saved_posts") || "[]");
      setIsBookmarked(bookmarks.some((b: any) => b.id === postId));
    } catch {
      setIsBookmarked(false);
    }
  }, [postId]);

  const handleLike = async () => {
    if (liking) return;
    // Optimistic update
    const prevLiked = liked;
    const prevCount = likeCount;
    setLiked(!liked);
    setLikeCount(liked ? likeCount - 1 : likeCount + 1);
    setLiking(true);
    try {
      const res = await Api.post<{ liked: boolean; like_count: number }>(
        `/posts/${postId}/like`
      );
      // Sync with server truth
      setLiked(res.data.liked);
      setLikeCount(res.data.like_count);
    } catch (err) {
      // Rollback on failure
      setLiked(prevLiked);
      setLikeCount(prevCount);
      console.error("Failed to toggle like:", err);
    } finally {
      setLiking(false);
    }
  };

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
          <button
            onClick={handleLike}
            disabled={liking}
            className="group flex items-center gap-1.5 transition-colors"
          >
            <div className={`p-2 rounded-full transition-colors ${
              liked
                ? "text-pink-500"
                : "group-hover:bg-pink-500/10 group-hover:text-pink-500"
            }`}>
              <Heart
                size={18}
                fill={liked ? "currentColor" : "none"}
                className={liked ? "text-pink-500" : ""}
              />
            </div>
            <span className={`text-sm font-medium transition-colors ${
              liked ? "text-pink-500" : "group-hover:text-pink-500"
            }`}>
              {likeCount}
            </span>
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
  const [mediaFile, setMediaFile] = useState<File | null>(null);
  const [mediaPreview, setMediaPreview] = useState<string | null>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);

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

  const handleAddComment = async () => {
    if (!text.trim() && !mediaFile) return;
    if (posting) return;
    setPosting(true);

    try {
      const formData = new FormData();
      formData.append("content", text);
      if (mediaFile) formData.append("media", mediaFile);

      const res = await Api.post(`/posts/${postId}/comments`, formData, {
        headers: { "Content-Type": "multipart/form-data" },
      });

      setComments((prev) => [...prev, res.data]);
      setText("");
      setMediaFile(null);
      setMediaPreview(null);
      if (fileInputRef.current) fileInputRef.current.value = "";
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
                {c.content && <p className="text-white mt-0.5 break-words">{c.content}</p>}
                {c.media_path && (
                  <img
                    src={c.media_path.startsWith("http") ? c.media_path : `http://localhost:8080/${c.media_path}`}
                    alt="comment media"
                    className="mt-2 max-h-48 rounded-xl object-cover border border-white/10"
                  />
                )}
              </div>
            </div>
          ))}
        </div>
      )}

      {/* Media preview */}
      {mediaPreview && (
        <div className="relative w-fit">
          <img src={mediaPreview} alt="preview" className="max-h-32 rounded-xl object-cover border border-white/10" />
          <button
            onClick={removeMedia}
            className="absolute -top-2 -right-2 bg-black/80 rounded-full p-0.5 hover:bg-black transition-colors"
          >
            <X size={14} className="text-white" />
          </button>
        </div>
      )}

      {/* Add Comment Input */}
      <div className="flex items-center gap-3 mt-2">
        <input
          ref={fileInputRef}
          type="file"
          accept="image/jpeg,image/png,image/gif"
          className="hidden"
          onChange={handleMediaSelect}
        />
        <button
          onClick={() => fileInputRef.current?.click()}
          className="p-2 rounded-full hover:bg-[--primary-theme]/10 text-[--primary-theme] transition-colors shrink-0"
          title="Attach image"
        >
          <ImageIcon size={18} />
        </button>
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
          disabled={posting || (!text.trim() && !mediaFile)}
        >
          {posting ? "..." : "Reply"}
        </button>
      </div>
    </div>
  );
}