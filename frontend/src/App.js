import React from 'react';
import Sidebar from './components/Sidebar';
import UpcomingShifts from './components/UpcomingShifts';
import Overview from './components/Overview';
import './App.css';

const App = () => {
  return (
    <div className="app-container">
      <Sidebar />
      <div className="main-content">
        <h2>Welcome back, Angela</h2>
        <div className="shifts-and-overview">
          <UpcomingShifts />
          <Overview />
        </div>
      </div>
    </div>
  );
};

export default App;
