import { createSlice, PayloadAction } from '@reduxjs/toolkit';

export interface AppState {
  credits: number;
  userProfile: {
    name: string;
    id: string;
  };
  aiPersonas: any[];
  isAuthenticated: boolean;
  subscriptionTier: 'free' | 'premium';
}

const initialState: AppState = {
  credits: 100, // Starts with 100 credits
  userProfile: {
    name: 'Real User',
    id: 'u1',
  },
  aiPersonas: [],
  isAuthenticated: false, // Default to not logged in
  subscriptionTier: 'free',
};

export const appSlice = createSlice({
  name: 'app',
  initialState,
  reducers: {
    deductCredits: (state, action: PayloadAction<number>) => {
      state.credits = Math.max(0, state.credits - action.payload);
    },
    addPersona: (state, action: PayloadAction<any>) => {
      state.aiPersonas.push(action.payload);
    },
    login: (state) => {
      state.isAuthenticated = true;
    },
    logout: (state) => {
      state.isAuthenticated = false;
    },
    setSubscriptionTier: (state, action: PayloadAction<'free' | 'premium'>) => {
      state.subscriptionTier = action.payload;
    },
  },
});

export const { deductCredits, addPersona, login, logout, setSubscriptionTier } = appSlice.actions;
export default appSlice.reducer;
