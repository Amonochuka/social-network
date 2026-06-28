"use client";

import { useEffect, useState } from "react";
import { BookMarkedIcon, Heart, MessageCircleReply } from "lucide-react";
import style from "@/styles/interactions.module.css";
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
    <div>
      <div className={style.displayInteractions}>
        <div className={style.likeCommentCont}>
          <div className={style.likeComment}>
            <Heart size={24} />
            <span>{likes}</span>
          </div>
          <div
            className={style.likeComment}
            onClick={() => setShowComments((prev) => !prev)}
            style={{ cursor: "pointer" }}
          >
            <MessageCircleReply size={24} />
            <span>{comments}</span>
          </div>
        </div>
        <span onClick={handleBookmark} style={{ cursor: "pointer" }}>
          <BookMarkedIcon
            size={24}
            className={isBookmarked ? "text-[--primary-theme]" : "text-gray-400"}
            fill={isBookmarked ? "currentColor" : "none"}
          />
        </span>
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
    <div className={style.commentsSection}>
      {loading ? (
        <p className={style.commentLoading}>Loading comments...</p>
      ) : comments.length === 0 ? (
        <p className={style.commentLoading}>No comments yet.</p>
      ) : (
        comments.map((c) => (
          <div key={c.id} className={style.commentRow}>
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
            <div className={style.commentBody}>
              <p className={style.commentAuthor}>{c.author_name}</p>
              <p className={style.commentText}>{c.content}</p>
            </div>
          </div>
        ))
      )}

      <div className={style.commentInputRow}>
        <input
          className={style.commentInput}
          type="text"
          placeholder="Write a comment..."
          value={text}
          onChange={(e) => setText(e.target.value)}
          onKeyDown={(e) => e.key === "Enter" && handleAddComment()}
        />
        <button
          className={style.commentSendBtn}
          onClick={handleAddComment}
          disabled={posting}
        >
          {posting ? "..." : "Send"}
        </button>
      </div>
    </div>
  );
}