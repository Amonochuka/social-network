import { Api } from "./axios";

export const followerService = {
  sendFollowRequest: async (receiverId: string) => {
    const res = await Api.post("/follow/requests", {
      receiver_id: receiverId,
    });
    return res.data;
  },

  acceptFollowRequest: async (requestId: string) => {
    const res = await Api.post(`/follow/requests/${requestId}/accept`);
    return res.data;
  },

  declineFollowRequest: async (requestId: string) => {
    const res = await Api.post(`/follow/requests/${requestId}/decline`);
    return res.data;
  },

  unfollowUser: async (userId: string) => {
    const res = await Api.delete(`/follow/${userId}`);
    return res.data;
  },

  getFollowers: async () => {
    const res = await Api.get("/followers");
    return res.data;
  },

  getFollowing: async () => {
    const res = await Api.get("/following");
    return res.data;
  },
};