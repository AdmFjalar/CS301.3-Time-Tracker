import React from 'react';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import { AuthProvider } from './components/AuthContext'; // Import AuthProvider
import Sidebar from './components/Sidebar';
import UpcomingShifts from './components/UpcomingShifts';
import HomePage from './components/HomePage';
import Overview from './components/Overview';
import WorkedTime from './components/WorkedTime';
import Login from './components/Login';
import Register from './components/Register';
import Activate from './components/Activate';
import './App.css';

const App = () => {
  return (
    <AuthProvider>
      <Router>
        <div className="app-container">
          <Sidebar />
          <div className="main-content">
            <Routes>
              <Route path="/" element={<HomePage />} />
              <Route path="/worked-time" element={<WorkedTime />} />
              <Route path="/login" element={<Login />} />
              <Route path="/register" element={<Register />} />
              <Route path="/confirm/:token" element={<Activate />} />
            </Routes>
          </div>
        </div>
      </Router>
    </AuthProvider>
  );
};

export default App;
