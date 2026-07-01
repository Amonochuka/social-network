"use client";

import { useEffect, useState } from "react";
import { Search } from "lucide-react";
import { useAppDispatch, useAppSelector } from "@/store/hooks";
import {
  searchUsers,
  clearResults,
  searchSelector,
} from "@/store/features/searchSlice";
import Link from "next/link";

export default function SearchBar() {
  const dispatch = useAppDispatch();
  const { results, loading } = useAppSelector(searchSelector);

  const [query, setQuery] = useState("");

  useEffect(() => {
    const timeout = setTimeout(() => {
      if (query.trim().length >= 2) {
        dispatch(searchUsers(query));
      } else {
        dispatch(clearResults());
      }
    }, 300);

    return () => clearTimeout(timeout);
  }, [query, dispatch]);

  return (
    <div className="relative w-full max-w-md">
      <div className="flex items-center rounded-xl bg-[#222] px-3 py-2">
        <Search size={18} className="text-gray-400" />

        <input
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          placeholder="Search users..."
          className="ml-2 w-full bg-transparent text-white outline-none placeholder:text-gray-500"
        />
      </div>

      {query.length >= 2 && (
        <div className="absolute z-50 mt-2 w-full rounded-xl bg-[#222] border border-white/10 shadow-xl">
          {loading ? (
            <div className="p-4 text-sm text-gray-400">
              Searching...
            </div>
          ) : results.length === 0 ? (
            <div className="p-4 text-sm text-gray-500">
              No users found
            </div>
          ) : (
            results.map((user) => (
              <Link
                key={user.id}
                href={`/view/profile/${user.id}`}
                className="flex items-center gap-3 p-3 hover:bg-white/5"
              >
                <div className="h-10 w-10 overflow-hidden rounded-full bg-gray-700">
                  {user.avatar ? (
                    <img
                      src={
                        user.avatar.startsWith("http")
                          ? user.avatar
                          : `http://localhost:8080/${user.avatar}`
                      }
                      alt={user.first_name}
                      className="h-full w-full object-cover"
                    />
                  ) : null}
                </div>

                <div>
                  <p className="text-sm font-semibold text-white">
                    {user.first_name} {user.last_name}
                  </p>

                  <p className="text-xs text-gray-400">
                    @{user.nickname}
                  </p>
                </div>
              </Link>
            ))
          )}
        </div>
      )}
    </div>
  );
}