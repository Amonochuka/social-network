"use client";

import { sidebarNav } from "@/data";
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

  return (
    <div className="flex flex-col h-full items-end xl:items-start pt-4 pb-6 px-2 xl:px-8">
      <Logo />
      <div className="flex flex-col gap-2 mt-4 w-full xl:w-auto items-end xl:items-start">
        {sidebarNav.map((v) => {
          const isLogout = v.text === "Logout";
          const isActive = v.text === "Home" 
            ? (path === "/" || path === "/view/Home") 
            : path === v.href;

          return (
            <Link
              href={v.href}
              key={v.id}
              onClick={isLogout ? handleLogout : undefined}
              className={`group flex items-center justify-center xl:justify-start gap-4 p-3 xl:px-5 xl:py-3.5 rounded-full transition-all duration-200 w-fit ${
                isActive 
                  ? "font-bold text-white" 
                  : "font-medium text-gray-300 hover:bg-white/10"
              }`}
            >
              <span className={`relative flex items-center justify-center transition-transform group-hover:scale-110 ${isActive ? "text-white" : ""}`}>
                {<v.icon size={26} strokeWidth={isActive ? 2.5 : 2} />}
                {/* Active indicator dot for mobile/xl if needed, but bold icon is usually enough */}
              </span>
              <span className={`hidden xl:inline text-xl ${isActive ? "font-bold" : "font-medium"}`}>
                {v.text}
              </span>
            </Link>
          );
        })}
      </div>
      
      {/* Example X-style Create Post button at bottom of sidebar */}
      <div className="hidden xl:block mt-8 w-full">
        <button className="w-full bg-[--primary-theme] hover:bg-[#129c94] text-white font-bold text-lg rounded-full py-3.5 transition-colors shadow-lg">
          Post
        </button>
      </div>
      <div className="xl:hidden mt-8 flex justify-end">
        <button className="bg-[--primary-theme] hover:bg-[#129c94] text-white rounded-full p-3 transition-colors shadow-lg">
           {/* Add a generic plus icon for mobile compose button */}
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