import React from 'react';
import { useLocation, Outlet } from 'react-router-dom'; // Import Outlet for nested routing
import Sidebar from './components/Sidebar'; // Import Sidebar component
import './App.css'; // Make sure the layout styles are applied

const AppLayout = () => {
  const location = useLocation();
  const showSidebar = location.pathname !== '/login'; // Only show sidebar when not on the login page

  return (
    <div className="app-container">
      {showSidebar && <Sidebar />} {/* Conditionally render Sidebar */}
      <div className="main-content">
        <Outlet /> {/* Render nested routes */}
      </div>
    </div>
  );
};

export default AppLayout;