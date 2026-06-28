"use client";

import UserProfileImage from "@/components/header/profile/userProfile";
import DefaultLayout from "@/components/layouts/defaultLayout";
import { Button } from "@/components/ui/button";
import { ButtonData } from "@/types";
import style from "@/styles/followers.module.css";
import { SearchUI } from "@/components/main/header";
import { useEffect, useState } from "react";
import { Api } from "@/services/axios";
import Link from "next/link";

type NetworkType = "Followers" | "Following";

interface FollowerProfile {
  user_id: string;
  first_name: string;
  last_name: string;
  avatar: string;
  nickname: string;
}

interface PeopleProp {
  data: FollowerProfile;
  networkType: NetworkType;
  onAction: (userId: string) => void;
  loading: string | null;
}

interface ActiveRoutesProps {
  count: number;
  active: NetworkType;
  setActive: React.Dispatch<React.SetStateAction<NetworkType>>;
}

export default function FollowersUI() {
  return (
    <DefaultLayout>
      <Followers />
    </DefaultLayout>
  );
}

export function Followers() {
  const [active, setActive] = useState<NetworkType>("Followers");
  const [followers, setFollowers] = useState<FollowerProfile[]>([]);
  const [following, setFollowing] = useState<FollowerProfile[]>([]);
  const [loading, setLoading] = useState(true);
  const [actionLoading, setActionLoading] = useState<string | null>(null);

  useEffect(() => {
    const fetchData = async () => {
      setLoading(true);
      try {
        const [followersRes, followingRes] = await Promise.all([
          Api.get<FollowerProfile[]>("/followers"),
          Api.get<FollowerProfile[]>("/following"),
        ]);
        setFollowers(followersRes.data ?? []);
        setFollowing(followingRes.data ?? []);
      } catch (err) {
        console.error("Failed to load network data:", err);
      } finally {
        setLoading(false);
      }
    };
    fetchData();
  }, []);

  const handleUnfollow = async (userId: string) => {
    setActionLoading(userId);
    try {
      await Api.delete(`/follow/${userId}`);
      setFollowing((prev) => prev.filter((f) => f.user_id !== userId));
    } catch (err) {
      console.error("Failed to unfollow:", err);
    } finally {
      setActionLoading(null);
    }
  };

  const handleFollow = async (userId: string) => {
    setActionLoading(userId);
    try {
      await Api.post("/follow/requests", { receiver_id: userId });
      // Move from followers to following if auto-followed (public user)
    } catch (err) {
      console.error("Failed to follow:", err);
    } finally {
      setActionLoading(null);
    }
  };

  const currentList = active === "Followers" ? followers : following;
  const handleAction = active === "Following" ? handleUnfollow : handleFollow;

  return (
    <div className={style.followMainCont}>
      <div className={style.followHeader}>
        <ActiveRoutes
          count={currentList.length}
          active={active}
          setActive={setActive}
        />
      </div>

      <div className={style.mappedFollowers}>
        {loading ? (
          <p style={{ textAlign: "center", color: "#888", padding: "2rem" }}>
            Loading...
          </p>
        ) : currentList.length === 0 ? (
          <p style={{ textAlign: "center", color: "#888", padding: "2rem" }}>
            {active === "Followers"
              ? "No followers yet."
              : "Not following anyone yet."}
          </p>
        ) : (
          currentList.map((person) => (
            <People
              key={person.user_id}
              data={person}
              networkType={active}
              onAction={handleAction}
              loading={actionLoading}
            />
          ))
        )}
      </div>
    </div>
  );
}

export function People({ data, networkType, onAction, loading }: PeopleProp) {
  const btnText = networkType === "Following" ? "Unfollow" : "Follow";
  const isLoading = loading === data.user_id;

  const avatarUrl = data.avatar
    ? data.avatar.startsWith("http")
      ? data.avatar
      : `http://localhost:8080/${data.avatar}`
    : undefined;

  const btnData: ButtonData = {
    text: isLoading ? "..." : btnText,
    onClick: () => onAction(data.user_id),
    style: {
      backgroundColor: networkType === "Following" ? "rgba(255,255,255,0.08)" : "var(--primary-theme)",
      padding: "0.5rem 1.25rem",
      borderRadius: "0.5rem",
      fontSize: "0.875rem",
      fontWeight: "500",
      color: networkType === "Following" ? "#e5e7eb" : undefined,
      border: networkType === "Following" ? "1px solid rgba(255,255,255,0.1)" : "none",
    },
  };

  return (
    <div className={style.peopleCont}>
      <Link
        href={`/view/profile/${data.user_id}`}
        style={{ display: "flex", alignItems: "center", gap: "0.875rem", flex: 1, minWidth: 0, textDecoration: "none" }}
      >
        <UserProfileImage url={avatarUrl} name={`${data.first_name} ${data.last_name}`} />
        <div className={style.personText}>
          <p className={style.personName}>{data.first_name} {data.last_name}</p>
          {data.nickname && (
            <p className={style.personLocation}>@{data.nickname}</p>
          )}
        </div>
      </Link>
      <div>
        <Button data={btnData} />
      </div>
    </div>
  );
}

export function ActiveRoutes({ active, count, setActive }: ActiveRoutesProps) {
  const routes: NetworkType[] = ["Followers", "Following"];

  return (
    <div className={style.networkHeader}>
      <div className={style.networkTitleRow}>
        <div className={style.networkTitleGroup}>
          <p className={style.FollowerText}>{active}</p>
          <span className={style.networkCount}>{count}</span>
        </div>
        <SearchUI placeholder={`Search ${active.toLowerCase()}…`} />
      </div>

      <div className={style.displayNetwork}>
        {routes.map((route) => (
          <button
            key={route}
            onClick={() => setActive(route)}
            className={`${style.networkTab} ${active === route ? style.activeNetwork : ""}`}
          >
            {route}
          </button>
        ))}
      </div>
    </div>
  );
}
