"use client";

import { useEffect, useState } from "react";
import { ShieldCheck, PaintBucket, Lock, User, Eye, Lock as LockIcon, Trash2 } from "lucide-react";
import { Api } from "@/services/axios";
import { useAppSelector, useAppDispatch } from "@/store/hooks";
import { authSelector, setSession } from "@/store/features/authSlice";

export default function SettingsPage() {
  const { user } = useAppSelector(authSelector);
  const dispatch = useAppDispatch();
  
  const [isPublic, setIsPublic] = useState(user?.is_public ?? true);
  const [updating, setUpdating] = useState(false);
  const [theme, setTheme] = useState<"light" | "dark">("dark");
  
  // Profile update fields
  const [firstName, setFirstName] = useState(user?.first_name ?? "");
  const [lastName, setLastName] = useState(user?.last_name ?? "");
  const [nickname, setNickname] = useState(user?.nickname ?? "");
  const [aboutMe, setAboutMe] = useState(user?.about_me ?? "");
  
  // Password fields
  const [currentPassword, setCurrentPassword] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [passwordStatus, setPasswordStatus] = useState<string | null>(null);

  useEffect(() => {
    if (user) {
      setIsPublic(user.is_public);
      setFirstName(user.first_name);
      setLastName(user.last_name);
      setNickname(user.nickname || "");
      setAboutMe(user.about_me || "");
    }
  }, [user]);

  const handlePrivacyToggle = async () => {
    setUpdating(true);
    try {
      const nextVal = !isPublic;
      // Fixed: The backend route is a PUT request!
      await Api.put("/profile/privacy", { is_public: nextVal });
      setIsPublic(nextVal);
      if (user) {
        dispatch(setSession({ ...user, is_public: nextVal }));
      }
    } catch (err) {
      console.error("Failed to update privacy:", err);
    } finally {
      setUpdating(false);
    }
  };

  const handleSaveProfile = async (e: React.FormEvent) => {
    e.preventDefault();
    setUpdating(true);
    try {
      const res = await Api.put("/profile", {
        first_name: firstName,
        last_name: lastName,
        date_of_birth: user?.date_of_birth || "",
        nickname: nickname,
        about_me: aboutMe
      });
      if (res.data) {
        dispatch(setSession(res.data));
        alert("Profile details updated successfully!");
      }
    } catch (err) {
      console.error("Failed to update profile info:", err);
      alert("Failed to update profile details.");
    } finally {
      setUpdating(false);
    }
  };

  const handlePasswordChange = (e: React.FormEvent) => {
    e.preventDefault();
    if (!currentPassword || !newPassword) {
      setPasswordStatus("Please fill in both fields.");
      return;
    }
    setPasswordStatus("Password changed successfully (mock demo).");
    setCurrentPassword("");
    setNewPassword("");
  };

  return (
    <main className="h-full w-full px-8 py-6 bg-[#181818] overflow-y-auto">
      <div className="max-w-2xl pb-12">
        {/* Header */}
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-white">Settings</h1>
          <p className="mt-1 text-sm text-gray-400">
            Manage your account settings, profile visibility, and privacy preferences.
          </p>
        </div>

        <div className="space-y-6">
          {/* Privacy Section */}
          <div className="rounded-2xl border border-white/10 bg-[#1e1e1e] p-6 shadow-lg">
            <div className="flex items-center gap-3 mb-4">
              <ShieldCheck className="text-[--primary-theme]" size={22} />
              <h2 className="text-lg font-bold text-white">Privacy Preferences</h2>
            </div>
            
            <div className="flex items-center justify-between gap-4 py-3">
              <div>
                <p className="text-sm font-semibold text-white">Public Account</p>
                <p className="mt-1 text-xs text-gray-400 max-w-md">
                  When enabled, any user can view your profile details and posts. When disabled, your profile is private and other users must request to follow you to view your posts.
                </p>
              </div>

              <button
                onClick={handlePrivacyToggle}
                disabled={updating}
                className={`relative h-7 w-12 rounded-full transition-colors shrink-0 ${
                  isPublic ? "bg-[--primary-theme]" : "bg-white/10"
                } disabled:opacity-50`}
              >
                <span
                  className={`absolute top-0.5 h-6 w-6 rounded-full bg-white shadow-md transition-transform ${
                    isPublic ? "translate-x-[22px]" : "translate-x-0.5"
                  }`}
                />
              </button>
            </div>
          </div>

          {/* Appearance Section */}
          <div className="rounded-2xl border border-white/10 bg-[#1e1e1e] p-6 shadow-lg">
            <div className="flex items-center gap-3 mb-4">
              <PaintBucket className="text-[--primary-theme]" size={22} />
              <h2 className="text-lg font-bold text-white">Appearance</h2>
            </div>

            <div className="flex items-center justify-between gap-4 py-3">
              <div>
                <p className="text-sm font-semibold text-white">Theme</p>
                <p className="mt-1 text-xs text-gray-400">Choose your preferred application theme.</p>
              </div>

              <div className="flex rounded-lg bg-[#111] p-1 w-48 shrink-0">
                {(["light", "dark"] as const).map((t) => (
                  <button
                    key={t}
                    onClick={() => setTheme(t)}
                    className={`flex-1 rounded-md py-1.5 text-xs font-semibold capitalize transition ${
                      theme === t
                        ? "bg-[--primary-theme] text-white"
                        : "text-gray-400 hover:text-white"
                    }`}
                  >
                    {t}
                  </button>
                ))}
              </div>
            </div>
          </div>

          {/* Profile Info Form */}
          <form onSubmit={handleSaveProfile} className="rounded-2xl border border-white/10 bg-[#1e1e1e] p-6 shadow-lg space-y-4">
            <div className="flex items-center gap-3 mb-2">
              <User className="text-[--primary-theme]" size={22} />
              <h2 className="text-lg font-bold text-white">Profile Details</h2>
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div>
                <label className="text-xs text-gray-400 block mb-1">First Name</label>
                <input
                  type="text"
                  value={firstName}
                  onChange={(e) => setFirstName(e.target.value)}
                  className="w-full bg-[#111] border border-white/10 rounded-xl px-4 py-2.5 text-sm text-white focus:outline-none focus:border-[--primary-theme]"
                  required
                />
              </div>
              <div>
                <label className="text-xs text-gray-400 block mb-1">Last Name</label>
                <input
                  type="text"
                  value={lastName}
                  onChange={(e) => setLastName(e.target.value)}
                  className="w-full bg-[#111] border border-white/10 rounded-xl px-4 py-2.5 text-sm text-white focus:outline-none focus:border-[--primary-theme]"
                  required
                />
              </div>
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div>
                <label className="text-xs text-gray-400 block mb-1">Nickname</label>
                <input
                  type="text"
                  value={nickname}
                  onChange={(e) => setNickname(e.target.value)}
                  className="w-full bg-[#111] border border-white/10 rounded-xl px-4 py-2.5 text-sm text-white focus:outline-none focus:border-[--primary-theme]"
                />
              </div>
              <div>
                <label className="text-xs text-gray-400 block mb-1">Email (Read Only)</label>
                <input
                  type="text"
                  value={user?.email || ""}
                  disabled
                  className="w-full bg-white/5 border border-white/5 rounded-xl px-4 py-2.5 text-sm text-gray-500 cursor-not-allowed"
                />
              </div>
            </div>

            <div>
              <label className="text-xs text-gray-400 block mb-1">About Me</label>
              <textarea
                value={aboutMe}
                onChange={(e) => setAboutMe(e.target.value)}
                rows={3}
                className="w-full bg-[#111] border border-white/10 rounded-xl px-4 py-2.5 text-sm text-white focus:outline-none focus:border-[--primary-theme] resize-none"
              />
            </div>

            <div className="flex justify-end">
              <button
                type="submit"
                disabled={updating}
                className="bg-[--primary-theme] hover:opacity-90 transition text-white text-xs font-semibold px-5 py-2.5 rounded-xl"
              >
                Save Details
              </button>
            </div>
          </form>

          {/* Password Management */}
          <form onSubmit={handlePasswordChange} className="rounded-2xl border border-white/10 bg-[#1e1e1e] p-6 shadow-lg space-y-4">
            <div className="flex items-center gap-3 mb-2">
              <Lock className="text-[--primary-theme]" size={22} />
              <h2 className="text-lg font-bold text-white">Security & Password</h2>
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div>
                <label className="text-xs text-gray-400 block mb-1">Current Password</label>
                <input
                  type="password"
                  value={currentPassword}
                  onChange={(e) => setCurrentPassword(e.target.value)}
                  placeholder="••••••••"
                  className="w-full bg-[#111] border border-white/10 rounded-xl px-4 py-2.5 text-sm text-white focus:outline-none focus:border-[--primary-theme]"
                />
              </div>
              <div>
                <label className="text-xs text-gray-400 block mb-1">New Password</label>
                <input
                  type="password"
                  value={newPassword}
                  onChange={(e) => setNewPassword(e.target.value)}
                  placeholder="••••••••"
                  className="w-full bg-[#111] border border-white/10 rounded-xl px-4 py-2.5 text-sm text-white focus:outline-none focus:border-[--primary-theme]"
                />
              </div>
            </div>

            {passwordStatus && (
              <p className="text-xs text-[--primary-theme] font-medium">{passwordStatus}</p>
            )}

            <div className="flex justify-end">
              <button
                type="submit"
                className="bg-[--primary-theme] hover:opacity-90 transition text-white text-xs font-semibold px-5 py-2.5 rounded-xl"
              >
                Change Password
              </button>
            </div>
          </form>

          {/* Danger Zone */}
          <div className="rounded-2xl border border-red-500/20 bg-red-500/5 p-6 shadow-lg">
            <h2 className="text-lg font-bold text-red-400 mb-2">Danger Zone</h2>
            <p className="text-xs text-gray-400 mb-4">
              Once you delete your account, there is no going back. Please be certain.
            </p>
            <button
              onClick={() => alert("Account deletion is disabled for demo purposes.")}
              className="bg-red-500/20 hover:bg-red-500/30 text-red-400 border border-red-500/30 rounded-xl px-5 py-2.5 text-xs font-semibold transition"
            >
              Delete Account
            </button>
          </div>
        </div>
      </div>
    </main>
  );
}