"use client";

import { useEffect, useMemo, useState } from "react";
import {
  Search,
  Plus,
  Users,
  X,
} from "lucide-react";
import { Api } from "@/services/axios";
import Link from "next/link";

interface Group {
  id: string;
  creator_id: string;
  title: string;
  description: string;
  created_at: string;
}

export default function GroupsPage() {
  const [query, setQuery] = useState("");
  const [groups, setGroups] = useState<Group[]>([]);
  const [loading, setLoading] = useState(true);
  const [showCreateModal, setShowCreateModal] = useState(false);

  const fetchGroups = async () => {
    try {
      const res = await Api.get<Group[]>("/groups");
      setGroups(res.data ?? []);
    } catch (err) {
      console.error("Failed to load groups:", err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchGroups();
  }, []);

  const filteredGroups = useMemo(() => {
    return groups.filter((group) =>
      group.title.toLowerCase().includes(query.toLowerCase())
    );
  }, [query, groups]);

  return (
    <main className="h-full w-full px-8 py-6">
      <section className="flex h-full w-full flex-col">
        {/* Header */}
        <div className="mb-6 flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold text-white">
              Groups
            </h1>
            <p className="mt-1 text-sm text-muted">
              Discover and join communities.
            </p>
          </div>

          <button
            onClick={() => setShowCreateModal(true)}
            className="
              flex items-center gap-2
              rounded-xl
              bg-(--primary-theme)
              px-5 py-2.5
              text-sm font-semibold
              text-white
              transition-opacity
              hover:opacity-90
            "
          >
            <Plus size={16} />
            Create Group
          </button>
        </div>

        {/* Search */}
        <div className="mb-6">
          <div
            className="
              flex max-w-md items-center gap-3
              rounded-xl
              border border-white/10
              bg-(--fade-background)
              px-4 py-3
            "
          >
            <Search
              size={18}
              className="text-muted-dim"
            />

            <input
              value={query}
              onChange={(e) => setQuery(e.target.value)}
              placeholder="Search groups..."
              className="
                w-full
                bg-transparent
                text-sm
                text-white
                placeholder:text-muted-dim
                focus:outline-none
              "
            />
          </div>
        </div>

        {/* Group List */}
        <div className="flex-1 overflow-y-auto">
          {loading ? (
            <div className="py-16 text-center text-muted">Loading groups...</div>
          ) : (
            <div className="space-y-4">
              {filteredGroups.map((group) => (
                <Link
                  href={`/view/Groups/${group.id}`}
                  key={group.id}
                  className="
                    flex items-center gap-5
                    rounded-2xl
                    border border-white/10
                    bg-panel-soft
                    p-5
                    transition-colors duration-200
                    hover:bg-(--fade-background)
                    no-underline
                  "
                  style={{ textDecoration: "none" }}
                >
                  {/* Icon placeholder */}
                  <div className="h-16 w-16 rounded-xl bg-gradient-to-br from-violet-600 to-indigo-500 flex items-center justify-center shrink-0">
                    <Users size={28} className="text-white" />
                  </div>

                  {/* Details */}
                  <div className="min-w-0 flex-1">
                    <h3 className="truncate text-base font-semibold text-white">
                      {group.title}
                    </h3>
                    <p className="mt-1 text-sm text-gray-400 line-clamp-1">
                      {group.description}
                    </p>
                    <div className="mt-2 flex items-center gap-2 text-sm text-muted">
                      <span className="text-xs text-gray-500">
                        Created {new Date(group.created_at).toLocaleDateString()}
                      </span>
                    </div>
                  </div>
                </Link>
              ))}

              {filteredGroups.length === 0 && !loading && (
                <div className="py-16 text-center text-muted">
                  {query ? `No groups found for "${query}"` : "No groups yet. Create the first one!"}
                </div>
              )}
            </div>
          )}
        </div>
      </section>

      {showCreateModal && (
        <CreateGroupModal
          onClose={() => setShowCreateModal(false)}
          onCreated={(group) => {
            setGroups((prev) => [group, ...prev]);
            setShowCreateModal(false);
          }}
        />
      )}
    </main>
  );
}

interface CreateGroupModalProps {
  onClose: () => void;
  onCreated: (group: Group) => void;
}

function CreateGroupModal({ onClose, onCreated }: CreateGroupModalProps) {
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [creating, setCreating] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleCreate = async () => {
    if (!title.trim() || !description.trim()) {
      setError("Title and description are required.");
      return;
    }
    setCreating(true);
    setError(null);

    try {
      const res = await Api.post("/groups", { title, description });
      onCreated(res.data);
    } catch (err) {
      console.error("Failed to create group:", err);
      setError("Failed to create group. Please try again.");
    } finally {
      setCreating(false);
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 p-4">
      <div className="w-full max-w-md rounded-2xl bg-[#1e1e1e] p-6 shadow-2xl">
        <div className="mb-5 flex items-center justify-between">
          <h2 className="text-lg font-bold text-white">Create Group</h2>
          <button
            onClick={onClose}
            className="rounded-lg p-1.5 text-gray-400 transition hover:bg-white/5 hover:text-white"
          >
            <X size={18} />
          </button>
        </div>

        {error && (
          <div className="mb-4 rounded-lg bg-red-500/10 px-3 py-2 text-sm text-red-400">
            {error}
          </div>
        )}

        <div className="space-y-3">
          <div>
            <label className="mb-1 block text-xs text-gray-400">Group name</label>
            <input
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              placeholder="e.g. Photography Lovers"
              className="w-full rounded-lg bg-[#262626] px-3 py-2.5 text-sm text-white outline-none focus:ring-1 focus:ring-[--primary-theme]"
            />
          </div>
          <div>
            <label className="mb-1 block text-xs text-gray-400">Description</label>
            <textarea
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              placeholder="What is this group about?"
              rows={3}
              className="w-full resize-none rounded-lg bg-[#262626] px-3 py-2.5 text-sm text-white outline-none focus:ring-1 focus:ring-[--primary-theme]"
            />
          </div>
        </div>

        <div className="mt-6 flex justify-end gap-2">
          <button
            onClick={onClose}
            className="rounded-lg px-4 py-2 text-sm font-semibold text-gray-300 transition hover:bg-white/5"
          >
            Cancel
          </button>
          <button
            onClick={handleCreate}
            disabled={creating}
            className="rounded-lg bg-[--primary-theme] px-4 py-2 text-sm font-semibold text-black transition hover:opacity-90 disabled:opacity-50"
          >
            {creating ? "Creating..." : "Create Group"}
          </button>
        </div>
      </div>
    </div>
  );
}