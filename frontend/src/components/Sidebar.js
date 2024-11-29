import React, { useEffect, useState } from 'react';
import { useAuth } from '../components/AuthContext'; // Import the useAuth hook
import './Sidebar.css';
import { useNavigate, useLocation } from 'react-router-dom'; // Import useNavigate and useLocation
import axios from 'axios';

const Sidebar = () => {
  const { authToken, signOut } = useAuth(); // Get the authToken and signOut function from the context
  const [userRoleLevel, setUserRoleLevel] = useState(null); // State to store the user's role level
  const navigate = useNavigate();
  const location = useLocation(); // Get the current location

  useEffect(() => {
    if (!authToken) {
      if (location.pathname === '/' || location.pathname === '/worked-time' || location.pathname === '/user-details') {
      navigate('/login');
      return;
      } // Stop execution if no authToken
    }
    // Fetch user data to determine their role
    const fetchUserRole = async () => {
      try {
        const response = await axios.get('http://localhost:8080/v1/', {
          headers: {
            Authorization: `Bearer ${authToken}`,
          },
        });
        setUserRoleLevel(response.data.data.role?.level); // Assume the role's level is part of the user data
      } catch (error) {
        console.error('Failed to fetch user role:', error);
      }
    };

    if (authToken) {
      fetchUserRole();
    }
  }, [authToken]);

  const handleSignOut = () => {
    signOut(); // Call the signOut function to clear the token
    navigate('/login'); // Redirect to login page
  };

  // Helper function to determine if a link is active
  const isActive = (path) => location.pathname === path;

  return (
    <div className="sidebar">
      <pre>
        <h1 className="page-title">☁ Thyme Flies</h1>
      </pre>
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
            <button
              className={`nav-link ${isActive('/user-details') ? 'active' : ''}`}
              onClick={() => navigate('/user-details')}
            >
              User Details
            </button>
            {userRoleLevel && userRoleLevel >= 2 && ( // Only show for managers (level 2) or admins (level 3+)
              <button
                className={`nav-link ${isActive('/users') ? 'active' : ''}`}
                onClick={() => navigate('/users')}
              >
                Manage Users
              </button>
            )}
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
