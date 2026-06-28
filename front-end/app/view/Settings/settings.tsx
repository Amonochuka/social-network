"use client";

import { useEffect, useState } from "react";
import { ShieldCheck } from "lucide-react";
import { Api } from "@/services/axios";
import { useAppSelector, useAppDispatch } from "@/store/hooks";
import { authSelector, setSession } from "@/store/features/authSlice";

export default function SettingsPage() {
  const { user } = useAppSelector(authSelector);
  const dispatch = useAppDispatch();
  const [isPublic, setIsPublic] = useState(user?.is_public ?? true);
  const [updating, setUpdating] = useState(false);

  useEffect(() => {
    if (user) {
      setIsPublic(user.is_public);
    }
  }, [user]);

  const handlePrivacyToggle = async () => {
    setUpdating(true);
    try {
      const nextVal = !isPublic;
      await Api.post("/profile/privacy", { is_public: nextVal });
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

  return (
    <main className="h-full w-full px-8 py-6 bg-[#181818] overflow-y-auto">
      <div className="max-w-2xl">
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

          {/* Account Details Section */}
          <div className="rounded-2xl border border-white/10 bg-[#1e1e1e] p-6 shadow-lg">
            <div className="flex items-center gap-3 mb-4">
              <h2 className="text-lg font-bold text-white">Account Info</h2>
            </div>

            <div className="space-y-4">
              <div className="flex justify-between py-2.5 border-b border-white/5">
                <div>
                  <p className="text-xs text-gray-500">Email Address</p>
                  <p className="text-sm font-medium text-white mt-0.5">{user?.email || "N/A"}</p>
                </div>
              </div>

              <div className="flex justify-between py-2.5 border-b border-white/5">
                <div>
                  <p className="text-xs text-gray-500">First Name</p>
                  <p className="text-sm font-medium text-white mt-0.5">{user?.first_name || "N/A"}</p>
                </div>
              </div>

              <div className="flex justify-between py-2.5 border-b border-white/5">
                <div>
                  <p className="text-xs text-gray-500">Last Name</p>
                  <p className="text-sm font-medium text-white mt-0.5">{user?.last_name || "N/A"}</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>
  );
}