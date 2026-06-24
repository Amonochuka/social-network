import { createSlice, createAsyncThunk } from "@reduxjs/toolkit";
import { Api } from "@/services/axios";
import { RootState } from "../store";

export interface User {
  id: string;
  email: string;
  first_name: string;
  last_name: string;
  date_of_birth: string;
  avatar: string;
  nickname: string;
  about_me: string;
  is_public: boolean;
}

interface AuthState {
  user: User | null;
  loading: boolean;
  error: string | null;
  isAuthenticated: boolean;
}

const initialState: AuthState = {
  user: null,
  loading: false,
  error: null,
  isAuthenticated: false,
};

export const loginUser = createAsyncThunk(
  "auth/login",
  async (data: { email: string; password: string }, { rejectWithValue }) => {
    try {
      const res = await Api.post("/auth/login", data);
      return res.data; // backend returns the user object directly, not wrapped in { user: ... }
    } catch (err: any) {
      return rejectWithValue(err.response?.data || "Login failed");
    }
  }
);

export const registerUser = createAsyncThunk(
  "auth/register",
  async (data: {
    email: string;
    password: string;
    first_name: string;
    last_name: string;
    date_of_birth: string;
    nickname?: string;
    about_me?: string;
  }, { rejectWithValue }) => {
    try {
      const res = await Api.post("/auth/register", data); // JSON, not FormData
      return res.data;
    } catch (err: any) {
      return rejectWithValue(err.response?.data || "Register failed");
    }
  }
);

const authSlice = createSlice({
  name: "auth",
  initialState,
  reducers: {
    setSession: (state, action) => {
      state.user = action.payload;
      state.isAuthenticated = true;
    },
    logout: (state) => {
      state.user = null;
      state.isAuthenticated = false;
      state.error = null;
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(loginUser.fulfilled, (state, action) => {
        state.user = action.payload;
        state.isAuthenticated = true;
        state.loading = false;
      })
      .addCase(loginUser.pending, (state) => {
        state.loading = true;
      })
      .addCase(loginUser.rejected, (state, action: any) => {
        state.loading = false;
        state.error = action.payload;
      })
      .addCase(registerUser.fulfilled, (state, action) => {
        state.user = action.payload;
        state.isAuthenticated = true;
        state.loading = false;
      })
      .addCase(registerUser.pending, (state) => {
        state.loading = true;
      })
      .addCase(registerUser.rejected, (state, action: any) => {
        state.loading = false;
        state.error = action.payload;
      });
  },
});

export const authSelector = (state: RootState) => state.auth;
export const { setSession, logout } = authSlice.actions;
export default authSlice.reducer;