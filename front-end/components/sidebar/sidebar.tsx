"use client";

import { useState, useRef, useEffect } from "react";
import { sidebarNav, sidebarMoreNav, sidebarMoreTrigger, sidebarLogout } from "@/data";
import { RootState } from "@/store/store";
import Image from "next/image";
import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import { useDispatch, useSelector } from "react-redux";
import { Api } from "@/services/axios";
import { logout } from "@/store/features/authSlice";

export default function NavSideBar() {
  const path = usePathname();
  const router = useRouter();
  const dispatch = useDispatch();
  const { toggle } = useSelector((state: RootState) => state.toggleSideBar);
  const [showMore, setShowMore] = useState(false);
  const moreRef = useRef<HTMLDivElement>(null);

  // close the popover when clicking outside it
  useEffect(() => {
    const handleClickOutside = (e: MouseEvent) => {
      if (moreRef.current && !moreRef.current.contains(e.target as Node)) {
        setShowMore(false);
      }
    };
    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  const handleLogout = async (e: React.MouseEvent) => {
    e.preventDefault();
    try {
      await Api.post("/auth/logout");
    } catch (err) {
      console.error("Logout failed:", err);
    } finally {
      dispatch(logout());
      router.push("/");
    }
  };

  const renderNavItem = (v: { id: number; icon: any; text: string; href: string }) => {
    const isActive = v.text === "Home"
      ? (path === "/" || path === "/view/Home")
      : path === v.href;

    return (
      <Link
        href={v.href}
        key={v.id}
        onClick={() => setShowMore(false)}
        className={`group flex items-center justify-center xl:justify-start gap-4 p-3 xl:px-5 xl:py-3.5 rounded-full transition-all duration-200 w-fit ${
          isActive
            ? "font-bold text-white"
            : "font-medium text-gray-300 hover:bg-white/10"
        }`}
      >
        <span className={`relative flex items-center justify-center transition-transform group-hover:scale-110 ${isActive ? "text-white" : ""}`}>
          {<v.icon size={26} strokeWidth={isActive ? 2.5 : 2} />}
        </span>
        <span className={`hidden xl:inline text-xl ${isActive ? "font-bold" : "font-medium"}`}>
          {v.text}
        </span>
      </Link>
    );
  };

  return (
    <div className="flex flex-col h-full items-end xl:items-start pt-4 pb-6 px-2 xl:px-8 relative">
      <Logo />
      <div className="flex flex-col gap-2 mt-4 w-full xl:w-auto items-end xl:items-start">
        {sidebarNav.map(renderNavItem)}

        {/* More toggle with flyout popover */}
        <div className="relative" ref={moreRef}>
          <button
            onClick={() => setShowMore((prev) => !prev)}
            className="group flex items-center justify-center xl:justify-start gap-4 p-3 xl:px-5 xl:py-3.5 rounded-full transition-all duration-200 w-fit font-medium text-gray-300 hover:bg-white/10"
          >
            <span className="relative flex items-center justify-center transition-transform group-hover:scale-110">
              <sidebarMoreTrigger.icon size={26} strokeWidth={2} />
            </span>
            <span className="hidden xl:inline text-xl font-medium">
              {sidebarMoreTrigger.text}
            </span>
          </button>

          {showMore && (
            <div className="absolute bottom-full left-0 mb-2 w-56 bg-[#222] border border-white/10 rounded-2xl shadow-2xl overflow-hidden py-2 z-50">
              {sidebarMoreNav.map((v) => {
                const isActive = path === v.href;
                return (
                  <Link
                    key={v.id}
                    href={v.href}
                    onClick={() => setShowMore(false)}
                    className={`flex items-center gap-3 px-4 py-3 transition-colors ${
                      isActive ? "text-white font-bold bg-white/5" : "text-gray-300 hover:bg-white/10"
                    }`}
                  >
                    <v.icon size={20} />
                    <span>{v.text}</span>
                  </Link>
                );
              })}
            </div>
          )}
        </div>

        {/* Logout always pinned, visible */}
        <Link
          href={sidebarLogout.href}
          onClick={handleLogout}
          className="group flex items-center justify-center xl:justify-start gap-4 p-3 xl:px-5 xl:py-3.5 rounded-full transition-all duration-200 w-fit font-medium text-gray-300 hover:bg-white/10"
        >
          <span className="relative flex items-center justify-center transition-transform group-hover:scale-110">
            <sidebarLogout.icon size={26} strokeWidth={2} />
          </span>
          <span className="hidden xl:inline text-xl font-medium">
            {sidebarLogout.text}
          </span>
        </Link>
      </div>

      {/* Create Post button at bottom of sidebar */}
      <div className="hidden xl:block mt-8 w-full">
        <button className="w-full bg-[--primary-theme] hover:bg-[#129c94] text-white font-bold text-lg rounded-full py-3.5 transition-colors shadow-lg">
          Post
        </button>
      </div>
      <div className="xl:hidden mt-8 flex justify-end">
        <button className="bg-[--primary-theme] hover:bg-[#129c94] text-white rounded-full p-3 transition-colors shadow-lg">
           <svg viewBox="0 0 24 24" aria-hidden="true" className="w-6 h-6 fill-current"><g><path d="M23 11h-10V1h-2v10H1v2h10v10h2V13h10v-2z"></path></g></svg>
        </button>
      </div>
    </div>
  );
}

export function Logo() {
  return (
    <div className="flex items-center justify-center xl:justify-start w-full xl:w-auto mb-2 p-3">
      <Link href={"/view/Home"} className="rounded-full p-2 hover:bg-[--primary-theme]/10 transition-colors">
        <Image
          src={"/social-network-logo.svg"}
          width={40}
          height={40}
          alt="social-app logo"
          loading="eager"
          className="xl:w-12 xl:h-12"
        />
      </Link>
    </div>
  );
}