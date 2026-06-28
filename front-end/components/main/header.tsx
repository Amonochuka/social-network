"use client";

import { useEffect, useState } from "react";
import { ArrowRight, Menu, Search } from "lucide-react";
import { FormState } from "@/types";
import { useActionState } from "react";
import { Home, Bell, MessageSquareDot } from "lucide-react";
import "@/styles/nav-side-bar.css";
import UserProfileImage from "../header/profile/userProfile";
import { useDispatch, useSelector } from "react-redux";
import { setToggle } from "@/store/features/toggleSideBarSlice";
import { RootState } from "@/store/store";
import { authSelector } from "@/store/features/authSlice";
import Link from "next/link";
import { Api } from "@/services/axios";

interface SearchUIProps {
  placeholder?: string;
  className?: string;
}

export const sendData = async (
  prevState: FormState,
  formData: FormData,
): Promise<FormState> => {
  const query = formData?.get("search") as string;
  return {
    success: true,
    message: `${query} sent successfully`,
    error: null,
  };
};

export function SearchUI({
  placeholder = "Search users, posts, groups…",
  className = "",
}: SearchUIProps) {
  const initialState: FormState = { success: false, message: "", error: null };
  const [, formAction] = useActionState<FormState, FormData>(sendData, initialState);

  return (
    <form action={formAction} className={`search-form ${className}`}>
      <div className="search-inner">
        <Search className="search-icon-svg" size={16} aria-hidden="true" />
        <input
          className="search-input"
          type="search"
          name="search"
          id="search"
          placeholder={placeholder}
          autoComplete="off"
        />
      </div>
    </form>
  );
}

export function HomeNavElements() {
  const [unreadCount, setUnreadCount] = useState(0);

  useEffect(() => {
    const fetchNotifications = async () => {
      try {
        const res = await Api.get<any[]>("/notifications");
        const unread = (res.data ?? []).filter((n) => !n.is_read).length;
        setUnreadCount(unread);
      } catch (err) {
        console.error("Failed to load notifications:", err);
      }
    };
    fetchNotifications();
    const interval = setInterval(fetchNotifications, 15000);
    return () => clearInterval(interval);
  }, []);

  return (
    <div className="display-nav-h-icons flex items-center gap-6">
      <Link href="/view/Home">
        <Home size={23} />
      </Link>
      <Link href="/view/Notifications" className="relative">
        <Bell size={23} />
        {unreadCount > 0 && (
          <span className="absolute -top-1.5 -right-1.5 flex h-4 w-4 items-center justify-center rounded-full bg-red-500 text-[10px] font-bold text-white">
            {unreadCount}
          </span>
        )}
      </Link>
      <Link href="/view/Messages">
        <MessageSquareDot size={23} />
      </Link>
    </div>
  );
}

export function ToggleSideBar() {
  const dispatch = useDispatch();
  const { toggle } = useSelector((state: RootState) => state.toggleSideBar);

  const styleCollapseIcon = {
    cursor: "pointer",
  };
  return (
    <div>
      {toggle ? (
        <ArrowRight
          width={27}
          strokeWidth={3}
          color="var(--primary-theme)"
          onClick={() => dispatch(setToggle())}
        />
      ) : (
        <Menu
          size={30}
          style={styleCollapseIcon}
          color="var(--primary-theme)"
          strokeWidth={2}
          onClick={() => dispatch(setToggle())}
        />
      )}
    </div>
  );
}

export default function HomeProfileUI() {
  const { user } = useSelector(authSelector);

  const avatarUrl = user?.avatar
    ? user.avatar.startsWith("http")
      ? user.avatar
      : `http://localhost:8080/${user.avatar}`
    : undefined;

  const fullName = user ? `${user.first_name} ${user.last_name}` : "?";

  return (
    <div className="home-d-n-main">
      <div>
        <ToggleSideBar />
      </div>
      <div className="display-home-nav">
        <SearchUI />
        <HomeNavElements />
        <Link href="/view/profile">
          <UserProfileImage url={avatarUrl} name={fullName} />
        </Link>
      </div>
    </div>
  );
}