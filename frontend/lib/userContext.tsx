"use client";
import React, { createContext, useContext } from 'react';
import { User } from '@/types';

// Context to provide the currently authenticated user throughout the app
export const UserContext = createContext<User | null>(null);
export const useCurrentUser = () => useContext(UserContext);

export const UserProvider = UserContext.Provider;

export default UserContext;
