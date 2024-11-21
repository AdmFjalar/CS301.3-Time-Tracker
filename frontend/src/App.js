import React from 'react';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import { AuthProvider } from './components/AuthContext'; // Import AuthProvider
import HomePage from './components/HomePage';
import WorkedTime from './components/WorkedTime';
import WorkedTimePage from './components/WorkedTimePage';
import Login from './components/Login';
import Register from './components/Register';
import Activate from './components/Activate';
import AppLayout from './AppLayout'; 
import './App.css';

const App = () => {
  return (
    <AuthProvider>
      <Router>
        <Routes>
          {/* Wrap all routes in the AppLayout which handles layout logic */}
          <Route path="/" element={<AppLayout />}>
            <Route path="/" element={<HomePage />} />
            <Route path="/worked-time" element={<WorkedTimePage />} />
            <Route path="/confirm/:token" element={<Activate />} />
            <Route path="/login" element={<Login />} /> 
            <Route path="/register" element={<Register />} />
          </Route>
        
        </Routes>
      </Router>
    </AuthProvider>
  );
};

export default App;

