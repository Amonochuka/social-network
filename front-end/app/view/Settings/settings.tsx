"use client";

import { useEffect, useState } from "react";
import { ShieldCheck, PaintBucket, Lock, User, Trash2, X, AlertTriangle, Sparkles, CheckCircle2, Eye } from "lucide-react";
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
  const [passwordStatus, setPasswordStatus] = useState<{msg: string, type: "success" | "error"} | null>(null);

  // Privacy Modal
  const [showPrivacyModal, setShowPrivacyModal] = useState(false);
  const [pendingPrivacyVal, setPendingPrivacyVal] = useState<boolean | null>(null);

  useEffect(() => {
    if (user) {
      setIsPublic(user.is_public);
      setFirstName(user.first_name);
      setLastName(user.last_name);
      setNickname(user.nickname || "");
      setAboutMe(user.about_me || "");
    }
  }, [user]);

  const handlePrivacyClick = () => {
    setPendingPrivacyVal(!isPublic);
    setShowPrivacyModal(true);
  };

  const confirmPrivacyToggle = async () => {
    if (pendingPrivacyVal === null) return;
    setUpdating(true);
    try {
      await Api.put("/profile/privacy", { is_public: pendingPrivacyVal });
      setIsPublic(pendingPrivacyVal);
      if (user) {
        dispatch(setSession({ ...user, is_public: pendingPrivacyVal }));
      }
      setShowPrivacyModal(false);
    } catch (err) {
      console.error("Failed to update privacy:", err);
    } finally {
      setUpdating(false);
      setPendingPrivacyVal(null);
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
      setPasswordStatus({ msg: "Please fill in both fields.", type: "error" });
      return;
    }
    setPasswordStatus({ msg: "Password changed successfully (mock demo).", type: "success" });
    setCurrentPassword("");
    setNewPassword("");
    setTimeout(() => setPasswordStatus(null), 3000);
  };

  return (
    <main className="relative h-full w-full px-4 md:px-8 py-8 bg-[#0f0f0f] overflow-y-auto overflow-x-hidden selection:bg-[--primary-theme] selection:text-white">
      {/* Dynamic Background Glows */}
      <div className="absolute top-[-10%] left-[-10%] w-[40%] h-[40%] rounded-full bg-[--primary-theme] blur-[150px] opacity-20 pointer-events-none" />
      <div className="absolute bottom-[-10%] right-[-10%] w-[30%] h-[30%] rounded-full bg-[#14afa7] blur-[150px] opacity-10 pointer-events-none" />

      <div className="relative max-w-3xl mx-auto pb-16 z-10">
        {/* Header */}
        <div className="mb-10 text-center md:text-left flex items-center justify-between">
          <div>
            <h1 className="text-4xl font-extrabold text-transparent bg-clip-text bg-gradient-to-r from-white to-gray-400 tracking-tight flex items-center gap-3">
              Settings <Sparkles className="text-[--primary-theme] animate-pulse" size={28} />
            </h1>
            <p className="mt-2 text-sm text-gray-400 font-medium max-w-md">
              Manage your digital footprint, customize your aesthetic, and secure your account.
            </p>
          </div>
        </div>

        <div className="space-y-8">
          {/* Privacy Section */}
          <section className="group rounded-3xl border border-white/5 bg-white/[0.02] backdrop-blur-xl p-7 shadow-2xl hover:border-[--primary-theme]/30 transition-all duration-500 relative overflow-hidden">
            <div className="absolute top-0 right-0 w-32 h-32 bg-gradient-to-br from-[--primary-theme]/20 to-transparent rounded-full blur-2xl opacity-0 group-hover:opacity-100 transition-opacity duration-500" />
            
            <div className="flex items-center gap-3 mb-5 relative z-10">
              <div className="p-2.5 rounded-xl bg-[--primary-theme]/10 text-[--primary-theme] shadow-[0_0_15px_rgba(20,175,167,0.2)]">
                <ShieldCheck size={24} />
              </div>
              <h2 className="text-xl font-bold text-white tracking-wide">Privacy Preferences</h2>
            </div>
            
            <div className="flex items-center justify-between gap-6 py-4 relative z-10 bg-black/20 rounded-2xl p-5 border border-white/5">
              <div className="flex-1">
                <p className="text-base font-semibold text-white flex items-center gap-2">
                  Public Account 
                  {isPublic && <span className="px-2 py-0.5 rounded-full bg-[--primary-theme]/20 text-[--primary-theme] text-[10px] uppercase font-bold tracking-wider">Active</span>}
                </p>
                <p className="mt-1.5 text-xs text-gray-400 leading-relaxed max-w-md">
                  When enabled, any user can view your profile details and posts. When disabled, your profile is strictly private.
                </p>
              </div>

              <button
                onClick={handlePrivacyClick}
                disabled={updating}
                className={`relative h-8 w-14 rounded-full transition-all duration-300 ease-out shrink-0 disabled:opacity-50 shadow-inner focus:outline-none focus:ring-2 focus:ring-[--primary-theme]/50 focus:ring-offset-2 focus:ring-offset-[#111] ${
                  isPublic ? "bg-[--primary-theme] shadow-[0_0_20px_rgba(20,175,167,0.4)]" : "bg-white/10"
                }`}
              >
                <span
                  className={`absolute top-[3px] left-[3px] h-[26px] w-[26px] rounded-full bg-white shadow-md transition-transform duration-300 ease-out flex items-center justify-center ${
                    isPublic ? "translate-x-6" : "translate-x-0"
                  }`}
                >
                  {isPublic ? (
                    <CheckCircle2 size={14} className="text-[--primary-theme]" />
                  ) : (
                    <Lock size={14} className="text-gray-400" />
                  )}
                </span>
              </button>
            </div>
          </section>

          {/* Appearance Section */}
          <section className="group rounded-3xl border border-white/5 bg-white/[0.02] backdrop-blur-xl p-7 shadow-2xl hover:border-[--primary-theme]/30 transition-all duration-500 relative overflow-hidden">
            <div className="flex items-center gap-3 mb-5 relative z-10">
              <div className="p-2.5 rounded-xl bg-[--primary-theme]/10 text-[--primary-theme]">
                <PaintBucket size={24} />
              </div>
              <h2 className="text-xl font-bold text-white tracking-wide">Appearance</h2>
            </div>

            <div className="flex items-center justify-between gap-6 py-4 relative z-10 bg-black/20 rounded-2xl p-5 border border-white/5">
              <div>
                <p className="text-base font-semibold text-white">Application Theme</p>
                <p className="mt-1 text-xs text-gray-400">Personalize your viewing experience.</p>
              </div>

              <div className="flex rounded-xl bg-black/50 p-1.5 w-56 shrink-0 border border-white/10 relative">
                <div 
                  className="absolute top-1.5 bottom-1.5 w-[calc(50%-6px)] bg-[--primary-theme] rounded-lg transition-all duration-300 ease-out shadow-[0_0_15px_rgba(20,175,167,0.4)]"
                  style={{ left: theme === "light" ? "6px" : "calc(50% + 0px)" }}
                />
                {(["light", "dark"] as const).map((t) => (
                  <button
                    key={t}
                    onClick={() => setTheme(t)}
                    className={`flex-1 relative z-10 rounded-lg py-2 text-sm font-bold capitalize transition-colors duration-300 ${
                      theme === t ? "text-white" : "text-gray-500 hover:text-gray-300"
                    }`}
                  >
                    {t}
                  </button>
                ))}
              </div>
            </div>
          </section>

          {/* Profile Info Form */}
          <form onSubmit={handleSaveProfile} className="group rounded-3xl border border-white/5 bg-white/[0.02] backdrop-blur-xl p-7 shadow-2xl hover:border-[--primary-theme]/30 transition-all duration-500 relative overflow-hidden space-y-6">
            <div className="flex items-center gap-3 relative z-10">
              <div className="p-2.5 rounded-xl bg-[--primary-theme]/10 text-[--primary-theme]">
                <User size={24} />
              </div>
              <h2 className="text-xl font-bold text-white tracking-wide">Profile Details</h2>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-5 relative z-10">
              <div className="space-y-1.5">
                <label className="text-[11px] uppercase tracking-wider font-bold text-gray-500 ml-1">First Name</label>
                <input
                  type="text"
                  value={firstName}
                  onChange={(e) => setFirstName(e.target.value)}
                  className="w-full bg-black/30 border border-white/10 rounded-xl px-4 py-3 text-sm text-white font-medium transition-all focus:outline-none focus:border-[--primary-theme] focus:ring-1 focus:ring-[--primary-theme]/50 placeholder:text-gray-600"
                  placeholder="John"
                  required
                />
              </div>
              <div className="space-y-1.5">
                <label className="text-[11px] uppercase tracking-wider font-bold text-gray-500 ml-1">Last Name</label>
                <input
                  type="text"
                  value={lastName}
                  onChange={(e) => setLastName(e.target.value)}
                  className="w-full bg-black/30 border border-white/10 rounded-xl px-4 py-3 text-sm text-white font-medium transition-all focus:outline-none focus:border-[--primary-theme] focus:ring-1 focus:ring-[--primary-theme]/50 placeholder:text-gray-600"
                  placeholder="Doe"
                  required
                />
              </div>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-5 relative z-10">
              <div className="space-y-1.5">
                <label className="text-[11px] uppercase tracking-wider font-bold text-gray-500 ml-1">Nickname</label>
                <input
                  type="text"
                  value={nickname}
                  onChange={(e) => setNickname(e.target.value)}
                  className="w-full bg-black/30 border border-white/10 rounded-xl px-4 py-3 text-sm text-white font-medium transition-all focus:outline-none focus:border-[--primary-theme] focus:ring-1 focus:ring-[--primary-theme]/50 placeholder:text-gray-600"
                  placeholder="@johndoe"
                />
              </div>
              <div className="space-y-1.5">
                <label className="text-[11px] uppercase tracking-wider font-bold text-gray-500 ml-1 flex items-center justify-between">
                  <span>Email Address</span>
                  <span className="text-[--primary-theme]">Read Only</span>
                </label>
                <input
                  type="text"
                  value={user?.email || ""}
                  disabled
                  className="w-full bg-white/[0.02] border border-white/5 rounded-xl px-4 py-3 text-sm text-gray-500 font-medium cursor-not-allowed"
                />
              </div>
            </div>

            <div className="space-y-1.5 relative z-10">
              <label className="text-[11px] uppercase tracking-wider font-bold text-gray-500 ml-1">About Me</label>
              <textarea
                value={aboutMe}
                onChange={(e) => setAboutMe(e.target.value)}
                rows={4}
                className="w-full bg-black/30 border border-white/10 rounded-xl px-4 py-3 text-sm text-white font-medium transition-all focus:outline-none focus:border-[--primary-theme] focus:ring-1 focus:ring-[--primary-theme]/50 placeholder:text-gray-600 resize-none"
                placeholder="Tell the world something interesting about yourself..."
              />
            </div>

            <div className="flex justify-end pt-2 relative z-10">
              <button
                type="submit"
                disabled={updating}
                className="group/btn relative overflow-hidden bg-[--primary-theme] hover:bg-[#129c94] transition-colors text-white text-sm font-bold px-8 py-3 rounded-xl shadow-[0_0_20px_rgba(20,175,167,0.3)] disabled:opacity-70 disabled:cursor-not-allowed"
              >
                <span className="relative z-10">Save Changes</span>
                <div className="absolute inset-0 h-full w-full bg-gradient-to-r from-transparent via-white/20 to-transparent -translate-x-full group-hover/btn:animate-[shimmer_1.5s_infinite]" />
              </button>
            </div>
          </form>

          {/* Password Management */}
          <form onSubmit={handlePasswordChange} className="group rounded-3xl border border-white/5 bg-white/[0.02] backdrop-blur-xl p-7 shadow-2xl hover:border-[--primary-theme]/30 transition-all duration-500 relative overflow-hidden space-y-6">
            <div className="flex items-center gap-3 relative z-10">
              <div className="p-2.5 rounded-xl bg-[--primary-theme]/10 text-[--primary-theme]">
                <Lock size={24} />
              </div>
              <h2 className="text-xl font-bold text-white tracking-wide">Security & Authentication</h2>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-5 relative z-10">
              <div className="space-y-1.5">
                <label className="text-[11px] uppercase tracking-wider font-bold text-gray-500 ml-1">Current Password</label>
                <input
                  type="password"
                  value={currentPassword}
                  onChange={(e) => setCurrentPassword(e.target.value)}
                  placeholder="••••••••"
                  className="w-full bg-black/30 border border-white/10 rounded-xl px-4 py-3 text-sm text-white font-medium transition-all focus:outline-none focus:border-[--primary-theme] focus:ring-1 focus:ring-[--primary-theme]/50 placeholder:text-gray-600"
                />
              </div>
              <div className="space-y-1.5">
                <label className="text-[11px] uppercase tracking-wider font-bold text-gray-500 ml-1">New Password</label>
                <input
                  type="password"
                  value={newPassword}
                  onChange={(e) => setNewPassword(e.target.value)}
                  placeholder="••••••••"
                  className="w-full bg-black/30 border border-white/10 rounded-xl px-4 py-3 text-sm text-white font-medium transition-all focus:outline-none focus:border-[--primary-theme] focus:ring-1 focus:ring-[--primary-theme]/50 placeholder:text-gray-600"
                />
              </div>
            </div>

            {passwordStatus && (
              <div className={`p-3 rounded-lg border text-sm font-semibold flex items-center gap-2 ${
                passwordStatus.type === "success" 
                  ? "bg-green-500/10 border-green-500/20 text-green-400" 
                  : "bg-red-500/10 border-red-500/20 text-red-400"
              }`}>
                {passwordStatus.type === "success" ? <CheckCircle2 size={18} /> : <AlertTriangle size={18} />}
                {passwordStatus.msg}
              </div>
            )}

            <div className="flex justify-end pt-2 relative z-10">
              <button
                type="submit"
                className="bg-white/5 hover:bg-white/10 border border-white/10 transition-colors text-white text-sm font-bold px-8 py-3 rounded-xl"
              >
                Update Password
              </button>
            </div>
          </form>

          {/* Danger Zone */}
          <section className="rounded-3xl border border-red-500/20 bg-[#1a0f0f] p-7 shadow-2xl relative overflow-hidden mt-12">
            <div className="absolute top-0 left-0 w-full h-1 bg-gradient-to-r from-red-500/0 via-red-500 to-red-500/0 opacity-50" />
            
            <div className="flex items-center gap-3 mb-3 relative z-10">
              <div className="p-2.5 rounded-xl bg-red-500/10 text-red-500">
                <AlertTriangle size={24} />
              </div>
              <h2 className="text-xl font-bold text-red-500 tracking-wide">Danger Zone</h2>
            </div>
            
            <p className="text-sm text-red-400/80 mb-6 max-w-2xl">
              Permanently remove your Personal Account and all of its contents from the Social Network platform. This action is not reversible, so please continue with caution.
            </p>
            
            <button
              onClick={() => alert("Account deletion is disabled for demo purposes.")}
              className="group flex items-center gap-2 bg-red-500/10 hover:bg-red-500/20 text-red-500 border border-red-500/30 rounded-xl px-6 py-3 text-sm font-bold transition-all hover:shadow-[0_0_20px_rgba(239,68,68,0.2)]"
            >
              <Trash2 size={16} className="group-hover:rotate-12 transition-transform" />
              Delete Account
            </button>
          </section>
        </div>
      </div>

      {/* Privacy Confirmation Modal (Bonus Requirement) */}
      {showPrivacyModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-sm animate-in fade-in duration-200">
          <div className="bg-[#1e1e1e] border border-white/10 rounded-3xl p-6 md:p-8 max-w-md w-full shadow-2xl transform transition-all scale-100 animate-in zoom-in-95 duration-200">
            <div className="flex items-start justify-between mb-5">
              <div className="p-3 rounded-full bg-[--primary-theme]/10 text-[--primary-theme]">
                {pendingPrivacyVal ? <Eye size={28} /> : <Lock size={28} />}
              </div>
              <button 
                onClick={() => {
                  setShowPrivacyModal(false);
                  setPendingPrivacyVal(null);
                }}
                className="p-1 rounded-full hover:bg-white/10 text-gray-400 hover:text-white transition-colors"
              >
                <X size={20} />
              </button>
            </div>
            
            <h3 className="text-2xl font-bold text-white mb-2">
              Make Account {pendingPrivacyVal ? "Public" : "Private"}?
            </h3>
            
            <p className="text-sm text-gray-400 mb-8 leading-relaxed">
              {pendingPrivacyVal 
                ? "Your profile will become public. Anyone will be able to see your posts, followers, and following list without needing your approval." 
                : "Your profile will become private. Only approved followers will be able to see your posts and profile details. Existing followers won't be affected."}
            </p>
            
            <div className="flex items-center justify-end gap-3">
              <button
                onClick={() => {
                  setShowPrivacyModal(false);
                  setPendingPrivacyVal(null);
                }}
                className="px-5 py-2.5 rounded-xl text-sm font-bold text-white hover:bg-white/10 transition-colors"
              >
                Cancel
              </button>
              <button
                onClick={confirmPrivacyToggle}
                disabled={updating}
                className="px-6 py-2.5 rounded-xl text-sm font-bold text-white bg-[--primary-theme] hover:bg-[#129c94] transition-colors shadow-lg disabled:opacity-50"
              >
                {updating ? "Updating..." : "Confirm Change"}
              </button>
            </div>
          </div>
        </div>
      )}
      
      {/* Required keyframes for animations */}
      <style dangerouslySetInnerHTML={{__html: `
        @keyframes shimmer {
          100% { transform: translateX(100%); }
        }
      `}} />
    </main>
  );
}