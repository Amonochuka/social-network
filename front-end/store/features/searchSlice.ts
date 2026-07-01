import { createAsyncThunk, createSlice } from "@reduxjs/toolkit";
import { RootState } from "../store";
import { searchService, SearchUser} from "@/services/searchService";

interface SearchState {
  results: SearchUser[];
  loading: boolean;
  error: string | null;
}

const initialState: SearchState = {
  results: [],
  loading: false,
  error: null,
};

export const searchUsers = createAsyncThunk(
  "search/searchUsers",
  async (query: string, { rejectWithValue }) => {
    try {
      return await searchService.searchUsers(query);
    } catch (err: any) {
      return rejectWithValue(err.response?.data || "Search failed");
    }
  }
);

const searchSlice = createSlice({
  name: "search",
  initialState,
  reducers: {
    clearResults(state) {
      state.results = [];
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(searchUsers.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(searchUsers.fulfilled, (state, action) => {
    state.loading = false;
    state.results = action.payload.users;
})
      .addCase(searchUsers.rejected, (state, action: any) => {
        state.loading = false;
        state.error = action.payload;
      });
  },
});

export const { clearResults } = searchSlice.actions;

export const searchSelector = (state: RootState) => state.search;

export default searchSlice.reducer;