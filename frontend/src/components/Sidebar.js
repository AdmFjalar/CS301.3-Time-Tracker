import React from 'react';
import { useAuth } from '../components/AuthContext'; // Import the useAuth hook
import './Sidebar.css';
import { useNavigate } from 'react-router-dom';

const Sidebar = () => {
  const { authToken, signOut } = useAuth(); // Get the authToken and signOut function from the context
  const navigate = useNavigate();

  const handleSignOut = () => {
    signOut(); // Call the signOut function to clear the token
    navigate('/login'); // Redirect to login page
  };

  return (
    <div className="sidebar">
      <h2 className="page-title">Thyme Flies</h2>
      {authToken && ( // Only display if user is signed in
        <>
          <nav>
            <a href="/">Home</a>
            <a href="/worked-time">Worked time</a>
          </nav>
          <button className="signout" onClick={handleSignOut}>
            Sign out
          </button>
        </>
      )}
      {/* Additional elements can be added below */}
    </div>
  );
};

export default Sidebar;
