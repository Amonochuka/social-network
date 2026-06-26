"use client";

import { sidebarNav } from "@/data";
import { RootState } from "@/store/store";
import "@/styles/nav-side-bar.css";
import { activeRoute, nonActive, navSideBar } from "@/styles/style";
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

  const newStyle = {
    ...navSideBar,
    flex: "0 0 5.5rem",
    display: "flex",
    flexDirection: "column",
    gap: "2rem",
    alignItems: "center",
    justifyContent: "start",
  };

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
    <div style={toggle ? newStyle : navSideBar}>
      <Logo />
      <div className="ns-main-cont">
        {sidebarNav.map((v) => {
          const isLogout = v.text === "Logout";

          return (
            <Link
              href={isLogout ? "#" : `/view/${v.text}`}
              key={v.id}
              onClick={isLogout ? handleLogout : undefined}
              style={
                v.text === "Home" && (path === "/" || path === "/view/Home")
                  ? activeRoute
                  : path === "/view/" + v.text
                    ? activeRoute
                    : nonActive
              }
            >
              <span>{<v.icon />}</span>
              {!toggle && <span>{v.text}</span>}
            </Link>
          );
        })}
      </div>
    </div>
  );
}

export function Logo() {
  const styleLogo = {
    padding: "2rem 1rem",
    display: "flex",
    alignItems: "center",
    justifyContent: "space-between",
  };

  return (
    <div  style={styleLogo}>
      <Link href={"/view/Home"}>
        <Image
          src={"/social-network-logo.svg"}
          width={50}
          height={50}
          alt="social-app logo"
          loading="eager"
        />
      </Link>
    </div>
  );
}