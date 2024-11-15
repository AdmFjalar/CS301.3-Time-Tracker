import React, { createContext, useState, useEffect, useContext } from 'react';

const AuthContext = createContext();

export const useAuth = () => useContext(AuthContext);

export const AuthProvider = ({ children }) => {
  const [authToken, setAuthToken] = useState(() => {
    // Retrieve the token from localStorage (if it exists) when the component mounts
    return localStorage.getItem('authToken') || null;
  });

  // Sign in function: Set the token and store it in localStorage
  const signIn = (token) => {
    setAuthToken(token); 
    localStorage.setItem('authToken', token); // Store the token in localStorage
  };

  // Sign out function: Clear the token from state and localStorage
  const signOut = () => {
    setAuthToken(null); 
    localStorage.removeItem('authToken'); // Remove token from localStorage
  };

  return (
    <AuthContext.Provider value={{ authToken, signIn, signOut }}>
      {children}
    </AuthContext.Provider>
  );
};
