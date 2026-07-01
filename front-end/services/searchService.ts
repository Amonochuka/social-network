import { Api } from "./axios";

export interface SearchUser {
  id: string;
  first_name: string;
  last_name: string;
  nickname: string;
  avatar: string;
  is_private: boolean;
}

export interface SearchGroup {
  id: string;
  title: string;
  avatar: string;
}

export interface SearchResponse {
  users: SearchUser[];
  groups: SearchGroup[];
}

export const searchService = {
  searchUsers: async (query: string): Promise<SearchResponse> => {
    const res = await Api.get(`/search?q=${encodeURIComponent(query)}`);
    return res.data;
  },
};