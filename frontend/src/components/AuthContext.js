import React, { createContext, useState, useContext } from 'react';

const AuthContext = createContext();

export const useAuth = () => useContext(AuthContext);

export const AuthProvider = ({ children }) => {
  const [authToken, setAuthToken] = useState(null); // Store auth token

  // Sign in function
  const signIn = (token) => {
    setAuthToken(token); // Set the auth token when the user logs in
  };

  // Sign out function
  const signOut = () => {
    setAuthToken(null); // Clear the auth token when the user signs out
  };

  return (
    <AuthContext.Provider value={{ authToken, signIn, signOut }}>
      {children}
    </AuthContext.Provider>
  );
};
