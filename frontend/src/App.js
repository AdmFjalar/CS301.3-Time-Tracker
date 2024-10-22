import React from 'react';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import Sidebar from './components/Sidebar';
import UpcomingShifts from './components/UpcomingShifts';
import Overview from './components/Overview';
import WorkedTime from './components/WorkedTime'; // New page
import './App.css';

const App = () => {
  return (
    <Router>
      <div className="app-container">
        <Sidebar />
        <div className="main-content">
          <Routes>
            <Route path="/" element={
              <>
                <h2>Welcome back, Angela</h2>
                <div className="shifts-and-overview">
                  <UpcomingShifts />
                  <Overview />
                </div>
              </>
            } />
            <Route path="/worked-time" element={<WorkedTime />} />
          </Routes>
        </div>
      </div>
    </Router>
  );
};

export default App;

