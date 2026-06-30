import { createSlice, createAsyncThunk } from "@reduxjs/toolkit";
import { RootState } from "../store";
import { followerService } from "@/services/followerService";

export interface Follower {
  user_id: string;
  first_name: string;
  last_name: string;
  avatar: string;
  nickname: string;
}

interface FollowerState {
  followers: Follower[];
  following: Follower[];
  loading: boolean;
  error: string | null;
}

const initialState: FollowerState = {
  followers: [],
  following: [],
  loading: false,
  error: null,
};

export const getFollowers = createAsyncThunk(
  "followers/getFollowers",
  async (_, { rejectWithValue }) => {
    try {
      return await followerService.getFollowers();
    } catch (err: any) {
      return rejectWithValue(err.response?.data || "Failed to load followers");
    }
  }
);

export const getFollowing = createAsyncThunk(
  "followers/getFollowing",
  async (_, { rejectWithValue }) => {
    try {
      return await followerService.getFollowing();
    } catch (err: any) {
      return rejectWithValue(err.response?.data || "Failed to load following");
    }
  }
);

export const sendFollowRequest = createAsyncThunk(
  "followers/sendFollowRequest",
  async (receiverId: string, { rejectWithValue }) => {
    try {
      return await followerService.sendFollowRequest(receiverId);
    } catch (err: any) {
      return rejectWithValue(
        err.response?.data || "Failed to send follow request"
      );
    }
  }
);

export const acceptFollowRequest = createAsyncThunk(
  "followers/acceptFollowRequest",
  async (requestId: string, { rejectWithValue }) => {
    try {
      return await followerService.acceptFollowRequest(requestId);
    } catch (err: any) {
      return rejectWithValue(
        err.response?.data || "Failed to accept follow request"
      );
    }
  }
);

export const declineFollowRequest = createAsyncThunk(
  "followers/declineFollowRequest",
  async (requestId: string, { rejectWithValue }) => {
    try {
      return await followerService.declineFollowRequest(requestId);
    } catch (err: any) {
      return rejectWithValue(
        err.response?.data || "Failed to decline follow request"
      );
    }
  }
);

export const unfollowUser = createAsyncThunk(
  "followers/unfollowUser",
  async (userId: string, { rejectWithValue }) => {
    try {
      return await followerService.unfollowUser(userId);
    } catch (err: any) {
      return rejectWithValue(err.response?.data || "Failed to unfollow user");
    }
  }
);

const followerSlice = createSlice({
  name: "follower",
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      // Get Followers
      .addCase(getFollowers.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(getFollowers.fulfilled, (state, action) => {
        state.loading = false;
        state.followers = action.payload;
      })
      .addCase(getFollowers.rejected, (state, action: any) => {
        state.loading = false;
        state.error = action.payload;
      })

      // Get Following
      .addCase(getFollowing.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(getFollowing.fulfilled, (state, action) => {
        state.loading = false;
        state.following = action.payload;
      })
      .addCase(getFollowing.rejected, (state, action: any) => {
        state.loading = false;
        state.error = action.payload;
      })

      // Send Follow Request
      .addCase(sendFollowRequest.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(sendFollowRequest.fulfilled, (state) => {
        state.loading = false;
      })
      .addCase(sendFollowRequest.rejected, (state, action: any) => {
        state.loading = false;
        state.error = action.payload;
      })

      // Accept Follow Request
      .addCase(acceptFollowRequest.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(acceptFollowRequest.fulfilled, (state) => {
        state.loading = false;
      })
      .addCase(acceptFollowRequest.rejected, (state, action: any) => {
        state.loading = false;
        state.error = action.payload;
      })

      // Decline Follow Request
      .addCase(declineFollowRequest.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(declineFollowRequest.fulfilled, (state) => {
        state.loading = false;
      })
      .addCase(declineFollowRequest.rejected, (state, action: any) => {
        state.loading = false;
        state.error = action.payload;
      })

      // Unfollow User
      .addCase(unfollowUser.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(unfollowUser.fulfilled, (state) => {
        state.loading = false;
      })
      .addCase(unfollowUser.rejected, (state, action: any) => {
        state.loading = false;
        state.error = action.payload;
      });
  },
});

export const followerSelector = (state: RootState) => state.followers;

export default followerSlice.reducer;