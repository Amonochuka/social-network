import { configureStore } from "@reduxjs/toolkit";
import toggleBarSlice from "@/store/features/toggleSideBarSlice";
import authSlice from "@/store/features/authSlice";
import followerSlice from "@/store/features/followerSlice";
import searchSlice from "@/store/features/searchSlice";

export const store = configureStore({
  reducer: {
  auth: authSlice,
  toggleSideBar: toggleBarSlice,
  followers: followerSlice,
  search: searchSlice,
},
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;