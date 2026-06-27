"use client";

import { useEffect, useState } from "react";
import UserProfileImage from "../header/profile/userProfile";
import style from "@/styles/all-post.module.css";
import { Button } from "../ui/button";
import { ButtonData } from "@/types";
import { createPostBtn } from "@/styles/style";
import { Clapperboard, Eye, LockIcon, Image as LucideImage } from "lucide-react";
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
  created_at: string;
  updated_at: string;
}

export default function UserPostUI() {
  const { user } = useAppSelector(authSelector);
  const [postText, setPostText] = useState("");
  const [posting, setPosting] = useState(false);
  const [refreshKey, setRefreshKey] = useState(0);

  const handleCreatePost = async () => {
    if (!postText.trim() || posting) return;
    setPosting(true);

    try {
      const formData = new FormData();
      formData.append("content", postText);
      formData.append("privacy", "public");

      await Api.post("/posts", formData, {
        headers: { "Content-Type": "multipart/form-data" },
      });

      setPostText("");
      setRefreshKey((k) => k + 1); // triggers AllPostUI to refetch
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
        <div className={style.iconPostDisplay}>
          <div className={style.iconFlex}>
            <LucideImage />
            <span>photos</span>
          </div>
          <div className={style.iconFlex}>
            <Clapperboard />
            <span>videos</span>
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
            description={post.content}
            postImage={
              post.media_path
                ? `http://localhost:8080/${post.media_path}`
                : undefined
            }
            likes={0}
            comments={0}
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
  description: string;
  postImage?: string;
  likes: number;
  comments: number;
}

export function UserPostContent({ description, postImage, likes, comments }: ContentInterface) {
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
          <PostInteractions comments={comments} likes={likes} />
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