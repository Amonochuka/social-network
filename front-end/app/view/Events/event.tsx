"use client";

import { useEffect, useState, useMemo } from "react";
import { Calendar, MapPin, Users, Check, XCircle, Clock } from "lucide-react";
import { Api } from "@/services/axios";

interface EventDetail {
  id: string;
  group_id: string;
  creator_id: string;
  creator_name: string;
  title: string;
  description: string;
  event_time: string;
  created_at: string;
  going_count: number;
  not_going_count: number;
  my_response: string;
  groupTitle?: string;
}

interface Group {
  id: string;
  creator_id: string;
  title: string;
  description: string;
  created_at: string;
}

type TabType = "Upcoming" | "Going" | "Not Going";

export default function EventsPage() {
  const [events, setEvents] = useState<EventDetail[]>([]);
  const [loading, setLoading] = useState(true);
  const [activeTab, setActiveTab] = useState<TabType>("Upcoming");

  const fetchAllEvents = async () => {
    setLoading(true);
    try {
      const groupsRes = await Api.get<Group[]>("/groups");
      const groups = groupsRes.data ?? [];
      const fetchedEvents: EventDetail[] = [];

      await Promise.all(
        groups.map(async (group) => {
          try {
            const eventsRes = await Api.get<EventDetail[]>(`/groups/${group.id}/events`);
            if (eventsRes.data) {
              const tagged = eventsRes.data.map((e) => ({
                ...e,
                groupTitle: group.title,
              }));
              fetchedEvents.push(...tagged);
            }
          } catch {
            // Ignore 403 or loading errors for groups user hasn't joined
          }
        })
      );

      // Sort by event time ascending
      fetchedEvents.sort(
        (a, b) => new Date(a.event_time).getTime() - new Date(b.event_time).getTime()
      );
      setEvents(fetchedEvents);
    } catch (err) {
      console.error("Failed to load global events:", err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchAllEvents();
  }, []);

  const handleRSVP = async (eventId: string, response: string) => {
    try {
      await Api.post(`/groups/events/${eventId}/rsvp`, { response });
      // Update local state instantly
      setEvents((prev) =>
        prev.map((ev) => {
          if (ev.id === eventId) {
            let goingDiff = 0;
            let notGoingDiff = 0;

            // adjust old response counts
            if (ev.my_response === "going") goingDiff--;
            if (ev.my_response === "not_going") notGoingDiff--;

            // adjust new response counts
            if (response === "going") goingDiff++;
            if (response === "not_going") notGoingDiff++;

            return {
              ...ev,
              my_response: response,
              going_count: Math.max(0, ev.going_count + goingDiff),
              not_going_count: Math.max(0, ev.not_going_count + notGoingDiff),
            };
          }
          return ev;
        })
      );
    } catch (err) {
      console.error("Failed to RSVP:", err);
    }
  };

  const filteredEvents = useMemo(() => {
    const now = new Date().getTime();
    return events.filter((ev) => {
      const eventTime = new Date(ev.event_time).getTime();
      if (activeTab === "Upcoming") {
        return eventTime >= now;
      } else if (activeTab === "Going") {
        return ev.my_response === "going";
      } else if (activeTab === "Not Going") {
        return ev.my_response === "not_going";
      }
      return true;
    });
  }, [events, activeTab]);

  return (
    <main className="h-full w-full px-8 py-6 bg-[#181818] overflow-y-auto">
      <div className="flex h-full w-full flex-col">
        {/* Header */}
        <div className="mb-6">
          <h1 className="text-3xl font-bold text-white">Events</h1>
          <p className="mt-1 text-sm text-gray-400">
            View and respond to events from all the groups you have joined.
          </p>
        </div>

        {/* Tabs */}
        <div className="mb-6 flex gap-6 border-b border-white/10 pb-2">
          {(["Upcoming", "Going", "Not Going"] as TabType[]).map((tab) => (
            <button
              key={tab}
              onClick={() => setActiveTab(tab)}
              className={`font-semibold text-sm pb-2 transition relative ${
                activeTab === tab ? "text-white" : "text-gray-500 hover:text-white"
              }`}
            >
              {tab}
              {activeTab === tab && (
                <div className="absolute bottom-0 left-0 right-0 h-0.5 bg-[--primary-theme]" />
              )}
            </button>
          ))}
        </div>

        {/* Events list */}
        <div className="flex-1 space-y-4">
          {loading ? (
            <div className="py-16 text-center text-gray-500">Loading events...</div>
          ) : filteredEvents.length === 0 ? (
            <div className="py-16 text-center text-gray-500 rounded-2xl border border-white/5 bg-[#1e1e1e]">
              No {activeTab.toLowerCase()} events found.
            </div>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {filteredEvents.map((event) => {
                const formattedTime = new Date(event.event_time).toLocaleString(undefined, {
                  weekday: "short",
                  month: "short",
                  day: "numeric",
                  hour: "2-digit",
                  minute: "2-digit",
                });

                return (
                  <div
                    key={event.id}
                    className="rounded-2xl border border-white/10 bg-[#1e1e1e] p-5 flex flex-col justify-between hover:border-white/20 transition duration-200"
                  >
                    <div>
                      {/* Group tag */}
                      {event.groupTitle && (
                        <span className="inline-block text-[10px] uppercase font-bold tracking-wider text-[--primary-theme] bg-[--primary-theme]/10 rounded-full px-2.5 py-1 mb-3">
                          {event.groupTitle}
                        </span>
                      )}

                      <h2 className="text-lg font-bold text-white leading-snug">
                        {event.title}
                      </h2>
                      <p className="mt-2 text-sm text-gray-400 leading-relaxed line-clamp-3">
                        {event.description}
                      </p>

                      <div className="mt-4 space-y-2">
                        <div className="flex items-center gap-2 text-xs text-gray-500">
                          <Calendar size={14} className="text-gray-600" />
                          <span>{formattedTime}</span>
                        </div>
                        <div className="flex items-center gap-2 text-xs text-gray-500">
                          <Users size={14} className="text-gray-600" />
                          <span>
                            Organized by <strong className="text-gray-400">{event.creator_name}</strong>
                          </span>
                        </div>
                      </div>
                    </div>

                    <div className="mt-5 pt-4 border-t border-white/5 flex items-center justify-between">
                      <div className="flex gap-3 text-xs text-gray-500">
                        <span>
                          <strong className="text-emerald-400">{event.going_count}</strong> going
                        </span>
                        <span>
                          <strong className="text-red-400">{event.not_going_count}</strong> not going
                        </span>
                      </div>

                      <div className="flex gap-2">
                        <button
                          onClick={() => handleRSVP(event.id, "going")}
                          className={`flex items-center gap-1 rounded-lg px-3 py-1.5 text-xs font-semibold transition ${
                            event.my_response === "going"
                              ? "bg-emerald-500/20 text-emerald-400 border border-emerald-500/30"
                              : "bg-white/5 text-gray-400 hover:bg-white/10"
                          }`}
                        >
                          <Check size={14} /> Going
                        </button>
                        <button
                          onClick={() => handleRSVP(event.id, "not_going")}
                          className={`flex items-center gap-1 rounded-lg px-3 py-1.5 text-xs font-semibold transition ${
                            event.my_response === "not_going"
                              ? "bg-red-500/20 text-red-400 border border-red-500/30"
                              : "bg-white/5 text-gray-400 hover:bg-white/10"
                          }`}
                        >
                          <XCircle size={14} /> Decline
                        </button>
                      </div>
                    </div>
                  </div>
                );
              })}
            </div>
          )}
        </div>
      </div>
    </main>
  );
}