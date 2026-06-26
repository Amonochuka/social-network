"use client";

import { useState } from "react";
import { useAppSelector } from "@/store/hooks";
import { authSelector } from "@/store/features/authSlice";
import LogIn from "@/components/auth/login";
import RegisterUI from "@/components/auth/register";
import HomeDisplayLayout from "@/components/main/display";

export default function Home() {
  const { isAuthenticated, loading } = useAppSelector(authSelector);
  const [showRegister, setShowRegister] = useState(false);

  if (loading) {
    return (
      <div style={{ color: "#fff", padding: "2rem" }}>
        Loading...
      </div>
    );
  }

  if (!isAuthenticated) {
    return showRegister ? (
      <RegisterUI setRegister={setShowRegister} />
    ) : (
      <LogIn setRegister={setShowRegister} />
    );
  }

  return <HomeDisplayLayout />;
}
