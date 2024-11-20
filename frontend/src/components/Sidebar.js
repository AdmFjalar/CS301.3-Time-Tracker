import React from 'react';
import { useAuth } from '../components/AuthContext'; // Import the useAuth hook
import './Sidebar.css';
import { useNavigate, useLocation } from 'react-router-dom'; // Import useNavigate and useLocation

const Sidebar = () => {
  const { authToken, signOut } = useAuth(); // Get the authToken and signOut function from the context
  const navigate = useNavigate();
  const location = useLocation(); // Get the current location

  const handleSignOut = () => {
    signOut(); // Call the signOut function to clear the token
    navigate('/login'); // Redirect to login page
  };

  // Helper function to determine if a link is active
  const isActive = (path) => location.pathname === path;

  return (
    <div className="sidebar">
      <h1 className="page-title">Thyme Flies</h1>
      <p className="page-subtitle">Time Tracker</p>
      {authToken && ( // Only display if user is signed in
        <>
          <nav>
            <button
              className={`nav-link ${isActive('/') ? 'active' : ''}`}
              onClick={() => navigate('/')}
            >
              Home
            </button>
            <button
              className={`nav-link ${isActive('/worked-time') ? 'active' : ''}`}
              onClick={() => navigate('/worked-time')}
            >
              Worked Time
            </button>
          </nav>
          <button className="signout" onClick={handleSignOut}>
            Sign out
          </button>
        </>
      )}
    </div>
  );
};

export default Sidebar;
