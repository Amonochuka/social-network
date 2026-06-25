"use client";

import { useState } from "react";
import {
  MoreHorizontal,
  Pencil,
  MapPin,
  Calendar,
} from "lucide-react";
import { useAppSelector } from "@/store/hooks";
import { authSelector } from "@/store/features/authSlice";
import { Api } from "@/services/axios";

export default function ProfilePage() {
  const { user } = useAppSelector(authSelector);

  const [isPublic, setIsPublic] = useState(user?.is_public ?? true);
  const [updating, setUpdating] = useState(false);

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
          <button className="flex items-center gap-2 rounded-xl border border-white/20 bg-black/30 px-4 py-2 text-sm font-semibold text-white backdrop-blur-sm transition hover:bg-black/50">
            <Pencil size={14} />
            Edit Profile
          </button>
          <button className="grid h-9 w-9 place-items-center rounded-xl border border-white/20 bg-black/30 text-white backdrop-blur-sm transition hover:bg-black/50">
            <MoreHorizontal size={17} />
          </button>

          {/* private profile toggle — wired to real backend */}
          <div className="flex items-center gap-2 rounded-xl border border-white/20 bg-black/30 px-3 py-1.5 backdrop-blur-sm">
            <span className="text-xs font-medium text-white/80">
              {isPublic ? "Public Profile" : "Private Profile"}
            </span>
            <button
              role="switch"
              aria-checked={!isPublic}
              disabled={updating}
              onClick={handlePrivacyToggle}
              className={`relative h-5 w-9 shrink-0 rounded-full transition ${
                !isPublic ? "bg-[--primary-theme]" : "bg-white/20"
              } ${updating ? "opacity-50" : ""}`}
            >
              <span
                className={`absolute top-0.5 h-4 w-4 rounded-full bg-white shadow transition-transform ${
                  !isPublic ? "translate-x-[17px]" : "translate-x-0.5"
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
      </div>

      {/* posts/followers will be wired in next */}
      <div className="flex-1 p-6">
        <div className="rounded-2xl bg-[#222] p-8 text-center text-sm text-gray-500">
          Posts, followers, and following counts coming next.
        </div>
      </div>
    </div>
  );
}