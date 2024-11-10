import React from 'react';
import './Sidebar.css'; 

const Sidebar = () => {
  return (
    <div className="sidebar">
      <h2>Thyme Flies</h2>
      <nav>
        <a href="/">Home</a>
        <a href="/worked-time">Worked time</a>
      </nav>
      <button className="signout">Sign out</button>
    </div>
  );
};

export default Sidebar;
