"use client";

import { useEffect, useState } from "react";
import { ButtonData } from "@/types";
import { Button } from "../ui/button";
import Image from "next/image";
import style from "@/styles/follow.module.css";
import { btnFollowStyle } from "@/styles/style";
import { Api } from "@/services/axios";

interface FollowerProfile {
  user_id: string;
  first_name: string;
  last_name: string;
  avatar: string;
  nickname: string;
}

export function RightSideBarUI() {
  return <Follow />;
}

export function Follow() {
  const [following, setFollowing] = useState<FollowerProfile[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchFollowing = async () => {
      try {
        const res = await Api.get<FollowerProfile[]>("/following");
        setFollowing(res.data ?? []);
      } catch (err) {
        console.error("Failed to load following:", err);
        setFollowing([]);
      } finally {
        setLoading(false);
      }
    };

    fetchFollowing();
  }, []);

  const handleUnfollow = async (userId: string) => {
    try {
      await Api.delete(`/follow/${userId}`);
      setFollowing((prev) => prev.filter((u) => u.user_id !== userId));
    } catch (err) {
      console.error("Failed to unfollow:", err);
    }
  };

  if (loading) {
    return (
      <div className={style.userSFollowMain}>
        <p className={style.people}>People you follow</p>
        <p style={{ color: "#888", padding: "1rem" }}>Loading...</p>
      </div>
    );
  }

  return (
    <div className={style.userSFollowMain}>
      <p className={style.people}>People you follow</p>
      {following.length === 0 ? (
        <p style={{ color: "#888", padding: "1rem" }}>
          You're not following anyone yet.
        </p>
      ) : (
        following.map((user) => (
          <FollowUserUI
            key={user.user_id}
            user={user}
            onUnfollow={() => handleUnfollow(user.user_id)}
          />
        ))
      )}
    </div>
  );
}

interface FollowUserUIProps {
  user: FollowerProfile;
  onUnfollow: () => void;
}

export function FollowUserUI({ user, onUnfollow }: FollowUserUIProps) {
  const buttonData: ButtonData = {
    text: "unfollow",
    style: btnFollowStyle,
    onClick: onUnfollow,
  };

  const avatarUrl = user.avatar
    ? user.avatar.startsWith("http")
      ? user.avatar
      : `http://localhost:8080/${user.avatar}`
    : "/maodongo.jpeg";

  return (
    <div className={style.userSFollowCont}>
      <div className={style.displayFollowProf}>
        <div className={style.followImage}>
          <Image
            src={avatarUrl}
            priority
            fill
            sizes="48px"
            alt={`${user.first_name} ${user.last_name} profile image`}
          />
        </div>
        <div>
          <p>
            {user.first_name} {user.last_name}
          </p>
          {user.nickname && <p>@{user.nickname}</p>}
        </div>
      </div>
      <div>
        <Button data={buttonData} />
      </div>
    </div>
  );
}