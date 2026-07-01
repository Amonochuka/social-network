"use client";

import { useEffect, useState } from "react";
import { Search } from "lucide-react";
import { useRouter } from "next/navigation";
import { searchService, SearchUser } from "@/services/searchService";

interface SearchResult {
  id: string;
  first_name: string;
  last_name: string;
  nickname: string;
  avatar: string;
}

interface SearchBarProps {
  className?: string;
}

export default function SearchBar({ className = "" }: SearchBarProps) {
  const router = useRouter();

  const [query, setQuery] = useState("");
  const [results, setResults] = useState<SearchUser[]>([]);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    console.log("Query changed:", query)
    if (query.trim().length < 2) {
      setResults([]);
      return;
    }

    const timer = setTimeout(async () => {
        
      setLoading(true);

      try {
        const data = await searchService.searchUsers(query);
        console.log(data);
        setResults(data.users);    
      } catch (err) {
        console.error(err);
      } finally {
        setLoading(false);
      }
    }, 300);

    return () => clearTimeout(timer);
  }, [query]);

  return (
    <div className={`relative ${className}`}>
      <div className="flex items-center rounded-full bg-[#202327] px-4 py-3">
        <Search size={18} className="text-gray-400" />

        <input
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          placeholder="Search users..."
          className="ml-3 flex-1 bg-transparent text-white outline-none placeholder:text-gray-500"
        />
      </div>

      {(loading || results.length > 0) && (
        <div className="absolute mt-2 w-full rounded-2xl bg-[#16181c] shadow-xl border border-white/10 overflow-hidden z-50">

          {loading && (
            <div className="p-4 text-gray-400 text-sm">
              Searching...
            </div>
          )}

          {!loading &&
            results.map((user) => (
              <button
                key={user.id}
                onClick={() => router.push(`/view/profile/${user.id}`)}
                className="flex w-full items-center gap-3 px-4 py-3 hover:bg-white/5 transition"
              >
                <div className="h-10 w-10 rounded-full bg-gray-700 overflow-hidden">
                  {user.avatar ? (
                    <img
                      src={
                        user.avatar.startsWith("http")
                          ? user.avatar
                          : `http://localhost:8080/${user.avatar}`
                      }
                      className="h-full w-full object-cover"
                    />
                  ) : null}
                </div>

                <div className="text-left">
                  <div className="font-semibold text-white">
                    {user.first_name} {user.last_name}
                  </div>

                  <div className="text-sm text-gray-500">
                    @{user.nickname}
                  </div>
                </div>
              </button>
            ))}
        </div>
      )}
    </div>
  );
}