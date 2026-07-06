"use client";

import UserPostUI from "../allPosts/allPosts";
import { Follow } from "../sidebar/follow";
import NavSideBar from "../sidebar/sidebar";
import SearchBar from "../header/SearchBar";
import UserProfileImage from "../header/profile/userProfile";
import { useSelector } from "react-redux";
import { authSelector } from "@/store/features/authSlice";

export default function HomeDisplayLayout() {
    const { user } = useSelector(authSelector);

    const avatarUrl = user?.avatar
        ? user.avatar.startsWith("http")
            ? user.avatar
            : `http://localhost:8080/${user.avatar}`
        : undefined;

    const fullName = user
        ? `${user.first_name} ${user.last_name}`
        : "?";

    return (
        <div className="min-h-screen w-full bg-black text-white flex justify-center selection:bg-[--primary-theme] selection:text-white">
            {/* Left Column: Navigation Sidebar */}
            <div className="hidden sm:flex w-[88px] xl:w-[275px] shrink-0 flex-col items-end xl:items-start border-r border-white/10">
                <div className="fixed h-screen w-[88px] xl:w-[275px]">
                    <NavSideBar />
                </div>
            </div>

            {/* Middle Column: Main Feed */}
            <div className="w-full sm:max-w-[600px] border-r border-white/10 min-h-screen flex flex-col relative">
                {/* Sticky Header */}
                <div className="sticky top-0 z-40 bg-black/70 backdrop-blur-md border-b border-white/10 px-4 py-2 flex items-center justify-end">
                    <div className="sm:hidden">
                        <UserProfileImage
                            url={avatarUrl}
                            name={fullName}
                        />
                    </div>
                </div>

                {/* Posts Feed */}
                <div className="w-full flex-1 pt-11 pb-20 sm:pb-0">
                    <UserPostUI />
                </div>
            </div>

            {/* Right Column: Explore & Follow */}
            <div className="hidden lg:block w-[350px] shrink-0 pl-8 py-4">
                <div className="sticky top-4 flex flex-col gap-6">
                    {/* Search Bar */}
                    <SearchBar className="w-full" />

                    {/* Follow widget */}
                    <div className="bg-[#16181c] rounded-2xl border border-white/5 overflow-hidden">
                        <h2 className="px-4 pt-4 pb-2 text-xl font-extrabold text-white">
                            Who to follow
                        </h2>
                        <Follow />
                    </div>
                </div>
            </div>
        </div>
    );
}