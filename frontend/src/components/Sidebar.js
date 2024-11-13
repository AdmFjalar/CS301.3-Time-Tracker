import React from 'react';
import { useAuth } from '../components/AuthContext'; // Import the useAuth hook
import './Sidebar.css'; 

const Sidebar = () => {
  const { signOut } = useAuth(); // Get the signOut function from the context

  const handleSignOut = () => {
    signOut(); // Call the signOut function to clear the token
    // You may also want to redirect the user to the login page after sign out
    window.location.href = '/login'; // Redirect to login page (or use navigate)
  };

  return (
    <div className="sidebar">
      <h2>Thyme Flies</h2>
      <nav>
        <a href="/">Home</a>
        <a href="/worked-time">Worked time</a>
      </nav>
      <button className="signout" onClick={handleSignOut}>
        Sign out
      </button>
    </div>
  );
};

export default Sidebar;
