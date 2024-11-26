import React, { createContext, useState, useEffect, useContext } from 'react';
import { useNavigate } from 'react-router-dom';
import { jwtDecode } from 'jwt-decode';

const AuthContext = createContext();

export const useAuth = () => useContext(AuthContext);

export const AuthProvider = ({ children }) => {
  const [authToken, setAuthToken] = useState(() => {
    // Retrieve the token from localStorage (if it exists) when the component mounts
    return localStorage.getItem('authToken') || null;
  });

  // Timer reference to clear on token change
  let tokenExpirationTimer = null;

  // Sign in function: Set the token and store it in localStorage
  const signIn = (token) => {
    setAuthToken(token);
    localStorage.setItem('authToken', token); // Store the token in localStorage
  };

  // Sign out function: Clear the token from state and localStorage and navigate to /login
  const signOut = (navigate) => {
    setAuthToken(null);
    localStorage.removeItem('authToken'); // Remove token from localStorage
    if (tokenExpirationTimer) {
      clearTimeout(tokenExpirationTimer); // Clear any existing timer
    }
    navigate('/login'); // Redirect to the login page
  };

  // Monitor token expiration
  useEffect(() => {
    if (authToken) {
      try {
        const decodedToken = jwtDecode(authToken);
        const currentTime = Date.now();
        const expiryTime = decodedToken.exp * 1000; // Convert exp to milliseconds
        const remainingTime = expiryTime - currentTime;

        if (remainingTime <= 0) {
          // Token already expired, sign out immediately
          signOut();
        } else {
          // Schedule sign-out when the token expires
          tokenExpirationTimer = setTimeout(signOut, remainingTime);
        }
      } catch (error) {
        console.error("Invalid token:", error);
        signOut(); // Handle invalid token by signing out
      }
    } else {
      if (tokenExpirationTimer) {
        clearTimeout(tokenExpirationTimer); // Clear timer if there's no token
      }
    }

    return () => {
      if (tokenExpirationTimer) {
        clearTimeout(tokenExpirationTimer); // Cleanup timer on unmount
      }
    };
  }, [authToken]);

  return (
    <AuthContext.Provider value={{ authToken, signIn, signOut }}>
      {children}
    </AuthContext.Provider>
  );
};
