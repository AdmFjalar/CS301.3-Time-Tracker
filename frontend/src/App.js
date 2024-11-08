import React from 'react';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import Sidebar from './components/Sidebar';
import UpcomingShifts from './components/UpcomingShifts';
import Overview from './components/Overview';
import WorkedTime from './components/WorkedTime';
import Login from './components/Login';
import Register from './components/Register';
import Activate from './components/Activate';
import './App.css';

const App = () => {
  return (
    <Router>
      <div className="app-container">
        <Sidebar />
        <div className="main-content">
          <Routes>
            {/* Home Route */}
            <Route path="/" element={
              <>
                <h2>Welcome back, Angela</h2>
                <div className="shifts-and-overview">
                  <UpcomingShifts />
                  <Overview />
                </div>
              </>
            } />

            {/* Worked Time Route */}
            <Route path="/worked-time" element={<WorkedTime />} />

            {/* Authentication Routes */}
            <Route path="/login" element={<Login />} />
            <Route path="/register" element={<Register />} />
            <Route path="/activate" element={<Activate />} />
          </Routes>
        </div>
      </div>
    </Router>
  );
};

export default App;


